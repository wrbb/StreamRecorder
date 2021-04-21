package util

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)


// InitLoggers initializes the loggers for the application
// There are 3 loggers one for each logging level:
// InfoLogger, WarningLogger and ErrorLogger
func InitLoggers() {
	// make log directory
	err := os.MkdirAll(viper.GetString(LogDirectory), 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Create the log file for warnings and info messages
	logFile, err := os.OpenFile(path.Join(viper.GetString(LogDirectory), "main.log" ), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Creates the log file for error message
	errorFile, err := os.OpenFile(path.Join(viper.GetString(LogDirectory) ,"error.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(logFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// ErrorLog logs a message to the ErrorLogger in addition to sending a
// message using the SlackClient
func ErrorLog(msg string) {
	err := SlackClient.SendMessage(fmt.Sprintf("Vortex Error: %s", msg))
	if err != nil {
		ErrorLogger.Print(err.Error())
	}
	ErrorLogger.Print(msg)
}