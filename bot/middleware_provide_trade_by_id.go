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
			b.SendMessage(u.Message.From.ID, fmt.Sprintf("Invalid trade ID: '%s'. Please provide a numeric ID.", args[0]))
			return
		}
		o, err := b.s.GetTrade(id)
		if err != nil {
			if err == sql.ErrNoRows {
				b.SendMessage(u.Message.From.ID, fmt.Sprintf("No trade found with ID: %d.", id))
				return
			}
			b.l.Printf("Error retrieving trade from database: %v", err)
			b.SendMessage(u.Message.From.ID, "An internal error occurred while fetching the trade. Please try again later.")
			return
		}
		ctx = context.WithValue(ctx, models.KeyTrade{}, *o)
		next(u, ctx)
	}
}
