package exchange

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"hft-arbitrage-bot/strategy"
)

type BybitBookTicker struct {
	Topic string `json:"topic"`
	Type  string `json:"type"`
	Data  []struct {
		Symbol string `json:"s"`
		Bid1Price string `json:"b"`
		Ask1Price string `json:"a"`
		Time int64 `json:"T"`
	} `json:"data"`
}

func Bybit(quoteChan chan<- strategy.Quote) {
	url := "wss://stream.bybit.com/v5/public/spot"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to Bybit WebSocket:", err)
	}
	defer conn.Close()

	subMsg := map[string]interface{}{
		"op": "subscribe",
		"args": []string{"tickers.DOGEUSDT"},
	}
	if err := conn.WriteJSON(subMsg); err != nil {
		log.Fatal("Bybit subscription failed:", err)
	}
	log.Println("🟠 Subscribed to Bybit DOGEUSDT ticker")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Bybit read error:", err)
			break
		}

		var ticker BybitBookTicker
		if err := json.Unmarshal(message, &ticker); err != nil {
			continue
		}
		if len(ticker.Data) == 0 {
			continue
		}
		bid, err1 := strconv.ParseFloat(ticker.Data[0].Bid1Price, 64)
		ask, err2 := strconv.ParseFloat(ticker.Data[0].Ask1Price, 64)
		if err1 != nil || err2 != nil {
			continue
		}
		quote := strategy.Quote{
			Exchange:  "bybit",
			Symbol:    "DOGEUSDT",
			Bid:       bid,
			Ask:       ask,
			Timestamp: time.Now(),
		}
		select {
		case quoteChan <- quote:
		default:
		}
		log.Printf("🟠 Bybit: Bid=%.6f, Ask=%.6f", bid, ask)
	}
} 