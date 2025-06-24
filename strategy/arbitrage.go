package strategy

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Quote represents a price quote from an exchange
type Quote struct {
	Exchange  string
	Symbol    string
	Bid       float64
	Ask       float64
	Timestamp time.Time
}

// ArbitrageOpportunity represents a potential arbitrage opportunity
type ArbitrageOpportunity struct {
	BuyExchange  string
	SellExchange string
	Symbol       string
	BuyPrice     float64
	SellPrice    float64
	Spread       float64
	SpreadPercent float64
	Timestamp    time.Time
}

// ArbitrageStrategy manages the arbitrage detection across multiple exchanges
type ArbitrageStrategy struct {
	quotes     map[string]Quote
	quotesLock sync.RWMutex
	minSpread  float64 // minimum spread percentage to consider arbitrage
	pnlManager *PnLManager // Add P&L manager
}

// NewArbitrageStrategy creates a new arbitrage strategy instance
func NewArbitrageStrategy(minSpreadPercent float64, initialBalance, tradeSize float64) *ArbitrageStrategy {
	return &ArbitrageStrategy{
		quotes:    make(map[string]Quote),
		minSpread: minSpreadPercent,
		pnlManager: NewPnLManager(initialBalance, tradeSize),
	}
}

// UpdateQuote updates the latest quote for an exchange
func (as *ArbitrageStrategy) UpdateQuote(quote Quote) {
	as.quotesLock.Lock()
	defer as.quotesLock.Unlock()
	
	as.quotes[quote.Exchange] = quote
}

// FindArbitrageOpportunities analyzes current quotes and finds arbitrage opportunities
func (as *ArbitrageStrategy) FindArbitrageOpportunities() []ArbitrageOpportunity {
	as.quotesLock.RLock()
	defer as.quotesLock.RUnlock()

	var opportunities []ArbitrageOpportunity
	exchanges := make([]string, 0, len(as.quotes))
	
	// Collect all exchanges with valid quotes
	for exchange, quote := range as.quotes {
		if quote.Bid > 0 && quote.Ask > 0 {
			exchanges = append(exchanges, exchange)
		}
	}

	// Need at least 2 exchanges to find arbitrage
	if len(exchanges) < 2 {
		return opportunities
	}

	// Compare all pairs of exchanges
	for i := 0; i < len(exchanges); i++ {
		for j := i + 1; j < len(exchanges); j++ {
			exchange1 := exchanges[i]
			exchange2 := exchanges[j]
			
			quote1 := as.quotes[exchange1]
			quote2 := as.quotes[exchange2]

			// Check if we can buy on exchange1 and sell on exchange2
			if quote1.Ask < quote2.Bid {
				spread := quote2.Bid - quote1.Ask
				spreadPercent := (spread / quote1.Ask) * 100
				
				if spreadPercent >= as.minSpread {
					opportunities = append(opportunities, ArbitrageOpportunity{
						BuyExchange:   exchange1,
						SellExchange:  exchange2,
						Symbol:        quote1.Symbol,
						BuyPrice:      quote1.Ask,
						SellPrice:     quote2.Bid,
						Spread:        spread,
						SpreadPercent: spreadPercent,
						Timestamp:     time.Now(),
					})
				}
			}

			// Check if we can buy on exchange2 and sell on exchange1
			if quote2.Ask < quote1.Bid {
				spread := quote1.Bid - quote2.Ask
				spreadPercent := (spread / quote2.Ask) * 100
				
				if spreadPercent >= as.minSpread {
					opportunities = append(opportunities, ArbitrageOpportunity{
						BuyExchange:   exchange2,
						SellExchange:  exchange1,
						Symbol:        quote2.Symbol,
						BuyPrice:      quote2.Ask,
						SellPrice:     quote1.Bid,
						Spread:        spread,
						SpreadPercent: spreadPercent,
						Timestamp:     time.Now(),
					})
				}
			}
		}
	}

	return opportunities
}

// PrintOpportunities prints arbitrage opportunities in a formatted way
func (as *ArbitrageStrategy) PrintOpportunities(opportunities []ArbitrageOpportunity) {
	if len(opportunities) == 0 {
		return
	}

	log.Println("=== ARBITRAGE OPPORTUNITIES ===")
	for _, opp := range opportunities {
		log.Printf("💰 BUY on %s at %.2f, SELL on %s at %.2f", 
			opp.BuyExchange, opp.BuyPrice, opp.SellExchange, opp.SellPrice)
		log.Printf("   Spread: $%.2f (%.2f%%)", opp.Spread, opp.SpreadPercent)
		log.Printf("   Time: %s", opp.Timestamp.Format("15:04:05.000"))
		
		// Execute the arbitrage opportunity
		err := as.pnlManager.ExecuteArbitrage(opp)
		if err != nil {
			log.Printf("❌ Failed to execute arbitrage: %v", err)
		}
		
		log.Println("---")
	}
}

// GetQuoteSummary returns a summary of all current quotes
func (as *ArbitrageStrategy) GetQuoteSummary() string {
	as.quotesLock.RLock()
	defer as.quotesLock.RUnlock()

	summary := "Current Quotes:\n"
	for exchange, quote := range as.quotes {
		summary += fmt.Sprintf("  %s: Bid=%.2f, Ask=%.2f, Spread=%.2f%%\n", 
			exchange, quote.Bid, quote.Ask, 
			((quote.Ask-quote.Bid)/quote.Bid)*100)
	}
	return summary
}

// RunArbitrageStrategy runs the main arbitrage strategy loop
func (as *ArbitrageStrategy) RunArbitrageStrategy(quoteChan <-chan Quote) {
	log.Println("Starting arbitrage strategy...")
	
	ticker := time.NewTicker(100 * time.Millisecond) // Check every 100ms
	pnlTicker := time.NewTicker(5 * time.Second) // Print P&L every 5 seconds
	defer ticker.Stop()
	defer pnlTicker.Stop()

	for {
		select {
		case quote := <-quoteChan:
			as.UpdateQuote(quote)
			
		case <-ticker.C:
			opportunities := as.FindArbitrageOpportunities()
			if len(opportunities) > 0 {
				as.PrintOpportunities(opportunities)
			}
			
		case <-pnlTicker.C:
			// Print P&L status periodically
			as.pnlManager.PrintPnLStatus()
		}
	}
}

// GetPnLManager returns the P&L manager for external access
func (as *ArbitrageStrategy) GetPnLManager() *PnLManager {
	return as.pnlManager
} 