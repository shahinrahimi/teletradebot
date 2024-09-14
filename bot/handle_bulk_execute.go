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
		var trades map[string]models.Trade
		var err error
		switch {
		case args[0] == "b" && args[1] == "s":
			trades, err = b.getAllUniqueTrades(types.ACCOUNT_B, types.SIDE_S, types.STATE_IDLE)
			if err != nil {
				return err
			}
		case args[0] == "b" && args[1] == "l":
			trades, err = b.getAllUniqueTrades(types.ACCOUNT_B, types.SIDE_L, types.STATE_IDLE)
			if err != nil {
				return err
			}
		case args[0] == "m" && args[1] == "s":
			trades, err = b.getAllUniqueTrades(types.ACCOUNT_M, types.SIDE_S, types.STATE_IDLE)
			if err != nil {
				return err
			}
		case args[0] == "m" && args[1] == "l":
			trades, err = b.getAllUniqueTrades(types.ACCOUNT_M, types.SIDE_L, types.STATE_IDLE)
			if err != nil {
				return err
			}
		default:
			msg := "Wrong arguments. Valid arguments are: b [s|l] and m [s|l]"
			b.MsgChan <- types.BotMessage{ChatID: userID, MsgStr: msg}
			return nil
		}
		for _, t := range trades {
			ctx = context.WithValue(ctx, models.KeyTrade{}, t)
			go b.HandleExecute(u, ctx)
		}
	}
	return nil
}
