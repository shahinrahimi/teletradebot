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

	b, err := bot.NewBot(logger, s, token)
	if err != nil {
		logger.Fatalf("error creating instance of bot: %v", err)
	}

	// global middleware
	b.Use(b.Logger)

	// routes
	r := b.NewRouter("trades")
	r.Handle(bot.HELP, b.MakeHandlerBotFunc(b.HandleHelp))
	r.Handle(bot.ADD, b.MakeHandlerBotFunc(b.HandleAdd))
	r.Handle(bot.REMOVE, b.MakeHandlerBotFunc(b.HandleRemove))
	r.Handle(bot.CHECK, b.MakeHandlerBotFunc(b.HandleCheck))
	r.Handle(bot.CANCEL, b.MakeHandlerBotFunc(b.HandleCancel))

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
