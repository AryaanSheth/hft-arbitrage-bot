# HFT Arbitrage Bot

A high-frequency trading arbitrage bot that monitors multiple cryptocurrency exchanges simultaneously to identify and alert on arbitrage opportunities.

## Features

- **Multi-Exchange Monitoring**: Simultaneously monitors Binance, Kraken, and OKX
- **Real-time Price Feeds**: Uses WebSocket connections for low-latency price updates
- **Arbitrage Detection**: Automatically identifies profitable arbitrage opportunities
- **Configurable Thresholds**: Set minimum spread percentages for arbitrage alerts
- **Thread-Safe**: Concurrent processing with proper synchronization

## Architecture

```
hft-arbitrage-bot/
├── exchange/          # Exchange-specific WebSocket implementations
│   ├── binance.go     # Binance BTC/USDT price feed
│   ├── kraken.go      # Kraken XBT/USD price feed
│   └── okx.go         # OKX BTC-USDT price feed
├── strategy/          # Arbitrage strategy implementation
│   └── arbitrage.go   # Main arbitrage detection logic
└── main.go           # Application entry point and coordination
```

## How It Works

1. **Exchange Connections**: Each exchange runs in its own goroutine, maintaining WebSocket connections
2. **Quote Aggregation**: All price quotes are sent to a central channel
3. **Arbitrage Analysis**: The strategy continuously analyzes quotes from all exchanges
4. **Opportunity Detection**: When profitable spreads are found, alerts are generated

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd hft-arbitrage-bot
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the application:
```bash
go build -o hft-bot .
```

## Usage

Run the arbitrage bot:
```bash
./hft-bot
```

The bot will:
- Connect to all three exchanges (Binance, Kraken, OKX)
- Start monitoring BTC price feeds
- Display real-time price updates
- Alert when arbitrage opportunities are found

Press `Ctrl+C` to gracefully shutdown the bot.

## Configuration

### Minimum Spread Threshold

The bot is configured with a 0.1% minimum spread threshold by default. This means arbitrage opportunities will only be reported if the potential profit is at least 0.1% of the purchase price.

To modify this threshold, edit the `main.go` file:

```go
// Create arbitrage strategy with 0.1% minimum spread
arbitrageStrategy := strategy.NewArbitrageStrategy(0.1) // Change this value
```

### Supported Trading Pairs

Currently, the bot monitors:
- **Binance**: BTC/USDT
- **Kraken**: XBT/USD  
- **OKX**: BTC-USDT

## Output Example

```
🚀 Starting HFT Arbitrage Bot
✅ All exchanges started successfully
📊 Monitoring for arbitrage opportunities...
💡 Minimum spread threshold: 0.1%

Connected to Binance stream for btcusdt
Subscribed to Kraken BTC/USD order book
Subscribed to OKX BTC-USDT order book

Binance: Bid=43250.50, Ask=43251.00
Kraken: Bid=43248.75, Ask=43249.25
OKX: Bid=43249.00, Ask=43249.50

=== ARBITRAGE OPPORTUNITIES ===
💰 BUY on binance at 43251.00, SELL on kraken at 43248.75
   Spread: $2.25 (0.005%)
   Time: 14:30:25.123
---
```

## Technical Details

### Concurrency Model

- Each exchange runs in its own goroutine
- Quotes are sent through a buffered channel (capacity: 1000)
- The arbitrage strategy processes quotes every 100ms
- Thread-safe quote storage with read-write mutex

### Error Handling

- Graceful WebSocket reconnection handling
- Channel overflow protection
- Invalid data filtering
- Graceful shutdown on interrupt signals

### Performance Considerations

- Non-blocking quote transmission
- Efficient arbitrage calculation algorithm
- Minimal memory allocation in hot paths
- Configurable update frequency

## Disclaimer

This bot is for educational purposes only. Cryptocurrency trading involves significant risk. Always:
- Test thoroughly before using real funds
- Understand the risks involved
- Consider transaction fees and slippage
- Comply with local regulations
- Never invest more than you can afford to lose

## License

This project is licensed under the MIT License - see the LICENSE file for details. 