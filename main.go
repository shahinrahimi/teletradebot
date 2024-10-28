package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/shahinrahimi/teletradebot/bot"
	"github.com/shahinrahimi/teletradebot/cash"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/exchange/binance"
	"github.com/shahinrahimi/teletradebot/exchange/bitmex"
	"github.com/shahinrahimi/teletradebot/logger"
	"github.com/shahinrahimi/teletradebot/store"
	"github.com/shahinrahimi/teletradebot/types"
)

func main() {
	// create custom logger
	l := logger.New("TELETRADE-BOT", config.SaveToDisk)
	if config.SaveToDisk {
		defer l.Close()
	}

	// create global context
	ctx := context.WithoutCancel(context.Background())

	// create global message channel
	msgChan := make(chan types.BotMessage)
	dbgChan := make(chan string)

	// create global store (Storage)
	s, err := store.NewSqliteStore(l.Log)
	if err != nil {
		l.Log.Fatalf("error creating new sqlite store instance: %v", err)
	}
	defer s.CloseDB()

	// init DB
	if err := s.Init(); err != nil {
		l.Log.Fatalf("error initializing DB: %v", err)
	}

	c := cash.NewCash(s, l.Log)

	// check .env file
	if err := godotenv.Load(); err != nil {
		l.Log.Fatalf("error loading environmental file: %v", err)
	}

	// check environmental variable for telegram bot
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		l.Log.Fatal("error wrong environmental variable")

	}

	// check environmental variable for binance api
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSec := os.Getenv("BINANCE_API_SEC")
	if apiKey == "" || apiSec == "" {
		l.Log.Fatal("error wrong environmental variable for binance client")
	}

	// check environmental variable for binance api
	apiKey2 := os.Getenv("BITMEX_API_ID")
	apiSec2 := os.Getenv("BITMEX_API_KEY")
	if apiKey2 == "" || apiSec2 == "" {
		l.Log.Fatal("error wrong environmental variable for bitmex client")
	}

	bc := binance.NewBinanceClient(l.Log, apiKey, apiSec, config.UseBinanceTestnet, dbgChan)
	mc := bitmex.NewBitmexClient(l.Log, apiKey2, apiSec2, config.UseBitmexTestnet, dbgChan)

	// start polling for binance
	bc.StartPolling(ctx)
	// start polling for bitmex
	mc.StartPolling(ctx)

	b, err := bot.NewBot(l.Log, c, bc, mc, token, msgChan, dbgChan)
	if err != nil {
		l.Log.Fatalf("error creating instance of bot: %v", err)
	}

	// start listening for messages
	b.StartMessageListener()

	// start binance ws
	bc.StartWebsocketService(ctx, b.WsHandler, b.WsErrHandler)

	// start bitmex ws
	mc.StartWebsocketService(ctx, b.WsHandlerBitmex)

	// global middleware
	b.Use(b.BanBots)
	b.Use(b.Logger)

	// routes
	// help route
	r := b.NewRouter("help")
	r.Handle(bot.HELP, b.MakeHandlerBotFunc(b.HandleHelp))
	r.Handle(bot.START, b.MakeHandlerBotFunc(b.HandleHelp))
	r.Handle(bot.INFO, b.MakeHandlerBotFunc(b.HandleInfo))

	// list route
	r0 := b.NewRouter("route-0")
	r0.Handle(bot.LIST, b.MakeHandlerBotFunc(b.HandleList))
	r0.Handle(bot.ALIAS, b.MakeHandlerBotFunc(b.HandleAlias))
	// for testing
	// r0.Handle("bulke", b.MakeHandlerBotFunc(b.HandleBulkExecute))
	// r0.Handle("bulkd", b.MakeHandlerBotFunc(b.HandleBulkDelete))
	// r0.Handle("bulka", b.MakeHandlerBotFunc(b.HandleBulkAdd))
	// r0.Handle("bulkr", b.MakeHandlerBotFunc(b.HandleBulkReset))
	// r0.Handle("bulkc", b.MakeHandlerBotFunc(b.HandleBulkClose))
	r0.Use(b.RequiredAuth)

	// new route
	r1 := b.NewRouter("route-1")
	r1.Handle(bot.ADD, b.MakeHandlerBotFunc(b.HandleAdd))
	r1.Use(b.RequiredAuth)
	r1.Use(b.ProvideAddTrade)

	// get delete cancel execute
	r2 := b.NewRouter("route-2")
	r2.Handle(bot.REMOVE, b.MakeHandlerBotFunc(b.HandleRemove))
	r2.Handle(bot.CHECK, b.MakeHandlerBotFunc(b.HandleCheck))
	r2.Handle(bot.EXECUTE, b.MakeHandlerBotFunc(b.HandleExecute))
	r2.Handle(bot.CLOSE, b.MakeHandlerBotFunc(b.HandleClose))
	r2.Handle(bot.VIEW, b.MakeHandlerBotFunc(b.HandleView))
	r2.Handle(bot.DESCRIBE, b.MakeHandlerBotFunc(b.HandleDescribe))
	r2.Handle(bot.RESET, b.MakeHandlerBotFunc(b.HandleReset))
	r2.Use(b.RequiredAuth)
	r2.Use(b.ProvideTradeByID)

	go func() {
		l.Log.Println("Bot running and listen for updates...")
		b.Start(ctx)
	}()

	// create signal
	cc := make(chan os.Signal, 1)
	// filter all other signal
	signal.Notify(cc, os.Interrupt)

	// block until a signal received
	rc := <-cc
	l.Log.Println("got signal", rc)

	// gracefully shutdown bot, waiting max 30 secs
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	b.Shutdown()

}
