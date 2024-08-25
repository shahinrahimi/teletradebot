package bot

import (
	"context"

	"gihub.com/shahinrahimi/teletradebot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleHelp(u *tgbotapi.Update, ctx context.Context) error {
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

func (b *Bot) HandleRemove(u *tgbotapi.Update, ctx context.Context) error {
	return nil
}

func (b *Bot) HandleExecute(u *tgbotapi.Update, ctx context.Context) error {
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
