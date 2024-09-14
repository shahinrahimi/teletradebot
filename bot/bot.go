package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/adshao/go-binance/v2/common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/exchange/binance"
	"github.com/shahinrahimi/teletradebot/exchange/bitmex"
	"github.com/shahinrahimi/teletradebot/store"
	"github.com/shahinrahimi/teletradebot/types"
)

type Bot struct {
	l           *log.Logger
	s           store.Storage
	api         *tgbotapi.BotAPI
	routers     map[string]*Router
	middlewares []Middleware
	bc          *binance.BinanceClient
	mc          *bitmex.BitmexClient
	MsgChan     chan types.BotMessage
}

func NewBot(l *log.Logger, s store.Storage, bc *binance.BinanceClient, mc *bitmex.BitmexClient, token string, msgChan chan types.BotMessage) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		l.Printf("error creating a new bot api: %v", err)
		return nil, err
	}
	return &Bot{
		l:           l,
		s:           s,
		api:         api,
		routers:     make(map[string]*Router),
		middlewares: []Middleware{},
		bc:          bc,
		mc:          mc,
		MsgChan:     msgChan,
	}, nil
}

func (b *Bot) NewRouter(routeName string) *Router {
	router := &Router{
		handlers: make(map[string]Handler),
	}
	b.routers[routeName] = router
	return router
}

func (b *Bot) Use(m Middleware) {
	b.middlewares = append(b.middlewares, m)
}

func (b *Bot) Start(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	us := b.api.GetUpdatesChan(u)
	go b.receiveUpdates(ctx, us)
}

// receivedUpdates check sign
func (b *Bot) receiveUpdates(ctx context.Context, us tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			b.l.Printf("recived a ctx.Done")
			return
		case u := <-us:
			if u.EditedMessage != nil {
				b.l.Printf("received an edited message, bypassing update")
				continue
			}
			b.handleUpdate(u, ctx)
		}
	}
}

func (b *Bot) handleUpdate(u tgbotapi.Update, ctx context.Context) {

	// Define the function that will handle the routing logic
	routerHandler := func(u *tgbotapi.Update, ctx context.Context) {
		command := u.Message.Command()
		for _, router := range b.routers {
			if handler, exists := router.handlers[command]; exists {
				// Start with the actual handler
				finalHandler := handler

				// Apply route-specific middlewares in reverse order
				for i := len(router.middlewares) - 1; i >= 0; i-- {
					finalHandler = router.middlewares[i](finalHandler)
				}

				// Execute the final composed handler
				finalHandler(u, ctx)
				return
			}
		}
		if command != "" {
			// Handle unknown command
			b.l.Printf("Unknown command: %s", command)
		}
	}

	// Start with the routing handler
	finalHandler := routerHandler

	// Apply global middlewares in reverse order
	for i := len(b.middlewares) - 1; i >= 0; i-- {
		finalHandler = b.middlewares[i](finalHandler)
	}

	// Execute the final composed handler (global middlewares + router handling)
	finalHandler(&u, ctx)
}

// SendMessage send message string to user and error does not returned
func (b *Bot) StartMessageListener() {
	go func() {
		for msg := range b.MsgChan {
			b.sendMessage(msg.ChatID, msg.MsgStr)
		}
	}()
}
func (b *Bot) sendMessage(chatID int64, msgStr string) {
	msg := tgbotapi.NewMessage(chatID, msgStr)
	if _, err := b.api.Send(msg); err != nil {
		b.l.Printf("error in sending message to user: %v", err)
	}
}

func (b *Bot) handleError(err error, userID int64, tradeID int64) {
	if apiErr, ok := err.(*common.APIError); ok {
		msg := fmt.Sprintf("Binance API:\n\nMessage: %s\nCode: %d\nTradeID: %d", apiErr.Message, apiErr.Code, tradeID)
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: msg,
		}
	} else {
		b.l.Printf("error casting error to Api error type: %T", err)
	}
}

// Shutdown stops the go routine which receives updates by simply call the StopReceivingUpdates
func (b *Bot) Shutdown() {
	b.l.Println("Shutting down ...")
	b.api.StopReceivingUpdates()
}
