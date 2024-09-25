package exchange

import (
	"context"

	"github.com/shahinrahimi/teletradebot/models"
)

type Exchange interface {
	// fundamental operations
	CheckSymbol(symbol string) bool
	FetchInterpreter(ctx context.Context, t *models.Trade) (*models.Interpreter, error)
	// basic operations
	PlaceStopOrder(ctx context.Context, po interface{}) (interface{}, error)
	PlaceTakeProfitOrder(ctx context.Context, po interface{}) (interface{}, error)
	CancelOrder(ctx context.Context, po interface{}) (interface{}, error)
	CloseOrder(ctx context.Context, po interface{}) (interface{}, error)
	GetOrder(ctx context.Context, po interface{}) (interface{}, error)
}
