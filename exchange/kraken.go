package exchange

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"hft-arbitrage-bot/strategy"
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

// Kraken starts the Kraken WebSocket connection and sends quotes to the provided channel
func Kraken(quoteChan chan<- strategy.Quote) {
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
			Depth: 10, // more depth = more data but slower
		},
	}

	err = conn.WriteJSON(subscribe)
	if err != nil {
		log.Fatal("Failed to subscribe:", err)
	}

	log.Println("Subscribed to Kraken BTC/USD order book")

	var currentBid, currentAsk float64

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
				if askFloat, err := strconv.ParseFloat(ask, 64); err == nil {
					currentAsk = askFloat
				}
			}
			if bids, ok := payload["b"]; ok {
				bid := bids.([]any)[0].([]any)[0].(string)
				if bidFloat, err := strconv.ParseFloat(bid, 64); err == nil {
					currentBid = bidFloat
				}
			}

			// Send quote if we have both bid and ask
			if currentBid > 0 && currentAsk > 0 {
				quote := strategy.Quote{
					Exchange:  "kraken",
					Symbol:    "XBTUSD",
					Bid:       currentBid,
					Ask:       currentAsk,
					Timestamp: time.Now(),
				}

				// Send quote to strategy
				select {
				case quoteChan <- quote:
				default:
					// Channel is full, skip this quote
				}

				log.Printf("Kraken: Bid=%.2f, Ask=%.2f", currentBid, currentAsk)
			}
		}
	}
}
