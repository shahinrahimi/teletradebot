package bot

import (
	"context"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/cash"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange"
	"github.com/shahinrahimi/teletradebot/types"
)

type Bot struct {
	l           *log.Logger
	c           *cash.Cash
	api         *tgbotapi.BotAPI
	routers     map[string]*Router
	middlewares []Middleware
	//bc          *binance.BinanceClient
	//mc          *bitmex.BitmexClient
	bc      exchange.Exchange
	mc      exchange.Exchange
	MsgChan chan types.BotMessage
	DbgChan chan string
}

func NewBot(l *log.Logger, c *cash.Cash, bc exchange.Exchange, mc exchange.Exchange, token string, msgChan chan types.BotMessage, dbgChan chan string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		l.Printf("error creating a new bot api: %v", err)
		return nil, err
	}
	return &Bot{
		l:           l,
		c:           c,
		api:         api,
		routers:     make(map[string]*Router),
		middlewares: []Middleware{},
		bc:          bc,
		mc:          mc,
		MsgChan:     msgChan,
		DbgChan:     dbgChan,
	}, nil
}

func (b *Bot) Debug(msg string) {
	if config.Debug {
		b.l.Println(msg)
	}
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
			b.l.Printf("received a ctx.Done")
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
		// remove case sensitivity for commands
		command = strings.ToLower(command)
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
	go b.startListenForDbgMessages()
	go b.startListenForBotMessages()
}

func (b *Bot) startListenForDbgMessages() {
	go func() {
		for msg := range b.DbgChan {
			b.Debug(msg)
		}
	}()
}

func (b *Bot) startListenForBotMessages() {
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

// Shutdown stops the go routine which receives updates by simply call the StopReceivingUpdates
func (b *Bot) Shutdown() {
	b.l.Println("Shutting down ...")
	b.api.StopReceivingUpdates()
}
