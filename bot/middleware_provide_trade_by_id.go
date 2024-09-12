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
		userID := u.Message.From.ID
		args := strings.Split(u.Message.CommandArguments(), " ")
		id, err := strconv.Atoi(args[0])
		if err != nil {
			msg := fmt.Sprintf("Invalid trade ID: '%s'. Please provide a numeric ID.", args[0])
			b.MsgChan <- BotMessage{
				ChatID: userID,
				MsgStr: msg,
			}
			return
		}
		o, err := b.s.GetTrade(int64(id))
		if err != nil {
			if err == sql.ErrNoRows {
				msg := fmt.Sprintf("No trade found with ID: %d.", id)
				b.MsgChan <- BotMessage{
					ChatID: userID,
					MsgStr: msg,
				}
				return
			}
			b.l.Printf("Error retrieving trade from database: %v", err)
			msg := "An internal error occurred while fetching the trade. Please try again later."
			b.MsgChan <- BotMessage{
				ChatID: userID,
				MsgStr: msg,
			}
			return
		}
		ctx = context.WithValue(ctx, models.KeyTrade{}, *o)
		next(u, ctx)
	}
}
