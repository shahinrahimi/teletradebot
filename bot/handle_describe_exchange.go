package bot

import (
	"context"

	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) handleDescribeExchange(ctx context.Context, t *models.Trade, userID int64, ex exchange.Exchange) {

	i, err := b.retry(config.MaxTries, config.WaitForNextTries, t, func() (interface{}, error) {
		return ex.FetchInterpreter(ctx, t)
	})
	if err != nil {
		b.l.Printf("error fetching interpreter: %v", err)
		b.handleError(err, userID, t.ID)
		return
	}
	interpreter, ok := i.(*models.Interpreter)
	if !ok {
		b.l.Panicf("unexpected error happened in casting error to *models.Interpreter: %T", interpreter)
	}
	b.MsgChan <- types.BotMessage{
		ChatID: userID,
		MsgStr: interpreter.Describe(false),
	}
}
