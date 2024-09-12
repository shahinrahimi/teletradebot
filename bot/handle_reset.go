package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/exchange/binance"
	"github.com/shahinrahimi/teletradebot/models"
)

func (b *Bot) HandleReset(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	userID := u.Message.From.ID
	if err := b.s.UpdateTradeIdle(&t); err != nil {
		b.l.Printf("Error updating the trade status: %v", err)
		return err
	}
	_, exist := binance.TradeDescribers[t.ID]
	if exist {
		delete(binance.TradeDescribers, t.ID)
	}
	msg := fmt.Sprintf("The trade has been successfully reset.\n\nTrade ID: %d", t.ID)
	b.MsgChan <- BotMessage{
		ChatID: userID,
		MsgStr: msg,
	}

	return nil
}
