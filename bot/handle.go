package bot

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2/common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	swagger "github.com/shahinrahimi/teletradebot/swagger"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) handleBinanceAPIError(err *common.APIError, userID int64) {
	msg := fmt.Sprintf("Binance API\n\ncode: %d\nmessage: %s", err.Code, err.Message)
	b.MsgChan <- types.BotMessage{
		ChatID: userID,
		MsgStr: msg,
	}
}

func (b *Bot) handleBitmexAPIError(err swagger.GenericSwaggerError, userID int64) {
	msg := fmt.Sprintf("Bitmex API\n\nBody: %s\nError: %s", err.Body(), err.Error())
	b.MsgChan <- types.BotMessage{
		ChatID: userID,
		MsgStr: msg,
	}
}

func (b *Bot) handleAPIError(err error, userID int64) {
	if apiErr, ok := err.(*common.APIError); ok {
		msg := fmt.Sprintf("Binance API:\n\\nMessage: %s\nCode: %d", apiErr.Message, apiErr.Code)
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
	} else {
		b.l.Printf("error casting error to Api error type: %T", err)
	}
}

func (b *Bot) MakeHandlerBotFunc(f ErrorHandler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		if err := f(u, ctx); err != nil {
			b.l.Printf("we have error %v", err)
			b.l.Printf("err type: %T", err)
			if ApiErr, ok := (err).(swagger.GenericSwaggerError); ok {
				b.l.Printf("error: %s , errof: %s", string(ApiErr.Body()), ApiErr.Error())
			}
			if apiErr, ok := err.(*common.APIError); ok {
				msg := fmt.Sprintf("Binance API:\ncould not place a order for trade\ncode:%d\nmessage: %s", apiErr.Code, apiErr.Message)
				b.l.Println(msg)
			}
		}
	}
}
