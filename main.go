package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"hft-arbitrage-bot/api"
	"hft-arbitrage-bot/exchange"
	"hft-arbitrage-bot/strategy"
)

func main() {
	log.Println("🚀 Starting HFT Arbitrage Bot")

	// Create a channel for quotes from all exchanges
	quoteChan := make(chan strategy.Quote, 1000) // Buffered channel to handle high-frequency updates

	// Create arbitrage strategy with 0.1% minimum spread, $1000 initial balance, $100 trade size
	arbitrageStrategy := strategy.NewArbitrageStrategy(0.1, 1000.0, 100.0)

	// Start the arbitrage strategy in a goroutine
	go arbitrageStrategy.RunArbitrageStrategy(quoteChan)

	// Start P&L API server
	pnlAPI := api.NewPnLAPI(arbitrageStrategy.GetPnLManager(), 8080)
	pnlAPI.Start()

	// Start all exchanges in separate goroutines
	var wg sync.WaitGroup

	// Start Binance
	wg.Add(1)
	go func() {
		defer wg.Done()
		exchange.Binance(quoteChan)
	}()

	// Start Kraken
	wg.Add(1)
	go func() {
		defer wg.Done()
		exchange.Kraken(quoteChan)
	}()

	// Start OKX
	wg.Add(1)
	go func() {
		defer wg.Done()
		exchange.OKX(quoteChan)
	}()

	log.Println("✅ All exchanges started successfully")
	log.Println("📊 Monitoring for arbitrage opportunities...")
	log.Println("💡 Minimum spread threshold: 0.1%")
	log.Println("💰 Initial balance: $1000.00")
	log.Println("📈 Trade size: $100.00")
	log.Println("🌐 P&L API available at http://localhost:8080")
	log.Println("")
	log.Println("💡 Commands:")
	log.Println("   - Press Enter to check P&L status")
	log.Println("   - Press Ctrl+C to stop the bot")
	log.Println("")
	log.Println("🌐 API Endpoints:")
	log.Println("   - GET /pnl - Full P&L status")
	log.Println("   - GET /summary - P&L summary")
	log.Println("   - GET /trades - Recent trades")
	log.Println("   - GET /health - Health check")

	// Start a goroutine to handle user input for P&L checking
	go handleUserInput(arbitrageStrategy)

	// Wait for interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("🛑 Shutting down HFT Arbitrage Bot...")
	
	// Stop the API server
	pnlAPI.Stop()
	
	// Close the quote channel to stop the strategy
	close(quoteChan)
	
	// Wait for all goroutines to finish
	wg.Wait()
	
	// Print final P&L status
	log.Println("")
	log.Println("=== FINAL P&L REPORT ===")
	arbitrageStrategy.GetPnLManager().PrintPnLStatus()
	
	log.Println("✅ HFT Arbitrage Bot stopped successfully")
}

// handleUserInput handles user input for checking P&L status
func handleUserInput(arbitrageStrategy *strategy.ArbitrageStrategy) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		
		switch input {
		case "":
			// Just Enter pressed - show P&L status
			arbitrageStrategy.GetPnLManager().PrintPnLStatus()
			
		case "pnl", "P&L", "status":
			// Explicit P&L request
			arbitrageStrategy.GetPnLManager().PrintPnLStatus()
			
		case "trades", "history":
			// Show recent trade history
			trades := arbitrageStrategy.GetPnLManager().GetTradeHistory(10)
			log.Println("=== RECENT TRADES ===")
			for i, trade := range trades {
				log.Printf("%d. %s %s %.4f %s at $%.2f on %s", 
					i+1, trade.Type, trade.Symbol, trade.Quantity, trade.Exchange, trade.Price, trade.Timestamp.Format("15:04:05"))
			}
			log.Println("====================")
			
		case "help":
			log.Println("Available commands:")
			log.Println("  Enter - Check P&L status")
			log.Println("  pnl   - Check P&L status")
			log.Println("  trades - Show recent trade history")
			log.Println("  help  - Show this help")
			
		default:
			log.Printf("Unknown command: %s. Type 'help' for available commands.", input)
		}
	}
}
