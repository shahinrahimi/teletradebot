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
	"github.com/shahinrahimi/teletradebot/utils"
)

const (
	wsBitmexURL          = "wss://www.bitmex.com"
	wsBitmexURLTestnet   = "wss://testnet.bitmex.com"
	endpoint             = "/realtime"
	Symbol               = "ETHUSDT"
	pingInterval         = 5 * time.Second
	pongWait             = 5 * time.Second
	reconnectingInterval = 5 * time.Second
)

// WebSocketConnection handles the lifecycle of the WebSocket connection
type WebSocketConnection struct {
	ws      *websocket.Conn
	ctx     context.Context
	handler func(ctx context.Context, od []swagger.OrderData)
}

// StartWebsocketService initializes the WebSocket connection and handles reconnection
func (mc *BitmexClient) StartWebsocketService(ctx context.Context, wsHandler func(ctx2 context.Context, od []swagger.OrderData)) {
	go mc.startUserDataStream(ctx, wsHandler)
}

// startUserDataStream initiates and manages the connection loop
func (mc *BitmexClient) startUserDataStream(ctx context.Context, wsHandler func(ctx2 context.Context, od []swagger.OrderData)) {
	for {
		err := mc.connectWebSocket(ctx, wsHandler)
		if err != nil {
			mc.l.Printf("Error in websocket connection: %v. Reconnecting in 5 seconds...", err)
			time.Sleep(5 * time.Second)
		}
	}
}

// connectWebSocket manages the WebSocket connection and its events
func (mc *BitmexClient) connectWebSocket(ctx context.Context, wsHandler func(ctx2 context.Context, od []swagger.OrderData)) error {
	websocketURL := wsBitmexURL + endpoint
	if config.UseBitmexTestnet {
		websocketURL = wsBitmexURLTestnet + endpoint
	}

	// Establish WebSocket connection
	ws, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		return fmt.Errorf("error dialing websocket: %v", err)
	}
	defer ws.Close() // Manually close on error or completion

	mc.l.Println("Connected to BitMEX WebSocket.")

	conn := &WebSocketConnection{ws: ws, ctx: ctx, handler: wsHandler}

	// Authenticate and subscribe
	if err := mc.authenticateAndSubscribe(conn); err != nil {
		return fmt.Errorf("failed to authenticate or subscribe: %v", err)
	}

	// Start ping/pong handling and read messages
	pingTicker := time.NewTicker(pingInterval)
	defer pingTicker.Stop()

	done := make(chan struct{})

	go conn.handleIncomingMessages(done)

	go conn.handlePingPong(pingTicker, done)

	// Wait for WebSocket close or error
	<-done
	mc.l.Printf("Bitmex WebSocket closed, Reconnecting after %s", utils.FriendlyDuration(reconnectingInterval))

	return nil
}

// authenticateAndSubscribe handles authentication and subscription to WebSocket streams
func (mc *BitmexClient) authenticateAndSubscribe(conn *WebSocketConnection) error {
	nonce := (time.Now().Unix() + 5) * 1000
	signature := generateSignature(mc.apiSec, "GET", endpoint, nonce)

	authMessage := fmt.Sprintf(`{"op": "authKeyExpires", "args": ["%s", %d, "%s"]}`, mc.apiKey, nonce, signature)
	if err := conn.ws.WriteMessage(websocket.TextMessage, []byte(authMessage)); err != nil {
		return err
	}

	publicSubMessage := `{"op": "subscribe", "args": ["instrument:DERIVATIVES"]}`
	privateSubMessage := `{"op": "subscribe", "args": ["execution", "order", "margin"]}`

	if err := conn.ws.WriteMessage(websocket.TextMessage, []byte(publicSubMessage)); err != nil {
		return err
	}
	if err := conn.ws.WriteMessage(websocket.TextMessage, []byte(privateSubMessage)); err != nil {
		return err
	}

	return nil
}

// handleIncomingMessages reads messages from WebSocket and processes them
func (conn *WebSocketConnection) handleIncomingMessages(done chan struct{}) {
	defer close(done)

	for {
		_, message, err := conn.ws.ReadMessage()
		if err != nil {
			conn.ws.Close()
			conn.logError("error reading websocket message", err)
			return
		}

		conn.processMessage(message)
	}
}

// processMessage unmarshals and processes incoming WebSocket messages
func (conn *WebSocketConnection) processMessage(message []byte) {
	var baseMessage struct {
		Info      string `json:"info"`
		Status    int    `json:"status"`
		Error     string `json:"error"`
		Success   bool   `json:"success"`
		Subscribe string `json:"subscribe"`
		Table     string `json:"table"`
	}

	if err := json.Unmarshal(message, &baseMessage); err != nil {
		conn.logError("error unmarshalling base message", err)
		return
	}

	switch baseMessage.Table {
	case "order":
		conn.handleOrderMessage(message)
	// Handle other cases (margin, execution, etc.)
	default:
		conn.logInfo("unknown message type", string(message))
	}
}

// handleOrderMessage processes order messages
func (conn *WebSocketConnection) handleOrderMessage(message []byte) {
	var orderTable swagger.OrderTable
	if err := json.Unmarshal(message, &orderTable); err == nil {
		conn.handler(conn.ctx, orderTable.Data)
	} else {
		conn.logError("error unmarshalling order message", err)
	}
}

// handlePingPong manages ping/pong mechanism for WebSocket
func (conn *WebSocketConnection) handlePingPong(ticker *time.Ticker, done chan struct{}) {
	pongReceived := make(chan struct{})

	conn.ws.SetPongHandler(func(appData string) error {
		select {
		case pongReceived <- struct{}{}:
		default:
		}
		return nil
	})

	for {
		select {
		case <-ticker.C:
			if err := conn.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				conn.logError("error sending ping", err)
				return
			}

			select {
			case <-pongReceived:
			case <-time.After(pongWait):
				conn.logInfo("pong not received in time, closing connection", "")
				conn.ws.Close()
				return
			}
		case <-done:
			return
		}
	}
}

// Utility functions for logging
func (conn *WebSocketConnection) logError(message string, err error) {
	fmt.Printf("%s: %v\n", message, err)
}

func (conn *WebSocketConnection) logInfo(message, data string) {
	//fmt.Printf("%s: %s\n", message, data)
}

// generateSignature creates the BitMEX API signature
func generateSignature(apiSecret, verb, endpoint string, nonce int64) string {
	message := verb + endpoint + strconv.FormatInt(nonce, 10)
	hash := hmac.New(sha256.New, []byte(apiSecret))
	hash.Write([]byte(message))
	return hex.EncodeToString(hash.Sum(nil))
}
