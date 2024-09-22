package bot

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleBulkReset(u *tgbotapi.Update, ctx context.Context) error {
	userID := u.Message.From.ID

	var trades []*models.Trade
	args := strings.Split(u.Message.CommandArguments(), " ")

	if len(args) == 2 {
		switch {
		case args[0] == "b" && args[1] == "s":
			trades = b.c.GetAllTrades(types.ACCOUNT_B, types.SIDE_S)
		case args[0] == "b" && args[1] == "l":
			trades = b.c.GetAllTrades(types.ACCOUNT_B, types.SIDE_L)
		case args[0] == "m" && args[1] == "s":
			trades = b.c.GetAllTrades(types.ACCOUNT_M, types.SIDE_S)
		case args[0] == "m" && args[1] == "l":
			trades = b.c.GetAllTrades(types.ACCOUNT_M, types.SIDE_L)
		default:
			msg := "Wrong arguments. Valid arguments are: b [s|l] and m [s|l]"
			b.MsgChan <- types.BotMessage{ChatID: userID, MsgStr: msg}
			return nil
		}

		for _, t := range trades {
			ctx = context.WithValue(ctx, models.KeyTrade{}, t)
			go b.HandleReset(u, ctx)
		}
	} else {
		trades := b.c.GetTrades()
		for _, t := range trades {
			ctx = context.WithValue(ctx, models.KeyTrade{}, *t)
			go b.HandleReset(u, ctx)
		}
	}
	return nil
}
