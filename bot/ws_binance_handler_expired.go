package bot

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) handleExpired(ctx context.Context, t *models.Trade, f futures.WsOrderTradeUpdate) {
	orderID := utils.ConvertBinanceOrderID(f.ID)
	t, orderIDType := b.c.GetTradeByAnyOrderID(orderID)
	switch orderIDType {
	case types.OrderIDTypeMain:
		b.l.Printf("the orderID is expired with trade ID: %s", orderID)
		b.handleNewExpired(t, f)
	case types.OrderIDTypeTakeProfit:
		b.l.Printf("the take-profit orderID is expired with trade ID: %s", orderID)
		//b.handleTPFilled(ctx, t, f)
	case types.OrderIDTypeStopLoss:
		b.l.Printf("the stop-loss orderID is expired with trade ID: %s", orderID)
		//b.handleSLFilled(ctx, t, f)
	default:
		b.l.Printf("the orderID expired is not associate with any trade: %s", orderID)
	}
}

func (b *Bot) handleNewExpired(t *models.Trade, f futures.WsOrderTradeUpdate) {

	// update trade state
	b.c.UpdateTradeExpired(t.ID)

	orderID := utils.ConvertBinanceOrderID(f.ID)

	// message the user
	msg := fmt.Sprintf("Order expired.\n\nOrder ID: %s\nTrade ID: %d", orderID, t.ID)
	b.MsgChan <- types.BotMessage{
		ChatID: t.UserID,
		MsgStr: msg,
	}
}
