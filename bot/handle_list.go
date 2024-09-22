package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleList(u *tgbotapi.Update, ctx context.Context) error {
	userID := u.Message.From.ID
	ts := b.c.GetTrades()
	msg := ""
	for _, t := range ts {
		msg = msg + t.ToListString() + "\n"
	}
	if len(ts) == 0 {
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: "No trades found.",
		}
		return nil
	}
	b.MsgChan <- types.BotMessage{
		ChatID: userID,
		MsgStr: "List of trades\n" + msg,
	}
	return nil
}
