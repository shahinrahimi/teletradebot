package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
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
