package internal

import (
	"fmt"
	"time"
)

func FormatDate(time time.Time) string {

	return fmt.Sprintf("%02d:%02d:%02d", time.Day(), time.Month(), time.Year())

}

func FormatTime(time time.Time) string {
	return fmt.Sprintf("%02d:%02d:%02d", time.Hour(), time.Minute(), time.Second())
}
