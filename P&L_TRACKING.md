# Profit/Loss Tracking Guide

This guide explains how to check your profit/loss status at any time when your HFT arbitrage bot is running.

## Overview

The bot now includes a comprehensive P&L tracking system that provides:

- Real-time profit/loss monitoring
- Trade execution tracking
- Performance statistics
- Multiple ways to check P&L status

## How to Check P&L

### 1. Interactive Console Commands

When the bot is running, you can use these commands in the console:

- **Press Enter** - Check current P&L status
- **Type `pnl`** - Check P&L status
- **Type `trades`** - Show recent trade history
- **Type `help`** - Show available commands

Example output:
```
=== PROFIT/LOSS STATUS ===
💰 Current Balance: $1050.25
📈 Total P&L: $50.25 (5.03%)
📊 Total Trades: 12
✅ Winning Trades: 8
❌ Losing Trades: 4
🎯 Win Rate: 66.7%
📈 Largest Win: $15.50
📉 Largest Loss: -$8.25
📊 Average P&L per Trade: $4.19
🕐 Last Update: 14:30:25
==========================
```

### 2. HTTP API Endpoints

The bot runs a web API server on port 8080. You can access these endpoints:

- **GET http://localhost:8080/pnl** - Full P&L status (JSON)
- **GET http://localhost:8080/summary** - P&L summary (JSON)
- **GET http://localhost:8080/trades** - Recent trades (JSON)
- **GET http://localhost:8080/health** - API health check

Example API response:
```json
{
  "status": "success",
  "data": {
    "current_balance": 1050.25,
    "total_pnl": 50.25,
    "total_pnl_percent": 5.03,
    "total_trades": 12,
    "win_rate": 66.7
  },
  "timestamp": 1703123425
}
```

### 3. Command-Line Client Tool

Build and use the P&L client tool:

```bash
# Build the client
make pnl-client

# Check P&L summary
./tools/pnl_client summary

# Get full P&L status
./tools/pnl_client pnl

# Get recent trades
./tools/pnl_client trades --limit 10

# Check API health
./tools/pnl_client health

# Use with different host
./tools/pnl_client summary --host 192.168.1.100:8080
```

### 4. Web Browser

Open your web browser and navigate to:
- http://localhost:8080/pnl
- http://localhost:8080/summary
- http://localhost:8080/trades

## Configuration

### Initial Settings

The bot starts with these default settings:
- **Initial Balance**: $1,000.00
- **Trade Size**: $100.00 per arbitrage
- **Minimum Spread**: 0.1%

You can modify these in `main.go`:

```go
// Create arbitrage strategy with custom settings
arbitrageStrategy := strategy.NewArbitrageStrategy(
    0.1,        // 0.1% minimum spread
    1000.0,     // $1000 initial balance
    100.0,      // $100 trade size
)
```

### API Port

The API server runs on port 8080 by default. You can change this in `main.go`:

```go
pnlAPI := api.NewPnLAPI(arbitrageStrategy.GetPnLManager(), 8080)
```

## P&L Metrics Explained

### Current Balance
Your current account balance after all trades.

### Total P&L
Total profit/loss in dollars and percentage since starting.

### Trade Statistics
- **Total Trades**: Number of arbitrage trades executed
- **Winning Trades**: Trades that resulted in profit
- **Losing Trades**: Trades that resulted in loss
- **Win Rate**: Percentage of profitable trades

### Performance Metrics
- **Largest Win**: Highest single trade profit
- **Largest Loss**: Highest single trade loss
- **Average P&L**: Average profit/loss per trade

## Monitoring Best Practices

### 1. Regular Checks
- Check P&L every few minutes during active trading
- Monitor win rate trends
- Watch for unusual losses

### 2. Performance Analysis
- Track win rate over time
- Monitor average P&L per trade
- Identify best performing arbitrage opportunities

### 3. Risk Management
- Set stop-loss limits
- Monitor largest loss values
- Track balance drawdown

## Troubleshooting

### API Not Responding
1. Check if the bot is running
2. Verify port 8080 is not blocked
3. Check firewall settings

### No Trades Showing
1. Verify exchanges are connected
2. Check minimum spread threshold
3. Ensure sufficient balance for trades

### Inaccurate P&L
1. Check trade execution logs
2. Verify exchange connections
3. Review arbitrage opportunity detection

## Building and Running

```bash
# Build everything
make all

# Run the bot
make run

# Test P&L client (requires bot to be running)
make test-pnl-client

# Clean build artifacts
make clean
```

## Security Notes

- The API server is for local monitoring only
- Don't expose port 8080 to the internet
- Consider adding authentication for production use
- Monitor API access logs

## Advanced Usage

### Custom Monitoring Scripts
You can create custom scripts to monitor P&L:

```bash
#!/bin/bash
# Monitor P&L every 30 seconds
while true; do
    ./tools/pnl_client summary
    sleep 30
done
```

### Integration with External Tools
The JSON API can be integrated with:
- Grafana dashboards
- TradingView alerts
- Slack notifications
- Email alerts

Example curl command:
```bash
curl -s http://localhost:8080/summary | jq '.data.total_pnl_percent'
``` 