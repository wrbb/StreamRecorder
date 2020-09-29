package main

import (
	"StreamRecorderGo/pkg/spinitron"
	"fmt"
	"net/http"
	"os"
)

func main() {
	// grab show data
	schedule := spinitron.ShowSchedule{}
	err := spinitron.GetShows(schedule)
	if err == nil {
		os.Exit(1)
	}

	go


	http.HandleFunc("/", Dashboard)
	http.ListenAndServe(":1049", nil)
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Vortex Stream Recorder is currently up and running!")
}
