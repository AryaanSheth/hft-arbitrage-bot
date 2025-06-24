.PHONY: build run clean pnl-client test

# Build the main bot
build:
	go build -o hft-bot main.go

# Build the P&L client tool
pnl-client:
	go build -o tools/pnl_client tools/pnl_client.go

# Run the bot
run: build
	./hft-bot

# Clean build artifacts
clean:
	rm -f hft-bot
	rm -f tools/pnl_client

# Test the P&L client (requires bot to be running)
test-pnl-client: pnl-client
	@echo "Testing P&L client..."
	@echo "Getting P&L summary:"
	./tools/pnl_client summary
	@echo ""
	@echo "Getting recent trades:"
	./tools/pnl_client trades --limit 5
	@echo ""
	@echo "Checking API health:"
	./tools/pnl_client health

# Install dependencies
deps:
	go mod download

# Run tests
test:
	go test ./...

# Build all
all: deps build pnl-client

# Help
help:
	@echo "Available targets:"
	@echo "  build           - Build the main bot"
	@echo "  pnl-client      - Build the P&L client tool"
	@echo "  run             - Build and run the bot"
	@echo "  clean           - Clean build artifacts"
	@echo "  test-pnl-client - Test the P&L client (requires bot to be running)"
	@echo "  deps            - Install dependencies"
	@echo "  test            - Run tests"
	@echo "  all             - Build everything"
	@echo "  help            - Show this help" 