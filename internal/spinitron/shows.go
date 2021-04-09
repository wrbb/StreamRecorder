package spinitron

import (
	"fmt"
	"os"
	"sync"
	"time"
	"wrbb-stream-recorder/internal/util"
)

// A Struct to represent a Spinitron show
type Show struct {
	// The Spinitron Id of the show
	Id int64
	// The start time of the show
	Start time.Time
	// The end time of the show
	End time.Time
	// The name of ths show
	Name string
	// The shows duration
	Duration time.Duration
}

// Returns true if the current show is live
func (s Show) IsLive() bool {
	now := time.Now()
	return now.Before(s.End) && now.After(s.Start)
}

// Returns true if the show has passed
func (s Show) HasPast() bool {
	now := time.Now()
	return now.After(s.End)
}

// Struct to represent the scheduled shows
type ShowSchedule struct {
	Schedule map[int]Show
	Mu sync.Mutex
}

// CreateSchedule creates a show schedule by requesting
// the Spinitron API
func CreateSchedule() (schedule *ShowSchedule, err error) {
	schedule = &ShowSchedule{
		Mu: sync.Mutex{},
	}
	err = FetchSchedule(schedule)
	return
}

// FetchSchedule gets the show schedule from Spinitron's API
// it writes the show the ShowSchedule passed in, using the mutex
func FetchSchedule(schedule *ShowSchedule) (err error) {
	util.InfoLogger.Println("Fetching Spinitron schedule")

	var response spinitronResponse
	response, err = getSpinitronSchedule()
	if err != nil {
		return
	}
	shows := convertShows(response)

	schedule.Mu.Lock()
	schedule.Schedule = shows
	schedule.Mu.Unlock()

	return nil
}

// convertShows converts all the shows in the Spinitron response
// to a map of ints representing the hour they start to Show structs
func convertShows(response spinitronResponse) (shows map[int]Show) {
	shows = map[int]Show{}
	for _, showResponse := range response.Shows {
		convertedShow, err := showResponse.convertToShow()
		fmt.Println(convertedShow.Start)
		fmt.Println(convertedShow.Name)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "could not convert show, %v", err.Error())
			os.Exit(1)
		}
		shows[convertedShow.Start.Hour()] = convertedShow
	}

	return shows
}
