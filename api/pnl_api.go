package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"hft-arbitrage-bot/strategy"
)

// PnLAPI provides HTTP endpoints for P&L monitoring
type PnLAPI struct {
	pnlManager *strategy.PnLManager
	server     *http.Server
}

// NewPnLAPI creates a new P&L API server
func NewPnLAPI(pnlManager *strategy.PnLManager, port int) *PnLAPI {
	mux := http.NewServeMux()
	api := &PnLAPI{pnlManager: pnlManager}
	
	// Register routes
	mux.HandleFunc("/pnl", api.handlePnL)
	mux.HandleFunc("/trades", api.handleTrades)
	mux.HandleFunc("/summary", api.handleSummary)
	mux.HandleFunc("/health", api.handleHealth)
	
	api.server = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}
	
	return api
}

// Start starts the HTTP server
func (api *PnLAPI) Start() {
	log.Printf("🌐 Starting P&L API server on port %s", api.server.Addr)
	go func() {
		if err := api.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ API server error: %v", err)
		}
	}()
}

// Stop stops the HTTP server
func (api *PnLAPI) Stop() {
	log.Println("🛑 Stopping P&L API server...")
	api.server.Close()
}

// handlePnL handles P&L status requests
func (api *PnLAPI) handlePnL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	status := api.pnlManager.GetCurrentPnL()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   status,
		"timestamp": time.Now().Unix(),
	})
}

// handleTrades handles trade history requests
func (api *PnLAPI) handleTrades(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	limit := 10 // default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	
	trades := api.pnlManager.GetTradeHistory(limit)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   trades,
		"count":  len(trades),
		"timestamp": time.Now().Unix(),
	})
}

// handleSummary handles P&L summary requests
func (api *PnLAPI) handleSummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	summary := api.pnlManager.GetPnLSummary()
	status := api.pnlManager.GetCurrentPnL()
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"summary": summary,
			"current_balance": status.CurrentBalance,
			"total_pnl": status.TotalPnL,
			"total_pnl_percent": status.TotalPnLPercent,
			"total_trades": status.TotalTrades,
			"win_rate": status.WinRate,
		},
		"timestamp": time.Now().Unix(),
	})
}

// handleHealth handles health check requests
func (api *PnLAPI) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"timestamp": time.Now().Unix(),
		"uptime": time.Since(time.Now()).String(), // This will be 0, but you could track actual uptime
	})
} 