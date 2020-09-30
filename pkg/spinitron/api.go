package spinitron

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type showResponse struct {
	Id 			int 	`json:"id"`
	Start		string 	`json:"start"`
	End 		string 	`json:"end"`
	Duration 	int 	`json:"duration"`
	Title 		string 	`json:"title"`
}

const DATE_FORMAT = "2006-01-02T15:04:05+0000"

func (s showResponse) convertToShow() Show {
	parsedStart, err := time.Parse(DATE_FORMAT, s.Start)
	if err != nil {
		fmt.Errorf("Unable to parse start time from spinitron: %v", s.Start)
	}
	parsedEnd, err := time.Parse(DATE_FORMAT, s.End)
	if err != nil {
		fmt.Errorf("Unable to parse end time from spinitron: %v", s.End)
	}
	return Show{
		Name: 		s.Title,
		Duration: 	time.Duration(1000 * 1000 * s.Duration),
		Start:  	parsedStart,
		End: 		parsedEnd,
	}
}

type spinitronResponse struct {
	Shows []showResponse `json:"items"`
}

func GetShows(schedule ShowSchedule) error {
	// Getch data from spinitron
	response, err := http.Get("https://spinitron.com/api/shows?access-token=ARdWnef9Fie7lKWspQzn5efv&count=1000")
	defer response.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	// Parse Response
	spinitronResponse := spinitronResponse{}
	body, err := ioutil.ReadAll(response.Body)
	print(string(body))
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = json.Unmarshal(body, &spinitronResponse)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Len: ", len(spinitronResponse.Shows))
	for _, showResponse := range spinitronResponse.Shows {
		convertedShow := showResponse.convertToShow()
		fmt.Printf("Show: %v, %v, %v, %v \n",convertedShow.Name, convertedShow.Duration.String(), convertedShow.Start.String(), convertedShow.End.String())
		schedule.AppendShow(showResponse.convertToShow())
	}

	return nil
}