package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleDescribe(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	if t.Account == types.ACCOUNT_B {
		dt, err := b.bc.GetTradeDescriber(ctx, &t)
		if err != nil {
			return err
		}
		sizeString := fmt.Sprintf("%.1f%%", float64(t.Size))
		slString := fmt.Sprintf("%.1f%%", float64((t.StopLoss - 100)))
		tpString := fmt.Sprintf("%.1f%%", float64((t.TakeProfit - 100)))
		msg := fmt.Sprintf("Trade ID %d\n\n", t.ID)
		msg = fmt.Sprintf("%s From: %s\n Till: %s\n Open: %s\n High: %s\n Low: %s\n Close: %s\n\n", msg, dt.From, dt.Till, dt.Open, dt.High, dt.Low, dt.Close)
		msg = fmt.Sprintf("%sTrading:\n", msg)
		msg = fmt.Sprintf("%s Entry %s at %s with %s of balance.\n", msg, dt.Side, dt.SP, sizeString)
		msg = fmt.Sprintf("%s TP at %s with %s.\n", msg, dt.TP, tpString)
		msg = fmt.Sprintf("%s SL at %s with %s.\n", msg, dt.SL, slString)
		b.SendMessage(u.Message.From.ID, msg)

	} else {
		return b.HandleUnderDevelopment(u, ctx)
	}
	return nil
}
