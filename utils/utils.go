package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func ConvertBinanceOrderID(orderID int64) string {
	return strconv.FormatInt(orderID, 10)
}

func ConvertOrderIDtoBinanceOrderID(orderID string) (int64, error) {
	return strconv.ParseInt(orderID, 10, 64)
}

func ConvertTime(timestamp int64) time.Time {
	// Convert milliseconds to seconds and nanoseconds
	seconds := timestamp / 1000
	nanoseconds := (timestamp % 1000) * 1e6
	t := time.Unix(seconds, nanoseconds).UTC()
	return t
}

func FormatTimestamp(timestamp int64) string {
	t := time.Unix(0, timestamp*int64(time.Millisecond))

	formattedTime := t.Format("2006-01-02 15:04:05")

	return formattedTime
}

func FriendlyDuration(duration time.Duration) string {
	// Convert to hours, minutes, seconds, etc.
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	milliseconds := int(duration.Milliseconds()) % 1000

	// Build a friendly string representation
	var friendlyDuration string
	if hours > 0 {
		friendlyDuration += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 {
		friendlyDuration += fmt.Sprintf("%dm ", minutes)
	}
	if seconds > 0 {
		friendlyDuration += fmt.Sprintf("%ds ", seconds)
	}
	if milliseconds > 0 {
		friendlyDuration += fmt.Sprintf("%dms", milliseconds)
	}

	return friendlyDuration
}

func PrintStructFields(s interface{}) {
	val := reflect.ValueOf(s)
	// check if the passed value is a pointer
	if val.Kind() == reflect.Ptr {
		// Dereference the pointer to get the underlying struct
		val = val.Elem()
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		// Get the field name
		fieldName := typ.Field(i).Name

		// Get the field type
		fieldType := typ.Field(i).Type

		// Get the field value
		fieldValue := val.Field(i)

		fmt.Printf("%s (%s) = %v\n", fieldName, fieldType, fieldValue)
	}
}
