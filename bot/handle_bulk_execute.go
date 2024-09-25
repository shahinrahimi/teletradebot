package bot

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleBulkExecute(u *tgbotapi.Update, ctx context.Context) error {
	args := strings.Split(u.Message.CommandArguments(), " ")
	userID := u.Message.From.ID
	if len(args) == 2 {
		var trades []*models.Trade
		switch {
		case args[0] == "b" && args[1] == "s":
			trades = b.c.GetAllUniqueTrades(types.ExchangeBinance, types.SideShort, types.StateIdle)
		case args[0] == "b" && args[1] == "l":
			trades = b.c.GetAllUniqueTrades(types.ExchangeBinance, types.SideLong, types.StateIdle)
		case args[0] == "m" && args[1] == "s":
			trades = b.c.GetAllUniqueTrades(types.ExchangeBitmex, types.SideShort, types.StateIdle)
		case args[0] == "m" && args[1] == "l":
			trades = b.c.GetAllUniqueTrades(types.ExchangeBitmex, types.SideLong, types.StateIdle)
		default:
			msg := "Wrong arguments. Valid arguments are: b [s|l] and m [s|l]"
			b.MsgChan <- types.BotMessage{ChatID: userID, MsgStr: msg}
			return nil
		}
		for _, t := range trades {
			ctx = context.WithValue(ctx, models.KeyTrade{}, *t)
			go b.HandleExecute(u, ctx)
		}
	}
	return nil
}
