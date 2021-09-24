// Package recording provides functions for recording shows
package recording

import (
	"fmt"
	"github.com/spf13/viper"
	"net"
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
		// Get the current hour
		hour := time.Now().Hour()
		// Lock the mutex
		schedule.Mu.Lock()
		// check if a show starts the current hour
		if show, starting := schedule.Schedule[hour]; starting {
			// If a show is starting now, delete show from schedule, then record
			delete(schedule.Schedule, hour)
			// Start a recording in a goroutine
			go func(s spinitron.Show) {
				// Print beginning of show
				util.InfoLogger.Printf("Starting to record show '%s'", s.Name)
				// Give 5 retires for recording
				retries := 5
				for {
					err := RecordShow(s)
					if err == nil {
						// If doesnt fail to record, print that the show is done and break from loop
						util.InfoLogger.Printf("Finished recording show '%s'", s.Name)
						break
					}

					if retries == 0 {
						// If fails after 5 times, state that the show failed to record, and break from loop
						util.ErrorLog(fmt.Sprintf("Unable to record show '%s': %s", s.Name, err.Error()))
						break
					}
					retries--
					util.WarningLogger.Printf("Failed to record show '%s': %s\tretrying...", s.Name, err.Error())
				}
			}(show)
		}
		// Unlock the mutex
		schedule.Mu.Unlock()
		// Sleep and check again in a minute
		time.Sleep(time.Minute)
	}
}

// createRequest creates a get request to the stream
// with a dial timeout/keep alive of timeout
func createRequest(timeout time.Duration) (*http.Response, error) {
	c := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: timeout,
			}).DialContext,
			TLSHandshakeTimeout: 0,
		},
	}
	return c.Get(viper.GetString(util.StreamUrl))
}

// createRecordingFile creates the file to record the show to.
func createRecordingFile(show spinitron.Show) (*os.File, error) {
	// Create the show directory
	showDirectory := path.Join(viper.GetString(util.StorageLocation), show.Name)
	if err := os.MkdirAll(showDirectory, 0755); err != nil {
		return nil, err
	}

	// create name of mp3 file
	year, month, day := show.Start.Date()
	filename := fmt.Sprintf("%s-%d-%d.mp3", month.String(), day, year)
	// Create and open the mp3
	return os.Create(path.Join(showDirectory,filename))
}

// RecordShow Records a given show from the StreamURL to an mp3 named the current date to
// a folder of the show names in the VortexStorageLocation directory for the shows duration
func RecordShow(show spinitron.Show) error {
	// Used to debug, if true, dont write show
	if !viper.GetBool(util.WriteShows) {
		return fmt.Errorf("not writing show due to debug flag")
	}

	ShowBeganRecording(show)
	defer ShowEndedRecording(show)
	// Get a connection to the stream
	response, err := createRequest(show.Duration)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Create file to save show to
	f, err := createRecordingFile(show)
	if err != nil {
		return err
	}

	// Begin writing show to disk
	if show.HasPast() {
		return nil
	}
	return copyShow(f, response.Body, show.End.Sub(time.Now()))
}

// currentRecordingStruct is the struct that represents the current
// recording list and mutex to access it
type currentRecordingStruct struct {
	mu    sync.Mutex
	shows map[string]spinitron.Show
}

// currentRecording is a list of the current shows being recorded
var currentRecording *currentRecordingStruct

// ShowBeganRecording adds show to the list of currently
// recording shows
func ShowBeganRecording(show spinitron.Show) {
	if currentRecording != nil {
		currentRecording.mu.Lock()
		currentRecording.shows[show.Name] = show
		currentRecording.mu.Unlock()
	}
}

// ShowEndedRecording removes show from the list of currently
// recording shows
func ShowEndedRecording(show spinitron.Show) {
	if currentRecording != nil {
		currentRecording.mu.Lock()
		delete(currentRecording.shows, show.Name)
		currentRecording.mu.Unlock()
	}
}

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
