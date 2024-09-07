package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
