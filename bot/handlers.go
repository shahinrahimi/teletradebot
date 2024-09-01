package bot

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gihub.com/shahinrahimi/teletradebot/config"
	"gihub.com/shahinrahimi/teletradebot/models"
	"gihub.com/shahinrahimi/teletradebot/types"
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

func (b *Bot) HandleExecute(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	if t.State != types.STATE_IDLE {
		b.SendMessage(u.Message.From.ID, "The trade could not be executed as it has already been executed once.")
		return nil
	}
	res, err := b.bc.TryPlaceOrderForTrade(&t)
	if err != nil {
		b.l.Printf("error in placing trade: %v", err)
		return err
	}
	orderID := strconv.FormatInt(res.OrderID, 10)
	msg := fmt.Sprintf("Order placed successfully. Order ID: %s", orderID)
	b.SendMessage(u.Message.From.ID, msg)

	// update trade for order_id
	t.OrderID = orderID
	t.State = types.STATE_PLACED
	t.UpdatedAt = time.Now().UTC()
	if err := b.s.UpdateTrade(&t); err != nil {
		msg := fmt.Sprintf("An important error occurred. The trade with ID '%d' could not be updated, which might cause tracking issues. Order ID: %s", t.ID, t.OrderID)
		b.SendMessage(u.Message.From.ID, msg)
		return err
	}

	return nil
}

func (b *Bot) HandleCancel(u *tgbotapi.Update, ctx context.Context) error {
	return nil
}

func (b *Bot) HandleCheck(u *tgbotapi.Update, ctx context.Context) error {
	return nil
}

func (b *Bot) MakeHandlerBotFunc(f ErrorHandler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		if err := f(u, ctx); err != nil {
			b.l.Printf("we have error %v", err)
		}
	}
}
