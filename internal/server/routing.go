package server

import (
	"html/template"
	"net/http"
	"regexp/syntax"
	"strings"
	"wrbb-stream-recorder/internal/recording"
	"wrbb-stream-recorder/internal/spinitron"
	"wrbb-stream-recorder/internal/util"
)

const DashboardTemplate = "web/template/dashboard.html"
const ScheduleTemplate = "web/template/schedule.html"

// InitServer creates the http server
func InitServer(schedule *spinitron.ShowSchedule) {
	http.HandleFunc("/", Dashboard)
	http.HandleFunc("/schedule", CreateScheduleViewHandler(schedule))
	_ = http.ListenAndServe(":1049", nil)
}

// CreateScheduleViewHandler creates the handler function for the schedule view
func CreateScheduleViewHandler(schedule *spinitron.ShowSchedule) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("Content Type", "text/html")
		tmpl := template.Must(template.New("schedule").Funcs(template.FuncMap{
			"nameFormat": spinitronNameFormat,
		}).ParseFiles(ScheduleTemplate))
		err := tmpl.Execute(w, schedule.Schedule)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Unable to parse template"))
			util.ErrorLogger.Printf("Unable to execute template: %s\n", err.Error())
		}
	}
}

// Dashboard in the handler function for the http server's
// index request
func Dashboard(w http.ResponseWriter, _ *http.Request) {
	currentShows := recording.GetCurrentShows()
	w.Header().Add("Content Type", "text/html")
	tmpl := template.Must(template.New("dashboard").Funcs(template.FuncMap{
		"nameFormat": spinitronNameFormat,
	}).ParseFiles(DashboardTemplate))
	err := tmpl.Execute(w, currentShows)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Unable to parse template"))
		util.ErrorLogger.Printf("Unable to execute template: %s\n", err.Error())
	}
}

// spinitronNameFormat formats a spinitron show name into the
// appropriate format for a spinitron show link
func spinitronNameFormat(input string) string {
	return removeNonWordCharacters(
		strings.ReplaceAll(
			strings.ReplaceAll(
				input, " ", "-"),
			"'", "-"))
}

func removeNonWordCharacters(input string) string {
	filter := func(r rune) rune {
		if syntax.IsWordChar(r) || r == '-' {
			return r
		}
		return -1
	}
	return strings.Map(filter, input)
}