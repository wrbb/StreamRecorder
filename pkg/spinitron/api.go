package spinitron

import (
	"encoding/json"
	"fmt"
	"time"
)

type Show struct {
	Start    time.Time
	End      time.Time
	Name     string
	Duration int
}

type ShowSchedule []Show

type SpinitronResponse struct {


}

func GetShows(schedule ShowSchedule) error {

	spinitronResponse := SpinitronResponse{}
	err := json.Unmarshal(response, &spinitronResponse)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}