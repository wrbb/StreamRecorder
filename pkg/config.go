package pkg

import "os"

type Config struct {
	StreamURL string
	StorageLocation string
	TurnOffWrite bool
}


func GetConfig() Config {
	return Config{
		StorageLocation: GetEnv("VORTEX_STORAGE_LOCATION", "."),
		StreamURL: GetEnv("WRBB_STREAM_URL", "http://stream.radiojar.com/9950r946bzzuv"),
		TurnOffWrite: GetEnv("WRITE_SHOWS", "1") == "0",
	}
}

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}