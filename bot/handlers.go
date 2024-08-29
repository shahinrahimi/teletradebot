package bot

import (
	"context"

	"gihub.com/shahinrahimi/teletradebot/models"
	"gihub.com/shahinrahimi/teletradebot/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleHelp(u *tgbotapi.Update, ctx context.Context) error {
	b.SendMessage(u.Message.From.ID, GetCommandHelp())
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
		b.SendMessage(u.Message.From.ID, "internal error creating a new trade")
		return err
	}
	b.SendMessage(u.Message.From.ID, "Successfully trade created!")
	return nil
}

func (b *Bot) HandleList(u *tgbotapi.Update, ctx context.Context) error {
	os, err := b.s.GetTrades()
	if err != nil {
		b.l.Printf("error getting trades: %v", err)
		b.SendMessage(u.Message.From.ID, "internal error listing trades")
		return err
	}
	msg := ""
	for _, o := range os {
		msg = msg + o.ToListString() + "\n"
	}
	if len(os) == 0 {
		b.SendMessage(u.Message.From.ID, "There is no trade found")
		return nil
	}
	b.SendMessage(u.Message.From.ID, "list of trades\n"+msg)
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
	b.bc.TrackOrder()
	return nil
}

func (b *Bot) HandleExecute(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	if t.State != types.STATE_IDLE {
		b.SendMessage(u.Message.From.ID, "the trade could not be executed since it is executed once")
		return nil
	}
	if err := b.bc.PlaceOrder(&t); err != nil {
		b.l.Printf("error in placing trade: %v", err)
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
