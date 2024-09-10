package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) HandleCancel(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	if t.State != types.STATE_PLACED {
		msg := fmt.Sprintf("The trade is currently in an invalid state [%s] for cancellation.", t.State)
		b.SendMessage(u.Message.From.ID, msg)
		return nil
	}
	// should not happened
	if t.OrderID == "" {
		b.l.Printf("Unexpected issue: the trade with a state of 'placed' is missing an OrderID")
		b.SendMessage(u.Message.From.ID, "Unable to find the Order ID for the trade.")
		return nil
	}
	if t.Account == types.ACCOUNT_B {
		orderID, err := utils.ConvertOrderIDtoBinanceOrderID(t.OrderID)
		if err != nil {
			b.l.Printf("Unexpected issue: the trade's OrderID is not in a valid format for conversion: %v", err)
			b.SendMessage(u.Message.From.ID, "The Order ID for the trade is not in a valid format.")
			return err
		}
		if _, err := b.bc.CancelOrder(ctx, orderID, t.Symbol); err != nil {
			b.handleAPIError(err, t.UserID)
			return err
		}
	} else {
		// the bitmex logic
		b.HandleUnderDevelopment(u, ctx)
		return nil
	}

	msg := fmt.Sprintf("The order has been successfully canceled.\n\nOrder ID: %s\nTrade ID: %d\n", t.OrderID, t.ID)
	b.SendMessage(u.Message.From.ID, msg)

	if err := b.s.UpdateTradeCancelled(&t); err != nil {
		b.l.Printf("Error updating the trade status: %v", err)
		return err
	}
	return nil
}
