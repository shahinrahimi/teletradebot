package bot

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
)

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
