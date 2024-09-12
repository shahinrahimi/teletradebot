package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleList(u *tgbotapi.Update, ctx context.Context) error {
	userID := u.Message.From.ID
	ts, err := b.s.GetTrades()
	if err != nil {
		b.l.Printf("error getting trades: %v", err)
		return err
	}
	msg := ""
	for _, t := range ts {
		msg = msg + t.ToListString() + "\n"
	}
	if len(ts) == 0 {
		b.MsgChan <- BotMessage{
			ChatID: userID,
			MsgStr: "No trades found.",
		}
		return nil
	}
	b.MsgChan <- BotMessage{
		ChatID: userID,
		MsgStr: "List of trades\n" + msg,
	}
	return nil
}
