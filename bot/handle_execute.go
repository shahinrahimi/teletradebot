package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) HandleExecute2(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	if t.State != types.STATE_IDLE {
		b.SendMessage(u.Message.From.ID, "The trade could not be executed as it has already been executed once.")
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
	b.SendMessage(u.Message.From.ID, msg)

	return nil
}

func (b *Bot) HandleExecute(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	if t.State != types.STATE_IDLE {
		b.SendMessage(u.Message.From.ID, "The trade could not be executed as it has already been executed once.")
		return nil
	}
	var orderID string
	if t.Account == types.ACCOUNT_B {
		// prepared trade for order
		po, err := b.bc.PrepareOrder(ctx, &t)
		if err != nil {
			b.l.Printf("trade could not be executed, error in preparing state: %v", err)
			return nil
		}

		b.l.Printf("Placing %s order with quantity %s and stop price %s expires in: %s", po.Side, po.Quantity, po.StopPrice, utils.FriendlyDuration(po.Expiration))
		res, err := b.bc.PlacePreparedOrder(ctx, po)
		if err != nil {
			b.handleAPIError(err, t.UserID)
			return err
		}

		// schedule order cancellation (it will raise error if currently filled)
		// if cancel successfully it will change trade state to replacing
		go b.scheduleOrderReplacement(ctx, po.Expiration, res.OrderID, &t)

		orderID = utils.ConvertBinanceOrderID(res.OrderID)

	} else {

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
		//TODO schedule order cancellation

		orderID = order.OrderID
	}

	msg := fmt.Sprintf("Order placed successfully\n\nOrder ID: %s\nTrade ID: %d", orderID, t.ID)

	b.SendMessage(u.Message.From.ID, msg)

	// update trade state
	b.l.Printf("try to update trade DB: %s", orderID)
	if err := b.s.UpdateTradePlaced(&t, orderID); err != nil {
		b.l.Printf("error updating trade to DB")
		return err
	}
	return nil
}
