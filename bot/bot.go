package bot

import (
	"context"
	"log"

	"gihub.com/shahinrahimi/teletradebot/exchange"
	"gihub.com/shahinrahimi/teletradebot/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Shortcuts = map[string]string{
	"1": "m XBTUSD long 1h 0.1 1 105 101 1",
	"2": "b BTCUSDT long 4h 0.2 1 105 101 1",
	"h": "m XBTUSD long 1h 0.1 1 105 101 1",
}

type Bot struct {
	l           *log.Logger
	s           store.Storage
	api         *tgbotapi.BotAPI
	routers     map[string]*Router
	middlewares []Middleware
	bc          *exchange.BinanceClient
	mc          *exchange.BitmexClient
}

func NewBot(l *log.Logger, s store.Storage, bc *exchange.BinanceClient, mc *exchange.BitmexClient, token string) (*Bot, error) {
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
			return
		case u := <-us:
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
func (b *Bot) SendMessage(userID int64, msgStr string) {
	msg := tgbotapi.NewMessage(userID, msgStr)
	if _, err := b.api.Send(msg); err != nil {
		b.l.Printf("error in sending message to user: %v", err)
	}
}

// Shutdown stops the go routine which receives updates by simply call the StopReceivingUpdates
func (b *Bot) Shutdown() {
	b.l.Println("Shutting down ...")
	b.api.StopReceivingUpdates()
}
