// Package recording provides functions for recording shows
package recording

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
	"wrbb-stream-recorder/internal/spinitron"
	"wrbb-stream-recorder/internal/util"
)

// ShowRecordingLoop is the goroutine used to continually
// check if a show is starting and record it
func ShowRecordingLoop(schedule *spinitron.ShowSchedule) {
	currentRecording = &currentRecordingStruct{
		mu:    sync.Mutex{},
		shows: map[string]spinitron.Show{},
	}
	for {
		// Get the current time, hour and minute
		now := time.Now()
		hour := now.Hour()
		// Lock the mutex
		schedule.Mu.Lock()
		// check if a show starts the current hour
		if show, starting := schedule.Schedule[hour]; starting {
			// If a show is starting now, delete show from schedule, then record
			delete(schedule.Schedule, hour)
			// Start a recording in a gorountine
			go func() {
				util.InfoLogger.Printf("Starting to record show %s", show.Name)
				err := RecordShow(show)
				if err != nil {
					util.ErrorLog(fmt.Sprintf("Unable to record show %s: %s\n", show.Name, err.Error()))
				}
			}()
		}
		// Unlock the mutex
		schedule.Mu.Unlock()
		// Sleep and check again in a minute
		time.Sleep(time.Minute)
	}
}

func createRequest() (*http.Response, error) {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 1024,
		TLSHandshakeTimeout: 0 * time.Second,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", viper.GetString(util.StreamUrl), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Keep-Alive", fmt.Sprintf("timeout=%d,max=0", int(time.Hour.Seconds())))
	return client.Do(req)
}

// RecordShow Records a given show from the StreamURL to an mp3 named the current date to
// a folder of the show names in the VortexStorageLocation directory for the shows duration
func RecordShow(show spinitron.Show) error {
	// Used to debug, if true, dont write show
	if !viper.GetBool(util.WriteShows) {
		return fmt.Errorf("not writing show due to debug flag")
	}

	// Get a connection to the stream
	response, err := createRequest()
	if err != nil {
		return err
	}

	// Create the show directory
	showDirectory := path.Join(viper.GetString(util.StorageLocation), show.Name)
	if err = os.MkdirAll(showDirectory, 0755); err != nil {
		return err
	}
	// Get date for name of mp3
	year, month, day := show.Start.Date()
	// Create and open the mp3
	f, err := os.Create(path.Join(showDirectory,
		fmt.Sprintf("%s-%d-%d.mp3", month.String(), day, year)))
	if err != nil {
		return err
	}

	currentRecording.mu.Lock()
	currentRecording.shows[show.Name] = show
	currentRecording.mu.Unlock()

	showTimeLeft := show.End.Sub(time.Now())
	copyShow(f, response.Body, showTimeLeft, show.Name)
	return nil
}

// currentRecordingStruct is the struct that represents the current
// recording list and mutex to access it
type currentRecordingStruct struct {
	mu    sync.Mutex
	shows map[string]spinitron.Show
}

// currentRecording is a list of the current shows being recorded
var currentRecording *currentRecordingStruct

// GetCurrentShows Gets the currently recording show
func GetCurrentShows() (shows []spinitron.Show) {
	shows = []spinitron.Show{}
	currentRecording.mu.Lock()
	for _, show := range currentRecording.shows {
		shows = append(shows, show)
	}
	currentRecording.mu.Unlock()
	return
}

// UpdateScheduleLoop is the goroutine to continually update the
// Spinitron spinitron.ShowSchedule at midnight every night
func UpdateScheduleLoop(schedule *spinitron.ShowSchedule) {
	timer := time.NewTimer(util.TimeUntilMidnight())
	for {
		select {
		case <-timer.C:
			err := spinitron.FetchSchedule(schedule)
			if err != nil {
				util.ErrorLog(fmt.Sprintf("Unable to fetch spinitron: %s\n", err.Error()))
			}
			timer.Reset(util.TimeUntilMidnight())
		}
	}
}
