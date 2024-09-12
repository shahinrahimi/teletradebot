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

func (b *Bot) ProvideAddTrade(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		var tradeArgs []string
		var userID int64 = u.Message.From.ID
		var chatID int64 = u.Message.Chat.ID
		args := strings.Split(u.Message.CommandArguments(), " ")
		// check shortcuts
		if len(args) == 1 {
			value, exist := config.Shortcuts[args[0]]
			if !exist {
				msg := fmt.Sprintf("Shortcut not found: '%s'", args[0])
				b.MsgChan <- BotMessage{
					ChatID: userID,
					MsgStr: msg,
				}
				return
			}
			tradeArgs = strings.Split(value, " ")
		} else {
			tradeArgs = args
		}
		t, err := utils.ParseTrade(tradeArgs)
		if err != nil {
			b.MsgChan <- BotMessage{
				ChatID: userID,
				MsgStr: err.Error(),
			}
			return
		}

		var isAvailable bool = false
		switch t.Account {
		case types.ACCOUNT_B:
			// check pair for binance
			isAvailable = b.bc.CheckSymbol(t.Symbol)
		case types.ACCOUNT_M:
			// check pair for bitmex
			isAvailable = b.mc.CheckSymbol(t.Symbol)
		default:
			// should never happen
			b.l.Panicf("Unexpected account type while checking symbol '%s'", t.Symbol)

		}

		if !isAvailable {
			b.l.Printf("Error checking symbol '%s' availability on %s", t.Symbol, t.Account)
			msg := fmt.Sprintf("the symbol '%s' is not available for exchange '%s'.", t.Symbol, t.Account)
			b.MsgChan <- BotMessage{
				ChatID: userID,
				MsgStr: msg,
			}
			return
		}

		// add UserID and ChatID
		t.UserID = userID
		t.ChatID = chatID

		ctx = context.WithValue(ctx, models.KeyTrade{}, *t)
		next(u, ctx)
	}
}
