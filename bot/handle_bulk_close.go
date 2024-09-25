package bot

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleBulkClose(u *tgbotapi.Update, ctx context.Context) error {
	userID := u.Message.From.ID
	var trades []*models.Trade
	args := strings.Split(u.Message.CommandArguments(), " ")
	if len(args) == 2 {
		switch {
		case args[0] == "b" && args[1] == "s":
			trades = b.c.GetAllTrades(types.ExchangeBinance, types.SideShort)
		case args[0] == "b" && args[1] == "l":
			trades = b.c.GetAllTrades(types.ExchangeBinance, types.SideLong)
		case args[0] == "m" && args[1] == "s":
			trades = b.c.GetAllTrades(types.ExchangeBitmex, types.SideShort)
		case args[0] == "m" && args[1] == "l":
			trades = b.c.GetAllTrades(types.ExchangeBitmex, types.SideLong)
		default:
			msg := "Wrong arguments. Valid arguments are: b [s|l] and m [s|l]"
			b.MsgChan <- types.BotMessage{ChatID: userID, MsgStr: msg}
			return nil
		}

		for _, t := range trades {
			ctx = context.WithValue(ctx, models.KeyTrade{}, t)
			go b.HandleClose(u, ctx)
		}
	} else {
		trades := b.c.GetTrades()
		for _, t := range trades {
			ctx = context.WithValue(ctx, models.KeyTrade{}, *t)
			go b.HandleClose(u, ctx)
		}
	}
	return nil
}
