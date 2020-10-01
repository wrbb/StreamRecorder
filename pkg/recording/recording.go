package recording

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"wrbb-stream-recorder/pkg"
	"wrbb-stream-recorder/pkg/spinitron"
)

func ScheduleLoop(schedule *spinitron.ShowSchedule, showChannel chan spinitron.Show) {
	for {
		if schedule.NextShowHasPassed() {
			_, _ = schedule.PopNextShow()
			fmt.Print("Next show has already occurred")
		} else if schedule.NextShowIsLive(){
			nextShow, err := schedule.PopNextShow()
			if err == nil {
				showChannel <- nextShow
			}
		} else {
			fmt.Println("No show found")
		}
		time.Sleep(5 * time.Minute)
	}
}

func RecordShow(config pkg.Config, show spinitron.Show) error {
	response, err := http.Get(config.StreamURL)
	if err != nil {
		fmt.Print(err)
		return err
	}
	showDirectory := fmt.Sprintf("%s/%s", config.StorageLocation, show.Name)
	if err = os.MkdirAll(showDirectory, 0755); err != nil {
		fmt.Print(err)
		return err
	}
	f, err := os.Create(fmt.Sprintf("%s/%s-%d-%d.mp3", showDirectory, show.Start.Month(), show.Start.Day(), show.Start.Year()) )
	if err != nil {
		fmt.Print(err)
		return err
	}
	if _, err := copyShow(f,response.Body, show); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Finsihed recording %v", show.Name)
	return nil
}

func RecordShowRoutine(config pkg.Config, showChannel chan spinitron.Show) {
	for {
		select {
		case show := <-showChannel:
			err := RecordShow(config, show)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to record show %v", show.Name)
			}
		default:
			// Sleep & try again
			time.Sleep(time.Second)
		}
	}
}