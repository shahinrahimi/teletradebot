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
	userID := u.Message.From.ID
	if t.State != types.STATE_PLACED {
		msg := fmt.Sprintf("The trade is currently in an invalid state [%s] for cancellation.", t.State)
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
		return nil
	}
	// should not happened
	if t.OrderID == "" {
		b.l.Printf("Unexpected issue: the trade with a state of 'placed' is missing an OrderID")
		msg := "Unable to find the Order ID for the trade."
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
		return nil
	}
	if t.Account == types.ACCOUNT_B {
		orderID, err := utils.ConvertOrderIDtoBinanceOrderID(t.OrderID)
		if err != nil {
			b.l.Printf("Unexpected issue: the trade's OrderID is not in a valid format for conversion: %v", err)
			msg := "The Order ID for the trade is not in a valid format."
			b.MsgChan <- types.BotMessage{
				ChatID: userID,
				MsgStr: msg,
			}
			return err
		}
		if _, err := b.bc.CancelOrder(ctx, orderID, t.Symbol); err != nil {
			b.handleAPIError(err, t.UserID)
			return err
		}
	} else {
		// the bitmex logic
		return b.HandleUnderDevelopment(u, ctx)
	}

	msg := fmt.Sprintf("The order has been successfully canceled.\n\nOrder ID: %s\nTrade ID: %d\n", t.OrderID, t.ID)
	b.MsgChan <- types.BotMessage{
		ChatID: userID,
		MsgStr: msg,
	}

	if err := b.s.UpdateTradeCancelled(&t); err != nil {
		b.l.Printf("Error updating the trade status: %v", err)
		return err
	}
	return nil
}
