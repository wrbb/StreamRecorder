package main

import (
	"wrbb-stream-recorder/pkg/recording"
	"wrbb-stream-recorder/pkg/server"
	"wrbb-stream-recorder/pkg/spinitron"
	"wrbb-stream-recorder/pkg/util"
)

// main is the Starting point of the application
// it initializes the config and the loggers, then
// fetches the Spinitron schedule and starts two goroutines
// one to update the schedule at midnight and one to record
// shows, it then starts an http server to view the health of the application
func main() {
	// Load config
	util.InitConfig()
	// Load loggers
	util.InitLoggers()

	// Create schedule
	schedule, err := spinitron.CreateSchedule()
	if err != nil {
		util.ErrorLogger.Printf("Unable to fetch spinitron schedule: %s\n", err.Error())
	}

	// Starts the show recording loop
	go recording.ShowRecordingLoop(schedule)

	// Starts the schedule updating loop
	go recording.UpdateScheduleLoop(schedule)

	// Start server
	server.InitServer()
}
