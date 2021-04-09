package server

import (
	"html/template"
	"net/http"
	"wrbb-stream-recorder/internal/recording"
	"wrbb-stream-recorder/internal/util"
)

// InitServer creates the http server
func InitServer() {
	http.HandleFunc("/", Dashboard)
	_ = http.ListenAndServe(":1049", nil)
}

// Dashboard in the handler function for the http server's
// index request
func Dashboard(w http.ResponseWriter, _ *http.Request) {
	currentShows := recording.GetCurrentShows()
	w.Header().Add("Content Type", "text/html")
	tmpl := template.Must(template.ParseFiles("web/template/dashboard.html"))
	err := tmpl.Execute(w, currentShows)
	if err != nil {
		w.WriteHeader(500)
		util.ErrorLogger.Printf("Unable to execute template: %s\n", err.Error())
	}
}