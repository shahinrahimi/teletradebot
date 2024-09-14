package bot

import (
	"context"
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
		var trades map[string]models.Trade
		var err error
		switch {
		case args[0] == "b" && args[1] == "s":
			trades, err = b.getAllUniqueRawTrades(rawTrades, types.ACCOUNT_B, types.SIDE_S)
			if err != nil {
				return err
			}
		case args[0] == "b" && args[1] == "l":
			trades, err = b.getAllUniqueRawTrades(rawTrades, types.ACCOUNT_B, types.SIDE_L)
			if err != nil {
				return err
			}
		case args[0] == "m" && args[1] == "s":
			trades, err = b.getAllUniqueRawTrades(rawTrades, types.ACCOUNT_M, types.SIDE_S)
			if err != nil {
				return err
			}
		case args[0] == "m" && args[1] == "l":
			trades, err = b.getAllUniqueRawTrades(rawTrades, types.ACCOUNT_M, types.SIDE_L)
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
			go b.HandleAdd(u, ctx)
		}
	}
	return nil
}
