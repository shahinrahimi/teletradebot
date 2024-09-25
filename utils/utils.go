package utils

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/swagger"
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

func ExtractOrderIDStr(orderIDorOrderRes interface{}) string {
	var orderIDStr string
	if orderID, ok := orderIDorOrderRes.(string); ok {
		orderIDStr = orderID
	} else if orderID, ok := orderIDorOrderRes.(int64); ok {
		orderIDStr = ConvertBinanceOrderID(orderID)
	} else if order, ok := orderIDorOrderRes.(*futures.CreateOrderResponse); ok {
		orderIDStr = ConvertBinanceOrderID(order.OrderID)
	} else if order, ok := orderIDorOrderRes.(*swagger.Order); ok {
		orderIDStr = order.OrderID
	} else {
		log.Panicf("unexpected error happened in casting order response to *futures.CreateOrderResponse or *swagger.Order: %T", orderIDorOrderRes)
	}
	return orderIDStr
}

func ExtractOrderStatus(orderRes interface{}) string {
	if order, ok := (orderRes).(*futures.Order); ok {
		return string(order.Status)
	} else if order, ok := (orderRes).(*swagger.Order); ok {
		return string(order.OrdStatus)
	} else {
		log.Panicf("unexpected error happened in casting error to futures.Order: %T", orderRes)
	}
	return ""
}
