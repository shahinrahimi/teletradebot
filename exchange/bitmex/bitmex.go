package bitmex

import (
	"context"
	"fmt"
	"log"

	"github.com/shahinrahimi/teletradebot/swagger"
)

type BitmexClient struct {
	l                 *log.Logger
	client            *swagger.APIClient
	activeInstruments []swagger.Instrument
	Verbose           bool
	apiKey            string
	apiSec            string
}

func NewBitmexClient(l *log.Logger, apiKey string, apiSec string, UseTestnet bool) *BitmexClient {
	cfg := swagger.NewConfiguration()
	if UseTestnet {
		cfg.BasePath = "https://testnet.bitmex.com/api/v1"
	}
	client := swagger.NewAPIClient(cfg)
	return &BitmexClient{
		l:       l,
		client:  client,
		apiKey:  apiKey,
		apiSec:  apiSec,
		Verbose: true,
	}
}

func (mc *BitmexClient) CheckSymbol(symbol string) bool {
	if mc.activeInstruments == nil {
		mc.l.Printf("exchange info not available right now please try after some time")
		return false
	}
	for _, s := range mc.activeInstruments {
		if s.Symbol == symbol {
			return true
		}
	}
	return false
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
		Key:    mc.apiKey,
		Secret: mc.apiSec,
	})
}
