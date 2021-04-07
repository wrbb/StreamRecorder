package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"text/template"
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
	// Fetchs the schedule for the show queue
	err := spinitron.FetchSchedule(&schedule)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to fetch data from spinitron: %v", err.Error())
	}

	// Sets a cronjob for every night at midnight to update the show schedule
	c := cron.New()
	_, err = c.AddFunc("@midnight", func() {
		err = spinitron.FetchSchedule(&schedule)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "unable to fetch data from spinitron: %v", err.Error())
		}
	})

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to start cronjob for spinitron schedule: %v", err.Error())
		os.Exit(1)
	}
	// Starts the cronjob
	c.Start()

	// Create the show channel
	showChannel := make(chan spinitron.Show)

	// Starts the recording goroutine
	go recording.RecordShowRoutine(config, showChannel)
	// Starts the schedule loop
	go recording.ScheduleLoop(&schedule, showChannel)

	http.HandleFunc("/", Dashboard)
	_ = http.ListenAndServe(":1049", nil)
}


type DashboardData struct {
	IsRecording bool
	ShowName string
}

func Dashboard(w http.ResponseWriter, _ *http.Request) {
	currentShow := recording.GetCurrentShow()
	dashboardData := DashboardData{
		IsRecording: currentShow != (spinitron.Show{}),
		ShowName:  currentShow.Name,
	}
	w.Header().Add("Content Type", "text/html")
	// The template name "template" does not matter here
	tmpl := template.Must(template.ParseFiles("web/template/dashboard.html"))
	tmpl.Execute(w, dashboardData)
}
