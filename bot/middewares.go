package bot

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gihub.com/shahinrahimi/teletradebot/config"
	"gihub.com/shahinrahimi/teletradebot/models"
	"gihub.com/shahinrahimi/teletradebot/types"
	"gihub.com/shahinrahimi/teletradebot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) Logger(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		var userID int64 = u.Message.From.ID
		if u.Message.Command() == "" {
			b.l.Printf("Received message (%d): %s", userID, u.Message.Text)
		} else {
			b.l.Printf("Received command (%d): %s\t args: %s", userID, u.Message.Command(), u.Message.CommandArguments())
		}
		next(u, ctx)
	}
}

func (b *Bot) RequiredAuth(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		var userID int64 = u.Message.From.ID
		for _, id := range config.UserIDs {
			if id == userID {
				next(u, ctx)
				return
			}
		}
		b.SendMessage(userID, "You are not allowed")
	}
}

func (b *Bot) BanBots(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		if u.Message.From.IsBot {
			return
		}
		next(u, ctx)
	}
}

func (b *Bot) ProvideAddTrade(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		var tradeArgs []string
		var userID int64 = u.Message.From.ID
		args := strings.Split(u.Message.CommandArguments(), " ")
		// check shortcuts
		if len(args) == 1 {
			value, exist := config.Shortcuts[args[0]]
			if !exist {
				b.SendMessage(userID, fmt.Sprintf("shortcut is not found: '%s'", args[0]))
				return
			}
			tradeArgs = strings.Split(value, " ")
		} else {
			tradeArgs = args
		}
		t, err := utils.ParseTrade(tradeArgs)
		if err != nil {
			b.SendMessage(userID, err.Error())
			return
		}

		var isAvailable bool = false
		switch t.Account {
		case types.ACCOUNT_B:
			// check pair for binance
			isAvailable = b.bc.CheckSymbol(t.Symbol)
			if !isAvailable {
				b.l.Printf("error checking symbol '%s' availability symbol", t.Symbol)
			}
		case types.ACCOUNT_M:
			// check pair for bitmex
			// TODO implement a method to check if pair is available in bitmex
			c := http.Client{}
			resp, err := c.Get(fmt.Sprintf("https://www.bitmex.com/api/udf/symbols?symbol=%s", t.Symbol))
			if resp.StatusCode == 200 || err != nil {
				isAvailable = true
			} else {
				b.l.Panicf("error checking symbol '%s' availability symbol", t.Symbol)
			}
		default:
			// should never happen

			b.l.Panicf("error checking symbol '%s' availability symbol", t.Symbol)

		}

		if !isAvailable {
			b.SendMessage(userID, fmt.Sprintf("the symbol '%s' is not available for exchange '%s'.", t.Symbol, t.Account))
			return
		}

		// add UserID
		t.UserID = userID

		ctx = context.WithValue(ctx, models.KeyTrade{}, *t)
		next(u, ctx)
	}
}

func (b *Bot) ProvideTradeByID(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		args := strings.Split(u.Message.CommandArguments(), " ")
		id, err := strconv.Atoi(args[0])
		if err != nil {
			b.SendMessage(u.Message.From.ID, fmt.Sprintf("the id is not valid: %s", args[0]))
			return
		}
		o, err := b.s.GetTrade(id)
		if err != nil {
			if err == sql.ErrNoRows {
				b.SendMessage(u.Message.From.ID, fmt.Sprintf("the trade not found with id: %d", id))
				return
			}
			b.l.Printf("internal error getting the trade from DB")
			b.SendMessage(u.Message.From.ID, fmt.Sprintf("the id is not valid: %d", id))
			return
		}
		ctx = context.WithValue(ctx, models.KeyTrade{}, *o)
		next(u, ctx)
	}
}
