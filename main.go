package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type CoinbaseAuthSubscription struct {
	Type       string   `json:"type"`
	Channels   []string `json:"channels"`
	Key        string   `json:"key"`
	Passphrase string   `json:"passphrase"`
	Timestamp  string   `json:"timestamp"`
	Signature  string   `json:"signature"`
	ProductIDs []string `json:"product_ids"`
}

type CoinbaseSnapshot struct {
	Type      string     `json:"type"`
	ProductID string     `json:"product_id"`
	Bids      [][]string `json:"bids"`
	Asks      [][]string `json:"asks"`
}

type CoinbaseL2Update struct {
	Type      string     `json:"type"`
	ProductID string     `json:"product_id"`
	Changes   [][]string `json:"changes"`
}

type Quote struct {
	Exchange  string
	Symbol    string
	Bid       float64
	Ask       float64
	Timestamp time.Time
}

func generateSignature(secret, timestamp, method, requestPath, body string) string {
	message := timestamp + method + requestPath + body
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func main() {
	// Load from env or hardcode (do NOT commit real secrets)
	apiKey := os.Getenv("COINBASE_API_KEY")
	apiSecret := os.Getenv("COINBASE_API_SECRET")
	passphrase := os.Getenv("COINBASE_PASSPHRASE")

	if apiKey == "" || apiSecret == "" || passphrase == "" {
		log.Fatal("Missing Coinbase API credentials")
	}

	url := "wss://ws-feed.exchange.coinbase.com"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to Coinbase WebSocket:", err)
	}
	defer conn.Close()

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	method := "GET"
	requestPath := "/users/self/verify"
	body := ""

	signature := generateSignature(apiSecret, timestamp, method, requestPath, body)

	subscribe := CoinbaseAuthSubscription{
		Type:       "subscribe",
		Channels:   []string{"level2"},
		Key:        apiKey,
		Passphrase: passphrase,
		Timestamp:  timestamp,
		Signature:  signature,
		ProductIDs: []string{"BTC-USD"},
	}

	err = conn.WriteJSON(subscribe)
	if err != nil {
		log.Fatal("Error subscribing:", err)
	}

	log.Println("Subscribed to Coinbase BTC-USD level2 (authenticated)")

	var bestBid, bestAsk float64

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		log.Printf("RAW: %s\n", message)

		var generic map[string]any
		err = json.Unmarshal(message, &generic)
		if err != nil {
			log.Println("Error parsing generic message:", err)
			continue
		}

		switch generic["type"] {
		case "snapshot":
			var snapshot CoinbaseSnapshot
			err := json.Unmarshal(message, &snapshot)
			if err != nil {
				log.Println("Error parsing snapshot:", err)
				continue
			}
			if len(snapshot.Bids) > 0 {
				bestBid, _ = strconv.ParseFloat(snapshot.Bids[0][0], 64)
			}
			if len(snapshot.Asks) > 0 {
				bestAsk, _ = strconv.ParseFloat(snapshot.Asks[0][0], 64)
			}

		case "l2update":
			var update CoinbaseL2Update
			err := json.Unmarshal(message, &update)
			if err != nil {
				log.Println("Error parsing l2update:", err)
				continue
			}
			for _, change := range update.Changes {
				side := change[0]
				price, _ := strconv.ParseFloat(change[1], 64)
				size, _ := strconv.ParseFloat(change[2], 64)

				if side == "buy" && size > 0 && price > bestBid {
					bestBid = price
				}
				if side == "sell" && size > 0 && (bestAsk == 0 || price < bestAsk) {
					bestAsk = price
				}
			}
		}

		if bestBid > 0 && bestAsk > 0 {
			quote := Quote{
				Exchange:  "coinbase",
				Symbol:    "BTC-USD",
				Bid:       bestBid,
				Ask:       bestAsk,
				Timestamp: time.Now(),
			}
			log.Printf("Received quote: %+v\n", quote)
		}
	}
}
