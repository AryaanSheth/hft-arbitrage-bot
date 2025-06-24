package exchange

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"hft-arbitrage-bot/strategy"
)

type OKXSubscribe struct {
	Op   string            `json:"op"`
	Args []OKXSubscribeArg `json:"args"`
}

type OKXSubscribeArg struct {
	Channel string `json:"channel"`
	InstId  string `json:"instId"`
}

type OKXOrderBookMessage struct {
	Arg struct {
		Channel string `json:"channel"`
		InstId  string `json:"instId"`
	} `json:"arg"`
	Data []struct {
		Asks [][]string `json:"asks"` // [price, size, liquidity]
		Bids [][]string `json:"bids"`
		Ts   string     `json:"ts"`
	} `json:"data"`
}

// OKX starts the OKX WebSocket connection and sends quotes to the provided channel
func OKX(quoteChan chan<- strategy.Quote) {
	url := "wss://ws.okx.com:8443/ws/v5/public"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("WebSocket connection failed:", err)
	}
	defer conn.Close()

	subscribe := OKXSubscribe{
		Op: "subscribe",
		Args: []OKXSubscribeArg{
			{
				Channel: "books",
				InstId:  "BTC-USDT",
			},
		},
	}

	err = conn.WriteJSON(subscribe)
	if err != nil {
		log.Fatal("Subscription failed:", err)
	}

	log.Println("Subscribed to OKX BTC-USDT order book")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var msg OKXOrderBookMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		if len(msg.Data) == 0 {
			continue
		}
		ob := msg.Data[0]
		if len(ob.Bids) == 0 || len(ob.Asks) == 0 {
			continue
		}

		bidStr := ob.Bids[0][0]
		askStr := ob.Asks[0][0]

		bid, err1 := strconv.ParseFloat(bidStr, 64)
		ask, err2 := strconv.ParseFloat(askStr, 64)
		if err1 != nil || err2 != nil {
			log.Println("Error parsing bid/ask:", err1, err2)
			continue
		}

		quote := strategy.Quote{
			Exchange:  "okx",
			Symbol:    "BTCUSDT",
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

		log.Printf("OKX: Bid=%.2f, Ask=%.2f", bid, ask)
	}
}
