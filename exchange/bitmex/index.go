package bitmex

import (
	"context"
	"log"

	"github.com/shahinrahimi/teletradebot/store"
	swagger "github.com/shahinrahimi/teletradebot/swagger"
)

type BitmexClient struct {
	l      *log.Logger
	s      store.Storage
	client *swagger.APIClient
	auth   context.Context
}

func NewBitmexClient(l *log.Logger, s store.Storage, apiKey string, apiSec string, UseTestnet bool) *BitmexClient {
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
		s:      s,
		client: client,
		auth:   auth,
	}
}
