package util

import (
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
	err := os.MkdirAll(viper.GetString(LogDirectory), 0666)
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