package bot

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gihub.com/shahinrahimi/teletradebot/config"
	"gihub.com/shahinrahimi/teletradebot/models"
	"gihub.com/shahinrahimi/teletradebot/types"
	"gihub.com/shahinrahimi/teletradebot/utils"
	"github.com/adshao/go-binance/v2/common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleHelp(u *tgbotapi.Update, ctx context.Context) error {
	var userID int64 = u.Message.From.ID
	b.SendMessage(userID, GetCommandHelp())
	return nil
}

func (b *Bot) HandleInfo(u *tgbotapi.Update, ctx context.Context) error {
	var userID int64 = u.Message.From.ID
	var username string = u.Message.From.UserName
	msg := fmt.Sprintf("UserID:\t%d\nUsername:\t%s", userID, username)
	b.SendMessage(userID, msg)
	return nil
}

func (b *Bot) HandleView(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	b.SendMessage(u.Message.From.ID, t.ToViewString())
	return nil
}

func (b *Bot) HandleAdd(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	if err := b.s.CreateTrade(&t); err != nil {
		b.l.Printf("error creating a new trade: %v", err)
		b.SendMessage(u.Message.From.ID, "Internal error while creating a new trade.")
		return err
	}
	b.SendMessage(u.Message.From.ID, "Trade created successfully!")
	return nil
}

func (b *Bot) HandleList(u *tgbotapi.Update, ctx context.Context) error {
	os, err := b.s.GetTrades()
	if err != nil {
		b.l.Printf("error getting trades: %v", err)
		b.SendMessage(u.Message.From.ID, "Internal error while listing trades.")
		return err
	}
	msg := ""
	for _, o := range os {
		msg = msg + o.ToListString() + "\n"
	}
	if len(os) == 0 {
		b.SendMessage(u.Message.From.ID, "No trades found.")
		return nil
	}
	b.SendMessage(u.Message.From.ID, "list of trades\n"+msg)
	return nil
}

func (b *Bot) HandleAlias(u *tgbotapi.Update, ctx context.Context) error {
	var msg string = "aliases: \n"
	for key, value := range config.Shortcuts {
		msg = msg + "'" + key + "' => " + value + "\n"
	}
	b.SendMessage(u.Message.From.ID, msg)
	return nil
}

func (b *Bot) HandleRemove(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	if t.State != types.STATE_IDLE {
		b.SendMessage(u.Message.From.ID, "The trade could not be removed to it has not state of Idle")
		return nil
	}
	if err := b.s.DeleteTrade(t.ID); err != nil {
		return err
	}
	b.SendMessage(u.Message.From.ID, "Trade removed successfully!")
	return nil
}

func (b *Bot) HandleDescribe(u *tgbotapi.Update, ctx context.Context) error {
	// o := ctx.Value(models.KeyTrade{}).(models.Trade)
	// if _, err := b.bc.GetKline(&o); err != nil {
	// 	return err
	// }
	//b.bc.TrackOrder()
	return nil
}

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
	if _, err := b.bc.CancelOrder(orderID, t.Pair); err != nil {
		if apiErr, ok := (err).(*common.APIError); ok {
			msg := fmt.Sprintf("Binance API:\ncould not cancel a order\ncode:%d\nmessage: %s", apiErr.Code, apiErr.Message)
			b.l.Println(msg)
			b.SendMessage(t.UserID, msg)
			return nil
		}
		b.l.Printf("error canceling order_id: %d", orderID)
		return nil
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

func (b *Bot) HandleCheck(u *tgbotapi.Update, ctx context.Context) error {

	return nil
}

func (b *Bot) HandleExecute2(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	if t.State != types.STATE_IDLE {
		b.SendMessage(u.Message.From.ID, "The trade could not be executed as it has already been executed once.")
		return nil
	}
	return b.mc.TryPlaceOrderForTrade(&t)
}

func (b *Bot) HandleExecute(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	if t.State != types.STATE_IDLE {
		b.SendMessage(u.Message.From.ID, "The trade could not be executed as it has already been executed once.")
		return nil
	}

	// prepared trade for order

	po, err := b.bc.PrepareTradeForOrder(&t)
	if err != nil {
		b.l.Printf("trade could not be executed, error in preparing state: %v", err)
		return nil
	}

	b.l.Printf("Placing %s order with quantity %s and stop price %s expires in: %s", po.Side, po.Quantity, po.StopPrice, utils.FriendlyDuration(po.Expiration))
	res, err := b.bc.PlacePreparedOrder(po)
	if err != nil {
		utils.PrintStructFields(err)
		fmt.Printf("Type of err: %T\n", err)
		if apiErr, ok := err.(*common.APIError); ok {
			msg := fmt.Sprintf("Binance API:\ncould not place a order for trade\ncode:%d\nmessage: %s", apiErr.Code, apiErr.Message)
			b.l.Println(msg)
			b.SendMessage(t.UserID, msg)
		}
		return err
	}
	msg := fmt.Sprintf("Order placed successfully for trade: %d, Order ID: %s", t.ID, t.OrderID)
	b.SendMessage(u.Message.From.ID, msg)
	// schedule order cancellation (it will raise error if currently filled)
	// if cancel successfully it will change trade state to replacing
	go b.scheduleOrderCancellation(res.OrderID, res.Symbol, po.Expiration, &t)

	// update trade state
	t.OrderID = strconv.FormatInt(res.OrderID, 10)
	t.State = types.STATE_PLACED
	t.UpdatedAt = time.Now().UTC()
	if err := b.s.UpdateTrade(&t); err != nil {
		msg := fmt.Sprintf("An important error occurred. The trade with ID '%d' could not be updated, which might cause tracking issues. Order ID: %s", t.ID, t.OrderID)
		b.SendMessage(u.Message.From.ID, msg)
		return err
	}
	return nil
}

func (b *Bot) MakeHandlerBotFunc(f ErrorHandler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		if err := f(u, ctx); err != nil {
			b.l.Printf("we have error %v", err)
		}
	}
}
