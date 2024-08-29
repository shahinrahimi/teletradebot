package utils

import "time"

func ConvertTime(timestamp int64) time.Time {
	// Convert milliseconds to seconds and nanoseconds
	seconds := timestamp / 1000
	nanoseconds := (timestamp % 1000) * 1e6

	// Convert to time.Time
	return time.Unix(seconds, nanoseconds)
}
