package spinitron

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"
	"wrbb-stream-recorder/internal/util"
)

const (
	// DateFormat is used to format the Spinitron request/response dates
	DateFormat = "2006-01-02T15:04:05-0700"
	// URL is the Spinitron API url for viewing the schedule of shows
	URL = "https://spinitron.com/api/shows?access-token=%s&count=%d&end=%s"
	// Count is the maximum count of shows in the response
	Count = 1000
)

// spinitronResponse represents the response from the
// Spinitron API
type spinitronResponse struct {
	Shows []showResponse `json:"items"`
}

// showResponse is a struct to represent show object
// from Spinitron response
type showResponse struct {
	Id       int64  `json:"id"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Duration int    `json:"duration"`
	Title    string `json:"title"`
}

// convertToShow converts a Spinitron showResponse into a Show struct
func (s showResponse) convertToShow() (Show, error) {
	// Parse the start date and time of the show
	parsedStart, err := time.Parse(DateFormat, s.Start)
	if err != nil {
		return Show{}, fmt.Errorf("unable to parse start time: %v", s.Start)
	}
	// Parse the end date and time of the show
	parsedEnd, err := time.Parse(DateFormat, s.End)
	if err != nil {
		return Show{}, fmt.Errorf("unable to parse end time: %v", s.End)
	}
	// Creat the show object
	return Show{
		Id:       s.Id,
		Name:     s.Title,
		Duration: time.Duration(s.Duration) * time.Second,
		Start:    parsedStart.In(util.TimeLoc),
		End:      parsedEnd.In(util.TimeLoc),
	}, nil
}

// getSpinitronSchedule makes a call to the Spinitron API to get the
// current schedule from the current time till 00:00 the next day
func getSpinitronSchedule() (response spinitronResponse, err error) {
	// get time till midnight
	midnight := util.GetMidnight()
	// Subtract 1 second to miss show at midnight tomorrow
	midnight = midnight.Add(-1 * time.Second)
	// Create URL with time and API Key
	url := fmt.Sprintf(URL, viper.GetString(util.SpinitronAPIKey), Count, midnight.Format(DateFormat))
	// Make the GET request to the URL
	httpResponse, err := http.Get(url)
	if err != nil {
		return
	}

	// Close connection once the function exits
	defer httpResponse.Body.Close()
	// If the response status code is not 200, something went wrong
	if httpResponse.StatusCode != 200 {
		return response, fmt.Errorf("Given non 200 response")
	}
	// Parse Response Body
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}

	// Prarse response body into response struct
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	return
}
