package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleExecute2(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	userID := u.Message.From.ID
	if t.State != types.STATE_IDLE {
		msg := "The trade could not be executed as it has already been executed once."
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
		return nil
	}
	po, err := b.mc.PrepareOrder(ctx, &t)
	if err != nil {
		b.l.Printf("error preparing order: %v", err)
		return err
	}
	order, err := b.mc.PlacePreparedOrder(po)
	if err != nil {
		b.l.Printf("error placing order: %v", err)
		return err
	}

	msg := fmt.Sprintf("Order placed successfully\n\nOrder ID: %s\nTrade ID: %d", order.OrderID, t.ID)
	b.MsgChan <- types.BotMessage{
		ChatID: userID,
		MsgStr: msg,
	}

	return nil
}

func (b *Bot) HandleExecute(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	userID := u.Message.From.ID
	if t.State != types.STATE_IDLE {
		msg := "The trade could not be executed as it has already been executed once."
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
		return nil
	}
	if t.Account == types.ACCOUNT_B {
		go b.bc.ExecuteTrade(ctx, &t, false)
	} else {

		return b.HandleUnderDevelopment(u, ctx)

		// po, err := b.mc.PrepareOrder(ctx, &t)
		// if err != nil {
		// 	b.l.Printf("error preparing order: %v", err)
		// 	return err
		// }
		// order, err := b.mc.PlacePreparedOrder(po)
		// if err != nil {
		// 	b.l.Printf("error placing order: %v", err)
		// 	return err
		// }
		//TODO schedule order cancellation

		//orderID = order.OrderID
	}
	return nil
}
