package server

import (
	"html/template"
	"net/http"
	"regexp/syntax"
	"strings"
	"wrbb-stream-recorder/internal/recording"
	"wrbb-stream-recorder/internal/util"
)


const DashboardTemplate = "web/template/dashboard.html"
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
	tmpl := template.Must(template.New("dashboard").Funcs(template.FuncMap{
		"nameFormat":  spinitronNameFormat,
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
	return strings.TrimFunc(
		strings.ReplaceAll(
			strings.ReplaceAll(input, " ", "-"),
			"'", "-"),
			func(r rune) bool {
				return !syntax.IsWordChar(r)
		})
}
