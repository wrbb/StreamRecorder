package spinitron

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type showResponse struct {
	Id       int64  `json:"id"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Duration int    `json:"duration"`
	Title    string `json:"title"`
}

const DateFormat = "2006-01-02T15:04:05-0700"
const AccessToken = "ARdWnef9Fie7lKWspQzn5efv"
const Count = 1000

func (s showResponse) convertToShow() (Show, error) {
	parsedStart, err := time.Parse(DateFormat, s.Start)
	if err != nil {
		return Show{}, fmt.Errorf("unable to parse start time: %v", s.Start)
	}
	parsedEnd, err := time.Parse(DateFormat, s.End)
	if err != nil {
		return Show{}, fmt.Errorf("unable to parse end time: %v", s.End)
	}
	return Show{
		Id:       s.Id,
		Name:     s.Title,
		Duration: time.Duration(s.Duration) * time.Second,
		Start:    parsedStart,
		End:      parsedEnd,
	}, nil
}

type spinitronResponse struct {
	Shows []showResponse `json:"items"`
}

func getMidnightTomorrow() string {
	year, month, day := time.Now().Add(time.Hour * 24).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Format(DateFormat)
}

// Fetch's the spinitron schedule from the current time to midnight of current day
func FetchSchedule(schedule *ShowSchedule) error {
	// Get data from spinitron
	url := fmt.Sprintf("https://spinitron.com/api/shows?access-token=%s&count=%d&end=%s", AccessToken, Count, getMidnightTomorrow())
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer response.Body.Close()
	// Parse Response
	spinitronResponse := spinitronResponse{}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &spinitronResponse)
	if err != nil {
		return err
	}
	for _, showResponse := range spinitronResponse.Shows {
		convertedShow, err := showResponse.convertToShow()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "could not convert show, %v", err.Error())
		}
		schedule.AppendShow(convertedShow)
	}

	return nil
}
