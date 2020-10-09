package spinitron

import (
	"fmt"
	"time"
)

// A Struct to represent a spinitron show
type Show struct {
	// The spinitron Id of the show
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
	Shows []Show
}

// Appends a show to the schedule
func (s *ShowSchedule) AppendShow(show Show) {
	s.Shows = append(s.Shows, show)
}

// Returns true if the next show in schedule is currently live
func (s *ShowSchedule) NextShowIsLive() bool {
	if len(s.Shows) > 0 {
		return s.Shows[0].IsLive()
	}
	return false
}

// Returns true if the next show in schedule has already passed
func (s *ShowSchedule) NextShowHasPassed() bool {
	if len(s.Shows) > 0 {
		return s.Shows[0].HasPast()
	}
	return false
}

// Pops the next show in the schedule off the schedule queue
func (s *ShowSchedule) PopNextShow() (Show, error) {
	if len(s.Shows) < 1 {
		return Show{}, fmt.Errorf("No shows available from spinitron")
	}

	nextShow := s.Shows[0]
	s.Shows[0] = Show{}
	s.Shows = s.Shows[1:]

	return nextShow, nil
}
