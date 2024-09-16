package bot

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wsBitmexURL        string = "wss://www.bitmex.com/realtime"
	wsBitmexURLTestnet string = "wss://testnet.bitmex.com/realtime"
	Symbol             string = "ETHUSDT"
)

func (b *Bot) StartWebsocketServiceBitmex(ctx context.Context) {
	go b.startUserDataStreamBitmex(ctx)
}

func generateSignature(apiSecret, verb, endpoint string, nonce int64) string {
	message := verb + endpoint + strconv.FormatInt(nonce, 10)
	hash := hmac.New(sha256.New, []byte(apiSecret))
	hash.Write([]byte(message))
	return hex.EncodeToString(hash.Sum(nil))
}

func (b *Bot) startUserDataStreamBitmex(ctx context.Context) {
	ws, _, err := websocket.DefaultDialer.Dial(wsBitmexURLTestnet, nil)
	if err != nil {
		b.l.Fatal("error dialing bitmex websocket: %v", err)
	}

	defer ws.Close()

	nonce := time.Now().UnixMilli()
	signature := generateSignature(b.mc.ApiSec, "GET", "/api/v1/userDataStream", nonce)

	authMessage := fmt.Sprintf(`{"op": "authKeyExpires", "args": ["%s", %d, %s]}`, b.mc.ApiKey, nonce, signature)

	err = ws.WriteMessage(websocket.TextMessage, []byte(authMessage))
	if err != nil {
		b.l.Fatalf("error sending bitmex auth message: %v", err)
	}
	b.l.Printf("sent auth message: %s", authMessage)

	// subscrie to public channels
	publicSubMessage := fmt.Sprintf(`{"op": "subscribe", ["trade:%s", "instrument:%s"]}`, Symbol, Symbol)

	err = ws.WriteMessage(websocket.TextMessage, []byte(publicSubMessage))
	if err != nil {
		b.l.Fatalf("error sending bitmex public sub message: %v", err)
	}
	b.l.Printf("sent public sub message: %s", publicSubMessage)

	// subscribe to private channel
	privateSubMessage := fmt.Sprintf(`{"op": "subscribe", "args": ["%s" , "%s", "%s"]}`, "user", "order", "margin")

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
