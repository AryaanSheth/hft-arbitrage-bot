package exchange

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
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

func okx() {
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

		bid := ob.Bids[0][0]
		ask := ob.Asks[0][0]

		log.Printf("OKX BTC/USDT - Bid: %s, Ask: %s, Time: %s\n", bid, ask, time.Now().Format(time.RFC3339))
	}
}
