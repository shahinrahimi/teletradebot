package config

import "time"

var UseBinanceTestnet bool = true
var UseBitmexTestnet bool = true

var Shortcuts = map[string]string{
	"1":  "b BNBUSDT long 5m 0.1 1 105 105 1",
	"2":  "b BNBUSDT short 5m 0.1 1 105 105 1",
	"3":  "b SOLUSDT long 5m 0.1 1 105 105 1",
	"4":  "b SOLUSDT short 5m 0.1 1 105 105 1",
	"5":  "b AVAXUSDT long 5m 0.1 1 105 105 1",
	"6":  "b AVAXUSDT short 5m 0.1 1 105 105 1",
	"7":  "b ADAUSDT long 5m 0.1 1 105 105 1",
	"8":  "b ADAUSDT short 5m 0.1 1 105 105 1",
	"9":  "b XRPUSDT long 5m 0.1 1 105 105 1",
	"10": "b XRPUSDT short 5m 0.1 1 105 105 1",
	"11": "b BTCUSDT long 5m 0.1 1 105 105 1",
	"12": "b BTCUSDT short 5m 0.1 1 105 105 1",
	"13": "b ETHUSDT long 5m 0.1 1 105 105 1",
	"14": "b ETHUSDT short 5m 0.1 1 105 105 1",
	"15": "b LTCUSDT long 5m 0.1 1 105 105 1",
	"16": "b LTCUSDT short 5m 0.1 1 105 105 1",
	"17": "b TRXUSDT long 5m 0.1 1 105 105 1",
	"18": "b TRXUSDT short 5m 0.1 1 105 105 1",
	"19": "b BATUSDT long 5m 0.1 1 105 105 1",
	"20": "b BATUSDT short 5m 0.1 1 105 105 1",
	"s":  "b BTCUSDT short 15m 0.1 1 50 50 1",
	"l":  "b BTCUSDT long 15m 0.1 1 50 50 1",
}

var MaxTries int = 3
var WaitForNextTries time.Duration = time.Second * 3 // seconds

var UserIDs = []int64{
	104196468,
	539168576,
}
