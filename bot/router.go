package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler func(*tgbotapi.Update, context.Context)
type ErrorHandler func(*tgbotapi.Update, context.Context) error
type Middleware func(Handler) Handler
type Router struct {
	middlewares []Middleware
	handlers    map[string]Handler
}

func (r *Router) Use(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) Handle(command string, handler Handler) {
	r.handlers[command] = handler
}
