package exchange

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type BinanceBookTicker struct {
	Symbol   string `json:"s"`
	BidPrice string `json:"b"`
	BidQty   string `json:"B"`
	AskPrice string `json:"a"`
	AskQty   string `json:"A"`
}

type Quote struct {
	Exchange string
	Symbol   string
	Bid      float64
	Ask      float64
	Timestamp time.Time
}

func main() {
	url := "wss://stream.binance.com:9443/ws/btcusdt@bookTicker"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer conn.Close()

	log.Println("Connected to Binance stream for btcusdt")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// log.Println("RAW:", string(message))

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

		quote := Quote{
			Exchange:  "binance",
			Symbol:    ticker.Symbol,
			Bid:       bid,
			Ask:       ask,
			Timestamp: time.Now(),
		}

		log.Printf("Received quote: %+v\n", quote)
	}
}
