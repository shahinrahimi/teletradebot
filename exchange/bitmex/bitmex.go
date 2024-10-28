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
	apiKey            string
	apiSec            string
	DbgChan           chan string
}

func NewBitmexClient(l *log.Logger, apiKey string, apiSec string, UseTestnet bool, dbgChan chan string) *BitmexClient {
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
		DbgChan: dbgChan,
	}
}

func (mc *BitmexClient) CheckSymbol(ctx context.Context, symbol string) (bool, error) {
	if _, err := mc.GetSymbol(ctx, symbol); err != nil {
		return false, err
	}
	return true, nil
}

func (mc *BitmexClient) GetSymbol(ctx context.Context, symbol string) (*swagger.Instrument, error) {
	// fallback
	if mc.activeInstruments == nil {
		mc.l.Printf("activeInstruments is nil trying to fetch")
		instruments, _, err := mc.client.InstrumentApi.InstrumentGetActive(ctx)
		if err != nil {
			mc.l.Printf("error fetching exchange info: %v", err)
			return nil, fmt.Errorf("error fetching exchange info: %v", err)
		}
		mc.activeInstruments = instruments
		for _, s := range mc.activeInstruments {
			if s.Symbol == symbol {
				return &s, nil
			}
		}
		return nil, fmt.Errorf("instrument %s not found", symbol)
	} else {
		for _, s := range mc.activeInstruments {
			if s.Symbol == symbol {
				return &s, nil
			}
		}
		return nil, fmt.Errorf("instrument %s not found", symbol)
	}
}

func (mc *BitmexClient) getAuthContext(ctx context.Context) context.Context {

	return context.WithValue(ctx, swagger.ContextAPIKey, swagger.APIKey{
		Key:    mc.apiKey,
		Secret: mc.apiSec,
	})
}
