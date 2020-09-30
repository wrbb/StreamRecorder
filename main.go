package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"net/http"
	"os"
	"wrbb-stream-recorder/pkg"
	"wrbb-stream-recorder/pkg/recording"
	"wrbb-stream-recorder/pkg/spinitron"
)

func main() {
	// Load config
	config := pkg.GetConfig()

	// grab show data
	schedule := spinitron.ShowSchedule{}
	err := spinitron.GetShows(schedule)
	if err != nil {
		os.Exit(1)
	}

	c := cron.New()
	c.AddFunc("@midnight", func() {
		err = spinitron.GetShows(schedule)
		if err != nil {
			println("Unable to fetch Spinitron schedule")
		}
	})
	c.Start()

	showChannel := make(chan spinitron.Show)

	go recording.RecordShowRoutine(config, showChannel)
	go recording.ScheduleLoop(&schedule, showChannel)

	http.HandleFunc("/", Dashboard)
	http.ListenAndServe(":1049", nil)
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Vortex Stream Recorder is currently up and running!")
}
