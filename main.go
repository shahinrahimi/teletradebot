package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"gihub.com/shahinrahimi/teletradebot/bot"
	"gihub.com/shahinrahimi/teletradebot/store"
	"github.com/joho/godotenv"
)

func main() {
	// create custom logger
	logger := log.New(os.Stdout, "[TELETRADE-BOT] ", log.LstdFlags)

	// check .env file
	if err := godotenv.Load(); err != nil {
		logger.Fatalf("error loading environmental file: %v", err)
	}

	// check environmental variable
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		logger.Fatal("error wrong environmental variable")
	}

	// create a store
	s, err := store.NewSqliteStore(logger)
	if err != nil {
		logger.Fatalf("error creating new sqlite store instance: %v", err)
	}
	defer s.CloseDB()

	// init DB
	if err := s.Init(); err != nil {
		logger.Fatalf("error initializing DB: %v", err)
	}

	b, err := bot.NewBot(logger, s, token)
	if err != nil {
		logger.Fatalf("error creating instance of bot: %v", err)
	}

	// global middleware
	b.Use(b.Logger)

	// routes
	// help route
	r := b.NewRouter("help")
	r.Handle(bot.HELP, b.MakeHandlerBotFunc(b.HandleHelp))
	r.Handle(bot.START, b.MakeHandlerBotFunc(b.HandleHelp))
	// new route
	r1 := b.NewRouter("route-1")
	r1.Handle(bot.ADD, b.MakeHandlerBotFunc(b.HandleAdd))
	r1.Use(b.ProvideAddOrder)
	// get delete cancel execute
	r2 := b.NewRouter("route-2")
	r2.Handle(bot.LIST, b.MakeHandlerBotFunc(b.HandleList))
	r2.Handle(bot.REMOVE, b.MakeHandlerBotFunc(b.HandleRemove))
	r2.Handle(bot.CHECK, b.MakeHandlerBotFunc(b.HandleCheck))
	r2.Handle(bot.CANCEL, b.MakeHandlerBotFunc(b.HandleCancel))

	// create context bot to received updates and gracefully shutdown
	ctx := context.WithoutCancel(context.Background())
	go func() {
		logger.Println("Bot started and running and listen for updates.")
		b.Start(ctx)
	}()

	// create signal
	c := make(chan os.Signal, 1)
	// filter all other signal
	signal.Notify(c, os.Interrupt)

	// block until a signal received
	rc := <-c
	logger.Println("go signal", rc)

	// gracefully shutdown bot, waiting max 30 secs
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	b.Shutdown()

}
