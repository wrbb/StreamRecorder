package spinitron

import (
	"fmt"
	"time"
)

type Show struct {
	Start    time.Time
	End      time.Time
	Name     string
	Duration time.Duration
}

func (s Show) IsLive() bool {
	now := time.Now()
	return now.Before(s.End) && now.After(s.Start)
}

func (s Show) HasPast() bool {
	now := time.Now()
	return now.After(s.End)
}

type ShowSchedule struct {
	Shows []Show
}

func (s ShowSchedule) AppendShow(show Show) {
	s.Shows = append(s.Shows, show)
}

func (s ShowSchedule) NextShowIsLive() bool {
	if len(s.Shows) > 1 {
		return s.Shows[0].IsLive()
	}
	return false
}

func (s ShowSchedule) PopNextShow() (Show, error) {
	if len(s.Shows) < 1 {
		return Show{},  fmt.Errorf("No shows available")
	}

	nextShow := s.Shows[0]
	s.Shows[0] = Show{}
	s.Shows = s.Shows[1:]

	return nextShow,nil
}