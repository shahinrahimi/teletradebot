package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	binance "github.com/adshao/go-binance/v2"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const (
	binanceWsBaseURL = "wss://fstream.binance.com"
)

// Stream to subscribe to
var streams = []string{
	"btcusdt@aggTrade",
	"btcusdt@depth",
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	apiKey := os.Getenv("BIANCE_API_KEY_FUTURES_TESTNET")
	apiSec := os.Getenv("BIANCE_API_SEC_FUTURES_TESTNET")
	client := binance.NewClient(apiKey, apiSec)
	//fetureClient := binance.NewFuturesClient(apiKey, apiSec)
	prices, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, p := range prices {
		fmt.Println(p)
	}
}

func binance_websocket() {
	// Create a channel to handle interrupt signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Build the WebSocket URL
	wsURL := binanceWsBaseURL + "/stream?streams=" + streams[0]
	for _, stream := range streams[1:] {
		wsURL += "/" + stream
	}

	log.Printf("Connecting to %s", wsURL)

	// Connect to the WebSocket server
	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			log.Printf("Received: %s", message)
		}
	}()

	// Ping/Pong Handling
	c.SetPingHandler(func(appData string) error {
		log.Println("Ping received, sending pong")
		return c.WriteMessage(websocket.PongMessage, []byte(appData))
	})

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println("Ping error:", err)
				return
			}
			log.Printf("Ping sent at %v", t)
		case <-interrupt:
			log.Println("Interrupt received, closing connection")
			// Cleanly close the connection by sending a close message
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Write close error:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func bitmex_websocket() {
	// Define the WebSocket server URL
	serverURL := "wss://ws.bitmex.com/realtimePlatform"

	// Parse the URL
	u, err := url.Parse(serverURL)
	if err != nil {
		log.Fatal("Error parsing URL:", err)
	}

	// Create a new WebSocket dialer
	dialer := websocket.DefaultDialer

	// Connect to the WebSocket server
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()

	log.Printf("Connected to %s\n", serverURL)

	// Setup a channel to handle OS signals (e.g., to close the connection on interrupt)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Channel for received messages
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}
			log.Printf("Received: %s\n", message)
		}
	}()

	// Prepare the subscription message
	subscribeMessage := map[string]interface{}{
		"op":   "subscribe",
		"args": []string{"announcement", "chat", "connected", "publicNotifications"},
	}

	// Send the subscription message
	messageBytes, err := json.Marshal(subscribeMessage)
	if err != nil {
		log.Println("Error marshaling subscribe message:", err)
		return
	}

	// Send the subscription message
	err = conn.WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		log.Println("Error sending subscription message:", err)
		return
	}
	log.Printf("Subscribed to topics: %v\n", subscribeMessage["args"])

	// Handle incoming messages and OS signals
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Interrupt received, closing connection...")

			// Cleanly close the connection by sending a close message
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during close:", err)
				return
			}

			// Wait for the connection to close
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
