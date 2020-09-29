package recording

import (
	"StreamRecorderGo/pkg/spinitron"
	"time"
)

func ScheduleLoop(schedule *spinitron.ShowSchedule) {
	for {
		if len(*schedule) > 0{
			next_show := (*schedule)[0]
			if next_show.Start.Before(time.Now()) {
				if next_show.End.After(time.Now()) {
					schedule = &(*schedule)[1:]
				}

			}

		}

	}
}

func RecordShow(show spinitron.Show) {

}