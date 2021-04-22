# StreamRecorder

An internally built stream recorder for logging past shows

## Design

Basic design is to create 2 goroutines. One that fetches the spinitron scehdule every night at midnight, and another to check if there is a show that is suppose to be recording, and if true, start another routine to record the show to the storage directory

The design also includes an HTTP server to checking the health of the application by being able to monitor the schedule and currently recording show
![Design Doc](images/VortexRecorderDoc.png)

## Features
- Send errors to Slack channel for feedback to staff
- HTTP Endpoints to check health of application

## Running 
```shell
go run main.go
```

