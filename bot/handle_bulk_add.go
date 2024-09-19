package bot

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) HandleBulkAdd(u *tgbotapi.Update, ctx context.Context) error {
	args := strings.Split(u.Message.CommandArguments(), " ")
	userID := u.Message.From.ID
	chatID := u.Message.Chat.ID
	var rawTrades []models.Trade
	for _, value := range config.Shortcuts {
		tradeArgs := strings.Split(value, " ")
		t, err := utils.ParseTrade(tradeArgs)
		if err != nil {
			continue
		}
		rawTrades = append(rawTrades, models.Trade{
			UserID:            userID,
			ChatID:            chatID,
			Symbol:            t.Symbol,
			Account:           t.Account,
			State:             t.State,
			Side:              t.Side,
			Timeframe:         t.Timeframe,
			Offset:            t.Offset,
			Size:              t.Size,
			StopLoss:          t.StopLoss,
			TakeProfit:        t.TakeProfit,
			ReverseMultiplier: t.ReverseMultiplier,
		})
	}

	if len(args) == 2 {
		var trades []*models.Trade
		switch {
		case args[0] == "b" && args[1] == "s":
			trades = b.c.GetAllUniqueRawTrades(rawTrades, types.ACCOUNT_B, types.SIDE_S)
		case args[0] == "b" && args[1] == "l":
			trades = b.c.GetAllUniqueRawTrades(rawTrades, types.ACCOUNT_B, types.SIDE_L)
		case args[0] == "m" && args[1] == "s":
			trades = b.c.GetAllUniqueRawTrades(rawTrades, types.ACCOUNT_M, types.SIDE_S)
		case args[0] == "m" && args[1] == "l":
			trades = b.c.GetAllUniqueRawTrades(rawTrades, types.ACCOUNT_M, types.SIDE_L)
		default:
			msg := "Wrong arguments. Valid arguments are: b [s|l] and m [s|l]"
			b.MsgChan <- types.BotMessage{ChatID: userID, MsgStr: msg}
			return nil
		}
		for _, t := range trades {
			// check for symbol availability
			var isAvailable bool = false
			if t.Account == types.ACCOUNT_B {
				isAvailable = b.bc.CheckSymbol(t.Symbol)
			}
			if t.Account == types.ACCOUNT_M {
				isAvailable = b.mc.CheckSymbol(t.Symbol)
			}
			if !isAvailable {
				msg := fmt.Sprintf("Symbol %s not available for account %s", t.Symbol, t.Account)
				b.MsgChan <- types.BotMessage{
					ChatID: userID,
					MsgStr: msg,
				}
				continue
			}
			ctx = context.WithValue(ctx, models.KeyTrade{}, t)
			go b.HandleAdd(u, ctx)
		}
	}
	return nil
}
