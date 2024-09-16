package config

import "time"

var UseBinanceTestnet bool = true
var UseBitmexTestnet bool = true

var Shortcuts = map[string]string{
	"1":  "b BNBUSDT long 5m 0 1 105 105 1",
	"2":  "b BNBUSDT short 5m 0 1 105 105 1",
	"3":  "b SOLUSDT long 5m 0 1 105 105 1",
	"4":  "b SOLUSDT short 5m 0 1 105 105 1",
	"5":  "b AVAXUSDT long 5m 0 1 105 105 1",
	"6":  "b AVAXUSDT short 5m 0 1 105 105 1",
	"7":  "b ADAUSDT long 5m 0 1 105 105 1",
	"8":  "b ADAUSDT short 5m 0 1 105 105 1",
	"9":  "b XRPUSDT long 5m 0 1 105 105 1",
	"10": "b XRPUSDT short 5m 0 1 105 105 1",
	"11": "b BTCUSDT long 5m 0 1 105 105 1",
	"12": "b BTCUSDT short 5m 0 1 105 105 1",
	"13": "b ETHUSDT long 5m 0 1 105 105 1",
	"14": "b ETHUSDT short 5m 0 1 105 105 1",
	"15": "b LTCUSDT long 5m 0 1 105 105 1",
	"16": "b LTCUSDT short 5m 0 1 105 105 1",
	"17": "b TRXUSDT long 5m 0 1 105 105 1",
	"18": "b TRXUSDT short 5m 0 1 105 105 1",
	"19": "b BATUSDT long 5m 0 1 105 105 1",
	"20": "b BATUSDT short 5m 0 1 105 105 1",
	"21": "m BNBUSDT long 5m 0 1 105 105 1",
	"22": "m BNBUSDT short 5m 0 1 105 105 1",
	"23": "m SOLUSDT long 5m 0 1 105 105 1",
	"24": "m SOLUSDT short 5m 0 1 105 105 1",
	"25": "m AVAXUSDT long 5m 0 1 105 105 1",
	"26": "m AVAXUSDT short 5m 0 1 105 105 1",
	"27": "m ADAUSDT long 5m 0 1 105 105 1",
	"28": "m ADAUSDT short 5m 0 1 105 105 1",
	"29": "m XRPUSDT long 5m 0 1 105 105 1",
	"30": "m XRPUSDT short 5m 0 1 105 105 1",
	"31": "m XBTUSDT long 5m 0 1 105 105 1",
	"32": "m XBTUSDT short 5m 0 1 105 105 1",
	"33": "m ETHUSDT long 5m 0 1 105 105 1",
	"34": "m ETHUSDT short 5m 0 1 105 105 1",
	"35": "m LTCUSDT long 5m 0 1 105 105 1",
	"36": "m LTCUSDT short 5m 0 1 105 105 1",
	"37": "m TRXUSDT long 5m 0 1 105 105 1",
	"38": "m TRXUSDT short 5m 0 1 105 105 1",
	"39": "m LINKUSDT long 5m 0 1 105 105 1",
	"40": "m LINKUSDT short 5m 0 1 105 105 1",
	"s":  "b BTCUSDT short 15m 0.1 1 50 50 1",
	"l":  "b BTCUSDT long 15m 0.1 1 50 50 1",
}

var MaxTries int = 3
var WaitForNextTries time.Duration = time.Second * 3 // seconds

var UserIDs = []int64{
	104196468,
	539168576,
}
