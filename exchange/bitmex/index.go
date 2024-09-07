package bitmex

import (
	"context"
	"log"

	swagger "gihub.com/shahinrahimi/teletradebot/go-client"
)

type BitmexClient struct {
	l           *log.Logger
	client      *swagger.APIClient
	auth        context.Context
	margin      *swagger.Margin
	instruments []swagger.Instrument
}

func NewBitmexClient(l *log.Logger, apiKey string, apiSec string, UseTestnet bool) *BitmexClient {
	cfg := swagger.NewConfiguration()
	if UseTestnet {
		cfg.BasePath = "https://testnet.bitmex.com/api/v1"
	}
	client := swagger.NewAPIClient(cfg)
	auth := context.WithValue(context.TODO(), swagger.ContextAPIKey, swagger.APIKey{
		Key:    apiKey,
		Secret: apiSec,
	})

	return &BitmexClient{
		l:      l,
		client: client,
		auth:   auth,
	}
}
