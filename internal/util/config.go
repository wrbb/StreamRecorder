package util

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

const (
	StorageLocation = "storage_location"
	StreamUrl       = "stream_url"
	WriteShows      = "debug.write_shows"
	LogDirectory    = "debug.log_directory"
	SpinitronAPIKey = "spinitron.api_key"
)

// InitConfig initializes the viper config for the application
func InitConfig() {
	// Set name of the file
	viper.SetConfigName("config")
	// Set file type
	viper.SetConfigType("yaml")
	// Set paths to search for config
	viper.AddConfigPath("$HOME/.config/stream_recorder/")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath(".")
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
