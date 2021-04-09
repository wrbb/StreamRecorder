package util

import (
	"log"
	"time"
)


// TimeLoc is the time zone the stream recording computer is in
var TimeLoc *time.Location = MustLoadLocation("EST5EDT")

// GetMidnight returns a time.Time of 00:00 the next day
func GetMidnight() time.Time {
	year, month, day := time.Now().Add(time.Hour * 24).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

// TimeUntilMidnight returns a time.Duration of the time
// until midnight (defined as 00:00 the next day)
func TimeUntilMidnight() time.Duration {
	midnight := GetMidnight()
	return midnight.Sub(time.Now())
}


// MustLoadLocation loads a time.Location from a string
// and will exit with code 1 if the timezone is not loaded
func MustLoadLocation(location string) *time.Location {
	// Load the time
	loc, err := time.LoadLocation(location)
	if err != nil {
		log.Fatal("Unable to load timezone")
	}
	return loc
}
