package bitmexc

import (
	"context"
	"fmt"
	"log"

	"github.com/shahinrahimi/teletradebot/swagger"
)

type BitmexClient struct {
	l                 *log.Logger
	client            *swagger.APIClient
	auth              context.Context
	activeInstruments []swagger.Instrument
	Verbose           bool
	ApiKey            string
	ApiSec            string
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

func (mc *BitmexClient) GetSymbol(symbol string) (*swagger.Instrument, error) {
	if mc.activeInstruments == nil {
		return nil, fmt.Errorf("exchange info not available right now")
	}
	for _, s := range mc.activeInstruments {
		if s.Symbol == symbol {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("symbol %s not found", symbol)
}

func (mc *BitmexClient) getAuthContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, swagger.ContextAPIKey, swagger.APIKey{
		Key:    mc.ApiKey,
		Secret: mc.ApiSec,
	})
}
