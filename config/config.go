package config

import "time"

const (
	UseBinanceTestnet bool = true
	UseBitmexTestnet  bool = true
)

var Shortcuts = map[string]string{
	"1":  "m XBTUSDT long 15m 0 1 105 105 1",
	"2":  "m XBTUSDT short 15m 0 1 105 105 1",
	"3":  "m XBTUSDT long 1h 0 1 105 105 1",
	"4":  "m XBTUSDT short 1h 0 1 105 105 1",
	"5":  "m XBTUSDT long 15h 0 1 105 105 0",
	"6":  "m XBTUSDT short 15h 0 1 105 105 0",
	"7":  "m XBTUSDT long 15h 0 1 105 105 2",
	"8":  "m XBTUSDT short 15h 0 1 105 105 2",
	"10": "b BTCUSDT long 15m 0 1 105 105 1",
	"11": "b BTCUSDT short 15m 0 1 105 105 1",
	"12": "b BTCUSDT long 1h 0 1 105 105 1",
	"13": "b BTCUSDT short 1h 0 1 105 105 1",
	"14": "b BTCUSDT long 15h 0 1 105 105 0",
	"15": "b BTCUSDT short 15h 0 1 105 105 0",
	"16": "b BTCUSDT long 15h 0 1 105 105 2",
	"17": "b BTCUSDT short 15h 0 1 105 105 2",
	"s":  "b BTCUSDT short 15m 0.1 1 50 50 1",
	"l":  "b BTCUSDT long 15m 0.1 1 50 50 1",
}

var MaxTries int = 3
var WaitForNextTries time.Duration = time.Second * 3   // 3 seconds
var WaitForReplacement time.Duration = time.Second * 2 // 2 seconds

var UserIDs = []int64{
	104196468,
	539168576,
}

// hard coded contract size for each symbol for bitmex
var ContractSizes = map[string]float64{
	"XBTUSDT":  0.000001,
	"ETHUSDT":  0.00001,
	"ADAUSDT":  0.01,
	"AVAXUSDT": 0.0001,
	"BNBUSDT":  0.0001,
	"LTCUSDT":  0.0001,
	"LINKUSDT": 0.001,
	"SOLUSDT":  0.0001,
	"TRXUSDT":  0.1,
	"XRPUSDT":  0.01,
}
