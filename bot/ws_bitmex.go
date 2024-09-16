package bot

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shahinrahimi/teletradebot/utils"
)

const (
	wsBitmexURL        string        = "wss://www.bitmex.com/realtime"
	wsBitmexURLTestnet string        = "wss://testnet.bitmex.com"
	endpoint           string        = "/realtime"
	Symbol             string        = "ETHUSDT"
	pingInterval       time.Duration = 5 * time.Second
	pongWait           time.Duration = 5 * time.Second
)

func (b *Bot) StartWebsocketServiceBitmex(ctx context.Context) {
	go b.startUserDataStreamBitmexReconnect(ctx)
	//b.Test()
}

func generateSignature(apiSecret, verb, endpoint string, nonce int64) string {
	message := verb + endpoint + strconv.FormatInt(nonce, 10)
	hash := hmac.New(sha256.New, []byte(apiSecret))
	hash.Write([]byte(message))
	return hex.EncodeToString(hash.Sum(nil))
}

func (b *Bot) startUserDataStreamBitmex(ctx context.Context) {
	ws, _, err := websocket.DefaultDialer.Dial(wsBitmexURLTestnet+endpoint, nil)
	if err != nil {
		b.l.Fatal("error dialing bitmex websocket: %v", err)
	}

	defer ws.Close()

	nonce := (time.Now().Unix() + 5) * 1000
	signature := generateSignature(b.mc.ApiSec, "GET", endpoint, nonce)

	authMessage := fmt.Sprintf(`{"op": "authKeyExpires", "args": ["%s", %d, "%s"]}`, b.mc.ApiKey, nonce, signature)

	err = ws.WriteMessage(websocket.TextMessage, []byte(authMessage))
	if err != nil {
		b.l.Fatalf("error sending bitmex auth message: %v", err)
	}
	b.l.Printf("sent auth message: %s", authMessage)

	// subscribe to public channels
	// publicSubMessage := map[string]interface{}{
	// 	"op":   "subscribe",
	// 	"args": []string{"announcement", "chat", "connected", "publicNotifications"},
	// }
	// messageByte, err := json.Marshal(publicSubMessage)
	// if err != nil {
	// 	b.l.Fatalf("error sending bitmex public sub message: %v", err)
	// }
	publicSubMessage := fmt.Sprintf(`{"op": "subscribe", ["trade:%s", "instrument:%s"]}`, Symbol, Symbol)

	err = ws.WriteMessage(websocket.TextMessage, []byte(publicSubMessage))
	if err != nil {
		b.l.Fatalf("error sending bitmex public sub message: %v", err)
	}
	b.l.Printf("sent public sub message: %s", publicSubMessage)

	// subscribe to private channel
	privateSubMessage := fmt.Sprintf(`{"op": "subscribe", "args": ["%s" , "%s", "%s"]}`, "execution", "order", "margin")

	err = ws.WriteMessage(websocket.TextMessage, []byte(privateSubMessage))
	if err != nil {
		b.l.Fatalf("error sending bitmex private sub message: %v", err)
	}
	b.l.Printf("sent private sub message: %s", privateSubMessage)

	//start reading message
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			b.l.Fatalf("error reading bitmex message: %v", err)
		}
		b.l.Printf("received message: %s", message)
	}
}

func (b *Bot) startUserDataStreamBitmexReconnect(ctx context.Context) {
	for {
		ws, _, err := websocket.DefaultDialer.Dial(wsBitmexURLTestnet+endpoint, nil)
		if err != nil {
			b.l.Fatal("error dialing bitmex websocket: %v", err)
		}

		defer ws.Close()

		nonce := (time.Now().Unix() + 5) * 1000
		signature := generateSignature(b.mc.ApiSec, "GET", endpoint, nonce)

		authMessage := fmt.Sprintf(`{"op": "authKeyExpires", "args": ["%s", %d, "%s"]}`, b.mc.ApiKey, nonce, signature)

		err = ws.WriteMessage(websocket.TextMessage, []byte(authMessage))
		if err != nil {
			b.l.Fatalf("error sending bitmex auth message: %v", err)
		}
		b.l.Printf("sent auth message: %s", authMessage)

		publicSubMessage := fmt.Sprintf(`{"op": "subscribe", ["instrument:%s"]}`, Symbol)

		err = ws.WriteMessage(websocket.TextMessage, []byte(publicSubMessage))
		if err != nil {
			b.l.Fatalf("error sending bitmex public sub message: %v", err)
		}
		b.l.Printf("sent public sub message: %s", publicSubMessage)

		// subscribe to private channel
		privateSubMessage := fmt.Sprintf(`{"op": "subscribe", "args": ["%s" , "%s", "%s"]}`, "execution", "order", "margin")

		err = ws.WriteMessage(websocket.TextMessage, []byte(privateSubMessage))
		if err != nil {
			b.l.Fatalf("error sending bitmex private sub message: %v", err)
		}
		b.l.Printf("sent private sub message: %s", privateSubMessage)

		// Timer for heartbeat
		lastMessageTime := time.Now()
		pingTicker := time.NewTicker(pingInterval)
		defer pingTicker.Stop()

		// Track pong receipt
		pongReceived := make(chan struct{})

		ws.SetPongHandler(func(appData string) error {
			b.l.Printf("received pong from BitMEX: %s", appData)
			lastMessageTime = time.Now() // Reset the last message time
			// Signal pong received
			select {
			case pongReceived <- struct{}{}:
			default:
			}

			return nil
		})

		done := make(chan struct{})

		// WebSocket read loop
		var orderTable OrderTable
		var marginTable MarginTable
		var executionTable ExecutionTable
		go func() {
			defer close(done)
			for {
				_, message, err := ws.ReadMessage()
				if err != nil {
					b.l.Printf("error reading bitmex message: %v", err)
					return
				}
				lastMessageTime = time.Now()

				var baseMessage struct {
					Info      string `json:"info"`
					Status    int    `json:"status"`
					Error     string `json:"error"`
					Success   bool   `json:"success"`
					Subscribe string `json:"subscribe"`
					Table     string `json:"table"`
				}
				if err := json.Unmarshal(message, &baseMessage); err != nil {
					b.l.Printf("error unmarshalling bitmex message: %v", err)
					continue
				}
				switch {
				case baseMessage.Info != "":
					b.l.Printf("received message info: %s", baseMessage.Info)
				case baseMessage.Status != 0:
					b.l.Printf("received message status: %d with error: %s", baseMessage.Status, baseMessage.Error)
				case baseMessage.Success:
					b.l.Printf("received message success on %v", baseMessage.Subscribe)
				case baseMessage.Table != "":
					b.l.Printf("received message table: %s", baseMessage.Table)
					switch baseMessage.Table {
					case "order":
						if err := json.Unmarshal(message, &orderTable); err == nil {
							b.l.Printf("received message orderTable: %s", orderTable.Table)
						} else {
							b.l.Printf("error unmarshalling bitmex message: %v", err)
						}
					case "margin":
						if err := json.Unmarshal(message, &marginTable); err == nil {
							b.l.Printf("received message marginTable: %s", marginTable.Table)
						} else {
							b.l.Printf("error unmarshalling bitmex message: %v", err)
						}
					case "execution":
						if err := json.Unmarshal(message, &executionTable); err == nil {
							b.l.Printf("received message executionTable: %s", executionTable.Table)
						} else {
							b.l.Printf("error unmarshalling bitmex message: %v", err)
						}
					default:
						b.l.Printf("received unknown message table: %s", baseMessage.Table)

					}
				default:
					b.l.Printf("received unknown message type: %s", message)
				}

			}
		}()

		// Ping/Pong handleing
		go func() {
			for {
				select {
				case <-pingTicker.C:
					b.l.Printf("last message time: %s", utils.FriendlyDuration(time.Since(lastMessageTime)))
					if time.Since(lastMessageTime) >= pingInterval {
						b.l.Println("ping bitmex")

						err := ws.WriteMessage(websocket.PingMessage, nil)
						if err != nil {
							b.l.Printf("error sending bitmex ping: %v", err)
							return
						}
						// wait for pong
						pongWaitTimer := time.NewTimer(pongWait)
						select {
						case <-pongReceived:
							b.l.Println("pong received in time")
							pongWaitTimer.Stop()
						case <-pongWaitTimer.C:
							b.l.Println("pong timeout - reconnecting")
							ws.Close()
							return
						}
					}
				}
			}
		}()

		// wait for websocket to close or error out
		<-done
		b.l.Println("Connection lost. Reconnecting...")
		time.Sleep(5 * time.Second)

	}
}

func (b *Bot) Test() {
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
