package exchange

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"hft-arbitrage-bot/strategy"
)

type KuCoinMessage struct {
	Type string `json:"type"`
	Topic string `json:"topic"`
	Subject string `json:"subject"`
	Data struct {
		BestBid string `json:"bestBid"`
		BestAsk string `json:"bestAsk"`
		Symbol string `json:"symbol"`
	} `json:"data"`
}

func Kucoin(quoteChan chan<- strategy.Quote) {
	url := "wss://ws-api-spot.kucoin.com/endpoint"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to KuCoin WebSocket:", err)
	}
	defer conn.Close()

	subMsg := map[string]interface{}{
		"id":   "dogeusdt-arb",
		"type": "subscribe",
		"topic": "/market/ticker:DOGE-USDT",
		"privateChannel": false,
		"response": true,
	}
	if err := conn.WriteJSON(subMsg); err != nil {
		log.Fatal("KuCoin subscription failed:", err)
	}
	log.Println("🟢 Subscribed to KuCoin DOGE-USDT ticker")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("KuCoin read error:", err)
			break
		}
		var msg KuCoinMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}
		if msg.Type != "message" || msg.Subject != "trade.ticker" {
			continue
		}
		bid, err1 := strconv.ParseFloat(msg.Data.BestBid, 64)
		ask, err2 := strconv.ParseFloat(msg.Data.BestAsk, 64)
		if err1 != nil || err2 != nil {
			continue
		}
		quote := strategy.Quote{
			Exchange:  "kucoin",
			Symbol:    "DOGEUSDT",
			Bid:       bid,
			Ask:       ask,
			Timestamp: time.Now(),
		}
		select {
		case quoteChan <- quote:
		default:
		}
		log.Printf("🟢 KuCoin: Bid=%.6f, Ask=%.6f", bid, ask)
	}
} 