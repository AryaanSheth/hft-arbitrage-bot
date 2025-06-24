package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type PnLResponse struct {
	Status    string      `json:"status"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: pnl_client <command> [options]")
		fmt.Println("Commands:")
		fmt.Println("  pnl     - Get full P&L status")
		fmt.Println("  summary - Get P&L summary")
		fmt.Println("  trades  - Get recent trades")
		fmt.Println("  health  - Check API health")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("  --host <host> - API host (default: localhost:8080)")
		fmt.Println("  --limit <n>   - Number of trades to fetch (for trades command)")
		os.Exit(1)
	}

	command := os.Args[1]
	host := "localhost:8080"
	limit := "10"

	// Parse options
	for i := 2; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--host":
			if i+1 < len(os.Args) {
				host = os.Args[i+1]
				i++
			}
		case "--limit":
			if i+1 < len(os.Args) {
				limit = os.Args[i+1]
				i++
			}
		}
	}

	url := fmt.Sprintf("http://%s", host)

	switch command {
	case "pnl":
		getPnL(url)
	case "summary":
		getSummary(url)
	case "trades":
		getTrades(url, limit)
	case "health":
		getHealth(url)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func getPnL(url string) {
	resp, err := http.Get(url + "/pnl")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	var response PnLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	if response.Status != "success" {
		fmt.Printf("API Error: %v\n", response.Data)
		os.Exit(1)
	}

	// Pretty print the P&L data
	data, _ := json.MarshalIndent(response.Data, "", "  ")
	fmt.Printf("=== P&L STATUS ===\n")
	fmt.Printf("Timestamp: %s\n", time.Unix(response.Timestamp, 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("%s\n", string(data))
}

func getSummary(url string) {
	resp, err := http.Get(url + "/summary")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	var response PnLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	if response.Status != "success" {
		fmt.Printf("API Error: %v\n", response.Data)
		os.Exit(1)
	}

	// Pretty print the summary data
	data, _ := json.MarshalIndent(response.Data, "", "  ")
	fmt.Printf("=== P&L SUMMARY ===\n")
	fmt.Printf("Timestamp: %s\n", time.Unix(response.Timestamp, 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("%s\n", string(data))
}

func getTrades(url string, limit string) {
	resp, err := http.Get(fmt.Sprintf("%s/trades?limit=%s", url, limit))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	var response PnLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	if response.Status != "success" {
		fmt.Printf("API Error: %v\n", response.Data)
		os.Exit(1)
	}

	// Pretty print the trades data
	data, _ := json.MarshalIndent(response.Data, "", "  ")
	fmt.Printf("=== RECENT TRADES ===\n")
	fmt.Printf("Timestamp: %s\n", time.Unix(response.Timestamp, 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("%s\n", string(data))
}

func getHealth(url string) {
	resp, err := http.Get(url + "/health")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	var response PnLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("=== API HEALTH ===\n")
	fmt.Printf("Status: %s\n", response.Data.(map[string]interface{})["status"])
	fmt.Printf("Timestamp: %s\n", time.Unix(response.Timestamp, 0).Format("2006-01-02 15:04:05"))
} 