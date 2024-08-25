package exchange

import "log"

type BitmexClient struct {
	l *log.Logger
}

func NewBitmexClient(l *log.Logger) *BitmexClient {
	return &BitmexClient{l}
}
