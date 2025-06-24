package strategy

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Trade represents an executed trade
type Trade struct {
	ID          string
	Type        string // "BUY" or "SELL"
	Exchange    string
	Symbol      string
	Price       float64
	Quantity    float64
	Timestamp   time.Time
	OrderID     string
	Status      string // "PENDING", "FILLED", "CANCELLED", "FAILED"
}

// Position represents a current position in a symbol
type Position struct {
	Symbol    string
	Quantity  float64
	AvgPrice  float64
	PnL       float64
	UnrealizedPnL float64
	LastUpdate time.Time
}

// PnLManager manages profit/loss tracking and trade execution
type PnLManager struct {
	trades     []Trade
	positions  map[string]*Position
	balance    float64
	initialBalance float64
	mutex      sync.RWMutex
	
	// Configuration
	baseBalance    float64
	tradeSize      float64
	maxPositions   int
	
	// Statistics
	totalTrades    int
	winningTrades  int
	losingTrades   int
	totalPnL       float64
	largestWin     float64
	largestLoss    float64
}

// NewPnLManager creates a new P&L manager
func NewPnLManager(initialBalance, tradeSize float64) *PnLManager {
	return &PnLManager{
		trades:         make([]Trade, 0),
		positions:      make(map[string]*Position),
		balance:        initialBalance,
		initialBalance: initialBalance,
		baseBalance:    initialBalance,
		tradeSize:      tradeSize,
		maxPositions:   5,
	}
}

// ExecuteArbitrage executes an arbitrage opportunity
func (pm *PnLManager) ExecuteArbitrage(opp ArbitrageOpportunity) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Check if we have enough balance
	if pm.balance < pm.tradeSize {
		return fmt.Errorf("insufficient balance: %.2f < %.2f", pm.balance, pm.tradeSize)
	}

	// Calculate quantity based on trade size
	quantity := pm.tradeSize / opp.BuyPrice

	// Execute buy trade
	buyTrade := Trade{
		ID:        fmt.Sprintf("buy_%d", time.Now().UnixNano()),
		Type:      "BUY",
		Exchange:  opp.BuyExchange,
		Symbol:    opp.Symbol,
		Price:     opp.BuyPrice,
		Quantity:  quantity,
		Timestamp: time.Now(),
		Status:    "FILLED",
	}

	// Execute sell trade
	sellTrade := Trade{
		ID:        fmt.Sprintf("sell_%d", time.Now().UnixNano()),
		Type:      "SELL",
		Exchange:  opp.SellExchange,
		Symbol:    opp.Symbol,
		Price:     opp.SellPrice,
		Quantity:  quantity,
		Timestamp: time.Now(),
		Status:    "FILLED",
	}

	// Update balance and positions
	pm.balance -= buyTrade.Price * buyTrade.Quantity
	pm.balance += sellTrade.Price * sellTrade.Quantity

	// Calculate P&L for this arbitrage
	pnl := (sellTrade.Price - buyTrade.Price) * quantity
	pm.totalPnL += pnl

	// Update statistics
	pm.totalTrades += 2
	if pnl > 0 {
		pm.winningTrades += 1
		if pnl > pm.largestWin {
			pm.largestWin = pnl
		}
	} else {
		pm.losingTrades += 1
		if pnl < pm.largestLoss {
			pm.largestLoss = pnl
		}
	}

	// Add trades to history
	pm.trades = append(pm.trades, buyTrade, sellTrade)

	// Log the execution
	log.Printf("🔄 EXECUTED ARBITRAGE: Buy %.4f %s on %s at $%.2f, Sell on %s at $%.2f", 
		quantity, opp.Symbol, opp.BuyExchange, opp.BuyPrice, opp.SellExchange, opp.SellPrice)
	log.Printf("💰 P&L: $%.2f (%.2f%%)", pnl, (pnl/pm.tradeSize)*100)

	return nil
}

// GetCurrentPnL returns the current profit/loss status
func (pm *PnLManager) GetCurrentPnL() PnLStatus {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	winRate := 0.0
	if pm.totalTrades > 0 {
		winRate = float64(pm.winningTrades) / float64(pm.totalTrades) * 100
	}

	return PnLStatus{
		CurrentBalance:    pm.balance,
		InitialBalance:    pm.initialBalance,
		TotalPnL:          pm.totalPnL,
		TotalPnLPercent:   ((pm.balance - pm.initialBalance) / pm.initialBalance) * 100,
		TotalTrades:       pm.totalTrades,
		WinningTrades:     pm.winningTrades,
		LosingTrades:      pm.losingTrades,
		WinRate:           winRate,
		LargestWin:        pm.largestWin,
		LargestLoss:       pm.largestLoss,
		AveragePnL:        pm.getAveragePnL(),
		LastUpdate:        time.Now(),
	}
}

// getAveragePnL calculates the average P&L per trade
func (pm *PnLManager) getAveragePnL() float64 {
	if pm.totalTrades == 0 {
		return 0
	}
	return pm.totalPnL / float64(pm.totalTrades)
}

// GetTradeHistory returns the recent trade history
func (pm *PnLManager) GetTradeHistory(limit int) []Trade {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	if limit <= 0 || limit > len(pm.trades) {
		limit = len(pm.trades)
	}

	start := len(pm.trades) - limit
	if start < 0 {
		start = 0
	}

	result := make([]Trade, limit)
	copy(result, pm.trades[start:])
	return result
}

// PnLStatus represents the current P&L status
type PnLStatus struct {
	CurrentBalance    float64
	InitialBalance    float64
	TotalPnL          float64
	TotalPnLPercent   float64
	TotalTrades       int
	WinningTrades     int
	LosingTrades      int
	WinRate           float64
	LargestWin        float64
	LargestLoss       float64
	AveragePnL        float64
	LastUpdate        time.Time
}

// PrintPnLStatus prints the current P&L status in a formatted way
func (pm *PnLManager) PrintPnLStatus() {
	status := pm.GetCurrentPnL()
	
	log.Println("=== PROFIT/LOSS STATUS ===")
	log.Printf("💰 Current Balance: $%.2f", status.CurrentBalance)
	log.Printf("📈 Total P&L: $%.2f (%.2f%%)", status.TotalPnL, status.TotalPnLPercent)
	log.Printf("📊 Total Trades: %d", status.TotalTrades)
	log.Printf("✅ Winning Trades: %d", status.WinningTrades)
	log.Printf("❌ Losing Trades: %d", status.LosingTrades)
	log.Printf("🎯 Win Rate: %.1f%%", status.WinRate)
	log.Printf("📈 Largest Win: $%.2f", status.LargestWin)
	log.Printf("📉 Largest Loss: $%.2f", status.LargestLoss)
	log.Printf("📊 Average P&L per Trade: $%.2f", status.AveragePnL)
	log.Printf("🕐 Last Update: %s", status.LastUpdate.Format("15:04:05"))
	log.Println("==========================")
}

// GetPnLSummary returns a concise P&L summary
func (pm *PnLManager) GetPnLSummary() string {
	status := pm.GetCurrentPnL()
	return fmt.Sprintf("P&L: $%.2f (%.2f%%) | Trades: %d | Win Rate: %.1f%%", 
		status.TotalPnL, status.TotalPnLPercent, status.TotalTrades, status.WinRate)
} 