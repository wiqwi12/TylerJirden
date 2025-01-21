package internal

import (
	"fmt"
	"time"
)

func FormatDate(time time.Time) interface{} {
	return fmt.Sprintf("%02d:%02d:%02d", time.Day(), time.Month(), time.Year())
}

func FormatTime(time time.Time) interface{} {
	return fmt.Sprintf("%02d:%02d:%02d", time.Hour(), time.Minute(), time.Second())
}
