package bitmex

import (
	"context"
	"log"

	swagger "github.com/shahinrahimi/teletradebot/swagger"
)

type BitmexClient struct {
	l       *log.Logger
	client  *swagger.APIClient
	auth    context.Context
	Verbose bool
	ApiKey  string
	ApiSec  string
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
		l:       l,
		client:  client,
		auth:    auth,
		ApiKey:  apiKey,
		ApiSec:  apiSec,
		Verbose: true,
	}
}

func (mc *BitmexClient) GetAuthContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, swagger.ContextAPIKey, swagger.APIKey{
		Key:    mc.ApiKey,
		Secret: mc.ApiSec,
	})
}
