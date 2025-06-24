package exchange

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"hft-arbitrage-bot/strategy"
)

type BinanceBookTicker struct {
	Symbol   string `json:"s"`
	BidPrice string `json:"b"`
	BidQty   string `json:"B"`
	AskPrice string `json:"a"`
	AskQty   string `json:"A"`
}

// Binance starts the Binance WebSocket connection and sends quotes to the provided channel
func Binance(quoteChan chan<- strategy.Quote) {
	url := "wss://stream.binance.com:9443/ws/dogeusdt@bookTicker"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer conn.Close()

	log.Println("Connected to Binance stream for dogeusdt")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var ticker BinanceBookTicker
		err = json.Unmarshal(message, &ticker)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		bid, err1 := strconv.ParseFloat(ticker.BidPrice, 64)
		ask, err2 := strconv.ParseFloat(ticker.AskPrice, 64)
		if err1 != nil || err2 != nil {
			log.Println("Error parsing bid/ask:", err1, err2)
			continue
		}

		quote := strategy.Quote{
			Exchange:  "binance",
			Symbol:    ticker.Symbol,
			Bid:       bid,
			Ask:       ask,
			Timestamp: time.Now(),
		}

		// Send quote to strategy
		select {
		case quoteChan <- quote:
		default:
			// Channel is full, skip this quote
		}

		log.Printf("🟡 Binance: Bid=%.6f, Ask=%.6f", bid, ask)
	}
}
