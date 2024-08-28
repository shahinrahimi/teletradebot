package bot

import (
	"context"
	"fmt"

	"gihub.com/shahinrahimi/teletradebot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleHelp(u *tgbotapi.Update, ctx context.Context) error {
	b.SendMessage(u.Message.From.ID, GetCommandHelp())
	return nil
}

func (b *Bot) HandleView(u *tgbotapi.Update, ctx context.Context) error {
	o := ctx.Value(models.KeyOrder{}).(models.Order)
	b.SendMessage(u.Message.From.ID, o.ToViewString())
	return nil
}

func (b *Bot) HandleAdd(u *tgbotapi.Update, ctx context.Context) error {
	o := ctx.Value(models.KeyOrder{}).(models.Order)
	if err := b.s.CreateOrder(&o); err != nil {
		b.l.Printf("error creating a new order: %v", err)
		b.SendMessage(u.Message.From.ID, "internal error creating a new order")
		return err
	}
	b.SendMessage(u.Message.From.ID, "Successfully order created!")
	return nil
}

func (b *Bot) HandleList(u *tgbotapi.Update, ctx context.Context) error {
	os, err := b.s.GetOrders()
	if err != nil {
		b.l.Printf("error getting orders: %v", err)
		b.SendMessage(u.Message.From.ID, "internal error listing orders")
		return err
	}
	fmt.Println(len(os))
	msg := ""
	for _, o := range os {
		msg = msg + o.ToListString() + "\n"
	}
	b.SendMessage(u.Message.From.ID, "list o orders\n"+msg)
	return nil
}

func (b *Bot) HandleRemove(u *tgbotapi.Update, ctx context.Context) error {
	return nil
}

func (b *Bot) HandleDescribe(u *tgbotapi.Update, ctx context.Context) error {
	o := ctx.Value(models.KeyOrder{}).(models.Order)
	if err := b.bc.GetKline(&o); err != nil {
		return err
	}
	return nil
}

func (b *Bot) HandleExecute(u *tgbotapi.Update, ctx context.Context) error {
	o := ctx.Value(models.KeyOrder{}).(models.Order)
	if o.State != models.STATE_IDLE {
		b.SendMessage(u.Message.From.ID, "the order could not be executed since it is executed once")
		return nil
	}
	if err := b.bc.PlaceOrder(&o); err != nil {
		b.l.Printf("error in placing order: %v", err)
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
