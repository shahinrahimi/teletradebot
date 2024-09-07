package bot

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleCancel(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	if t.State == types.STATE_IDLE {
		b.SendMessage(u.Message.From.ID, "The trade already has state of idle, the trade not have not any order_id associate with to cancel")
		return nil
	}
	if t.State != types.STATE_PLACED {
		b.SendMessage(u.Message.From.ID, "The trade could not be canceled as it has filled order.")
		return nil
	}

	if t.OrderID == "" {
		b.l.Printf("error the trade has state of placed but does not have any order_id associate with")
		b.SendMessage(u.Message.From.ID, "Internal error")
		return nil
	}
	orderID, err := strconv.ParseInt(t.OrderID, 10, 64)
	if err != nil {
		b.l.Printf("error converting order_id string it int 64: %v", err)
		b.SendMessage(u.Message.From.ID, "Internal error")
		return nil
	}
	if _, err := b.bc.CancelOrder(ctx, orderID, t.Symbol); err != nil {
		b.handleAPIError(err, t.UserID)
		return err
	}
	msg := fmt.Sprintf("Placed order successfully canceled, trade: %d, Order ID: %s", t.ID, t.OrderID)
	b.SendMessage(u.Message.From.ID, msg)

	t.OrderID = ""
	t.State = types.STATE_IDLE
	if err := b.s.UpdateTrade(&t); err != nil {
		b.l.Printf("error updating trade: %v", err)
		return err
	}
	return nil
}
