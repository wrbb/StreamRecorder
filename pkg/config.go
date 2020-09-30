package pkg

import "os"

type Config struct {
	StreamURL string
	StorageLocation string
}


func GetConfig() Config {
	return Config{
		StorageLocation: GetEnv("VORTEX_STORAGE_LOCATION", "."),
		StreamURL: GetEnv("WRBB_STREAM_URL", "http://stream.radiojar.com/9950r946bzzuv"),
	}
}

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}