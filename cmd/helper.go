package cmd

import (
	"time"
)

func unixToDateString(timestamp int64) string {
	var t time.Time

	// If the timestamp has more than 10 digits, it's in milliseconds
	if timestamp > 1e11 {
		t = time.UnixMilli(timestamp)
	} else {
		t = time.Unix(timestamp, 0)
	}

	// Use Go's standard layout reference date: Mon Jan 2 15:04:05 MST 2006
	return t.Format("2006-01-02 15:04:05")
}