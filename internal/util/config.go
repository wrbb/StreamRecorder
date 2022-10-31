package util

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

const (
	StorageLocation = "storage_location"
	StreamUrl       = "stream_url"
	WriteShows      = "debug.write_shows"
	LogDirectory    = "debug.log_directory"
	SpinitronAPIKey = "spinitron.api_key"
)

var (
	_, b, _, _ = runtime.Caller(0)
	root       = filepath.Join(filepath.Dir(b), "../..")
)

// InitConfig initializes the viper config for the application
func InitConfig() {
	// add config file location
	viper.AddConfigPath(root)
	// Set name of the file
	viper.SetConfigName("config")
	// Set file type
	viper.SetConfigType("yaml")
	// Set paths to search for config
	// Set defaults
	viper.SetDefault(StorageLocation, "./storage")
	viper.SetDefault(StreamUrl, "http://stream.radiojar.com/9950r946bzzuv")
	viper.SetDefault(WriteShows, true)
	viper.SetDefault(LogDirectory, "./logs")
	viper.SetDefault(SpinitronAPIKey, os.Getenv("SPINITRON_API_KEY"))

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
