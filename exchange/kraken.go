package exchange

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type KrakenSubscribeMsg struct {
	Event        string       `json:"event"`
	Pair         []string     `json:"pair"`
	Subscription Subscription `json:"subscription"`
}

type Subscription struct {
	Name  string `json:"name"`
	Depth int    `json:"depth"`
}

func Kraken() {
	url := "wss://ws.kraken.com"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("WebSocket connection failed:", err)
	}
	defer conn.Close()

	subscribe := KrakenSubscribeMsg{
		Event: "subscribe",
		Pair:  []string{"XBT/USD"},
		Subscription: Subscription{
			Name:  "book",
			Depth: 10, // more detpth = more data but slower
		},
	}

	err = conn.WriteJSON(subscribe)
	if err != nil {
		log.Fatal("Failed to subscribe:", err)
	}

	log.Println("Subscribed to Kraken BTC/USD order book")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var data []any
		if err := json.Unmarshal(message, &data); err != nil {
			continue
		}

		if len(data) >= 2 {
			payload := data[1].(map[string]any)

			// Check for asks or bids
			if asks, ok := payload["a"]; ok {
				ask := asks.([]any)[0].([]any)[0].(string)
				log.Printf("Ask: %s at %s\n", ask, time.Now().Format(time.RFC3339))
			}
			if bids, ok := payload["b"]; ok {
				bid := bids.([]any)[0].([]any)[0].(string)
				log.Printf("Bid: %s at %s\n", bid, time.Now().Format(time.RFC3339))
			}
		}
	}
}
