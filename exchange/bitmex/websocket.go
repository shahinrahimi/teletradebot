package bitmex

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shahinrahimi/teletradebot/config"
	"github.com/shahinrahimi/teletradebot/swagger"
)

const (
	wsBitmexURL        string        = "wss://www.bitmex.com"
	wsBitmexURLTestnet string        = "wss://testnet.bitmex.com"
	endpoint           string        = "/realtime"
	Symbol             string        = "ETHUSDT"
	pingInterval       time.Duration = 5 * time.Second
	pongWait           time.Duration = 5 * time.Second
)

func (mc *BitmexClient) StartWebsocketService(ctx context.Context, wsHandler func(ctx context.Context, od []swagger.OrderData)) {
	go mc.startUserDataStream724(ctx, wsHandler)
	//b.Test()
}

func (mc *BitmexClient) startUserDataStream724(ctx context.Context, wsHandler func(ctx context.Context, od []swagger.OrderData)) {
	for {
		websocketURL := wsBitmexURL + endpoint
		if config.UseBitmexTestnet {
			websocketURL = wsBitmexURLTestnet + endpoint
		}
		ws, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
		if err != nil {
			mc.l.Fatalf("error dialing bitmex websocket: %v", err)
		}

		defer ws.Close()

		nonce := (time.Now().Unix() + 5) * 1000
		signature := generateSignature(mc.apiSec, "GET", endpoint, nonce)

		authMessage := fmt.Sprintf(`{"op": "authKeyExpires", "args": ["%s", %d, "%s"]}`, mc.apiKey, nonce, signature)

		err = ws.WriteMessage(websocket.TextMessage, []byte(authMessage))
		if err != nil {
			mc.l.Fatalf("error sending bitmex auth message: %v", err)
		}
		//b.l.Printf("sent auth message: %s", authMessage)

		publicSubMessage := `{"op": "subscribe", "args": ["instrument:DERIVATIVES"]}`
		err = ws.WriteMessage(websocket.TextMessage, []byte(publicSubMessage))
		if err != nil {
			mc.l.Fatalf("error sending bitmex public sub message: %v", err)
		}
		// b.l.Printf("sent public sub message: %s", publicSubMessage)

		// subscribe to private channel
		privateSubMessage := fmt.Sprintf(`{"op": "subscribe", "args": ["%s" , "%s", "%s"]}`, "execution", "order", "margin")

		err = ws.WriteMessage(websocket.TextMessage, []byte(privateSubMessage))
		if err != nil {
			mc.l.Fatalf("error sending bitmex private sub message: %v", err)
		}
		//b.l.Printf("sent private sub message: %s", privateSubMessage)

		// Timer for heartbeat
		lastMessageTime := time.Now()
		pingTicker := time.NewTicker(pingInterval)
		defer pingTicker.Stop()

		// Track pong receipt
		pongReceived := make(chan struct{})

		ws.SetPongHandler(func(appData string) error {
			//b.l.Printf("received pong from BitMEX: %s", appData)
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

		go func() {
			defer close(done)
			for {
				_, message, err := ws.ReadMessage()
				if err != nil {
					mc.l.Printf("error reading bitmex message: %v", err)
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
				var orderTable swagger.OrderTable
				var marginTable swagger.MarginTable
				var executionTable swagger.ExecutionTable
				var instrumentTable swagger.InstrumentTable
				if err := json.Unmarshal(message, &baseMessage); err != nil {
					mc.l.Printf("error unmarshalling bitmex message: %v", err)
					continue
				}
				switch {
				case baseMessage.Info != "":
					mc.l.Printf("received message info: %s", baseMessage.Info)
				case baseMessage.Status != 0:
					mc.l.Printf("received message status: %d with error: %s", baseMessage.Status, baseMessage.Error)
				case baseMessage.Success:
					mc.l.Printf("received message success on %v", baseMessage.Subscribe)
				case baseMessage.Table != "":
					//b.l.Printf("received message table: %s", baseMessage.Table)
					switch baseMessage.Table {
					case "order":
						if err := json.Unmarshal(message, &orderTable); err == nil {
							mc.l.Printf("received message orderTable: %s", orderTable.Table)
							wsHandler(ctx, orderTable.Data)
							// for i := range orderTable.Data {
							// 	b.l.Printf("orderTable: %v", orderTable.Data[i])
							// }

						} else {
							mc.l.Printf("error unmarshalling bitmex message: %v", err)
						}
					case "margin":
						if err := json.Unmarshal(message, &marginTable); err == nil {
							//b.l.Printf("received message marginTable: %s", marginTable.Table)
						} else {
							mc.l.Printf("error unmarshalling bitmex message: %v", err)
						}
					case "execution":
						if err := json.Unmarshal(message, &executionTable); err == nil {
							mc.l.Printf("received message executionTable: %s", executionTable.Table)
						} else {
							mc.l.Printf("error unmarshalling bitmex message: %v", err)
						}
					case "instrument":
						if err := json.Unmarshal(message, &instrumentTable); err == nil {
							//b.l.Printf("received message instrumentTable, symbol is: %s data length: %d", instrumentTable.Table, len(instrumentTable.Data))
							if len(instrumentTable.Data) >= 3 {
								mc.l.Printf("too many symbols: %d", len(instrumentTable.Data))
							} else {
								for _, i := range instrumentTable.Data {
									// if i.Symbol == "SOLUSDT" {
									// 	b.l.Printf("symbol: %s, markPrice: %0.5f , since: %s", i.Symbol, i.MarkPrice, utils.FriendlyDuration(time.Since(i.Timestamp)))
									// }
									if i.MarkPrice > 0 {
										go mc.updateCandles(i.Symbol, i.MarkPrice, i.Timestamp)
									}
									//trunc1min := i.Timestamp.Truncate(time.Minute).Local()
									//trunc15min := i.Timestamp.Truncate(time.Minute * 15).Local()
									//trunc1h := i.Timestamp.Truncate(time.Hour).Local()

									//b.l.Printf("symbol: %s, markPrice: %0.5f , since: %s \ntruncated1min: %s, truncated15min: %s, truncated1h: %s", i.Symbol, i.MarkPrice, utils.FriendlyDuration(time.Since(i.Timestamp)), trunc1min, trunc15min, trunc1h)
								}

							}

						} else {
							mc.l.Printf("error unmarshalling bitmex message: %v", err)
						}
					default:
						mc.l.Printf("received unknown message table: %s", baseMessage.Table)

					}
				default:
					mc.l.Printf("received unknown message type: %s", message)
				}

			}
		}()

		// Ping/Pong handleing
		go func() {
			for {
				select {
				case <-pingTicker.C:
					//b.l.Printf("last message time: %s", utils.FriendlyDuration(time.Since(lastMessageTime)))
					if time.Since(lastMessageTime) >= pingInterval {
						//b.l.Println("ping bitmex")

						err := ws.WriteMessage(websocket.PingMessage, nil)
						if err != nil {
							mc.l.Printf("error sending bitmex ping: %v", err)
							return
						}
						// wait for pong
						pongWaitTimer := time.NewTimer(pongWait)
						select {
						case <-pongReceived:
							//b.l.Println("pong received in time")
							pongWaitTimer.Stop()
						case <-pongWaitTimer.C:
							mc.l.Println("pong timeout - reconnecting")
							ws.Close()
							return
						}
					}
				}
			}
		}()

		// wait for websocket to close or error out
		<-done
		mc.l.Println("Connection lost. Reconnecting...")
		time.Sleep(5 * time.Second)

	}
}

func generateSignature(apiSecret, verb, endpoint string, nonce int64) string {
	message := verb + endpoint + strconv.FormatInt(nonce, 10)
	hash := hmac.New(sha256.New, []byte(apiSecret))
	hash.Write([]byte(message))
	return hex.EncodeToString(hash.Sum(nil))
}
