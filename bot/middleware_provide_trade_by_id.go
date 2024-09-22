package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) ProvideTradeByID(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		userID := u.Message.From.ID
		args := strings.Split(u.Message.CommandArguments(), " ")
		id, err := strconv.Atoi(args[0])
		if err != nil {
			msg := fmt.Sprintf("Invalid trade ID: '%s'. Please provide a numeric ID.", args[0])
			b.MsgChan <- types.BotMessage{
				ChatID: userID,
				MsgStr: msg,
			}
			return
		}
		ts := b.c.GetTrades()
		for _, t := range ts {
			if t.ID == int64(id) {
				ctx = context.WithValue(ctx, models.KeyTrade{}, *t)
				next(u, ctx)
				return
			}
		}
		msg := fmt.Sprintf("No trade found with ID: %d.", id)
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
	}
}
