package main

import (
	"log"
	"net/url"

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
}
