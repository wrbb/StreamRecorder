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
	err = os.MkdirAll(showDirectory, 0755)
	if err != nil {
		fmt.Print(err)
		return err
	}
	f, err := os.Create(fmt.Sprintf("%s/%s-%d-%d.mp3", showDirectory, show.Start.Month(), show.Start.Day(), show.Start.Year()) )
	if err != nil {
		fmt.Print(err)
		return err
	}
	for show.IsLive() {
		buffer := make([]byte, 1024)
		_, err = response.Body.Read(buffer)
		if err != nil {
			fmt.Print(err)
			return err
		}
		_, err = f.Write(buffer)
		if err != nil {
			fmt.Print(err)
			return err
		}
	}

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