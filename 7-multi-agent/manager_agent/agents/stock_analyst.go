package agents

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// ===== Stock Analyst Tool Structures =====

type getStockPriceArgs struct {
	Ticker string `json:"ticker"`
}

type getStockPriceResults struct {
	Status       string `json:"status"`
	Ticker       string `json:"ticker,omitempty"`
	Price        string `json:"price,omitempty"`
	Timestamp    string `json:"timestamp,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// ===== Tool Implementation =====

// getStockPrice retrieves current stock price using mock data
// Note: In production, replace with real stock API like Alpha Vantage or IEX Cloud
func getStockPrice(ctx tool.Context, input getStockPriceArgs) (getStockPriceResults, error) {
	fmt.Printf("--- Tool: get_stock_price called for %s ---\n", input.Ticker)

	// Mock stock prices for demonstration
	// In production, you would use a real stock API:
	// - Alpha Vantage: https://www.alphavantage.co/ (free tier: 5 API requests per minute)
	// - IEX Cloud: https://iexcloud.io/ (free tier available)
	// - Finnhub: https://finnhub.io/ (free tier available)
	mockPrices := map[string]string{
		"GOOG":  "175.34",
		"GOOGL": "175.50",
		"TSLA":  "156.78",
		"META":  "123.45",
		"AAPL":  "189.50",
		"MSFT":  "378.25",
		"AMZN":  "145.67",
	}

	price, exists := mockPrices[input.Ticker]
	if !exists {
		return getStockPriceResults{
			Status:       "error",
			ErrorMessage: fmt.Sprintf("Could not fetch price for %s. Available tickers: GOOG, GOOGL, TSLA, META, AAPL, MSFT, AMZN", input.Ticker),
		}, nil
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	return getStockPriceResults{
		Status:    "success",
		Ticker:    input.Ticker,
		Price:     price,
		Timestamp: currentTime,
	}, nil
}

// ===== Agent Creation =====

// NewStockAnalyst creates a specialized agent for stock market analysis
func NewStockAnalyst(ctx context.Context, mdl model.LLM) (agent.Agent, error) {
	// Create get_stock_price tool
	getStockPriceTool, err := functiontool.New(
		functiontool.Config{
			Name:        "get_stock_price",
			Description: "Retrieves current stock price for a given ticker symbol",
		},
		getStockPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to create get_stock_price tool: %w", err)
	}

	// Create stock analyst agent
	stockAnalyst, err := llmagent.New(llmagent.Config{
		Name:        "stock_analyst",
		Model:       mdl,
		Description: "An agent that can look up stock prices and track them over time.",
		Instruction: `You are a helpful stock market assistant that helps users track their stocks of interest.

When asked about stock prices:
1. Use the get_stock_price tool to fetch the latest price for the requested stock(s)
2. Format the response to show each stock's current price and the time it was fetched
3. If a stock price couldn't be fetched, mention this in your response

Example response format:
"Here are the current prices for your stocks:
- GOOG: $175.34 (updated at 2024-04-21 16:30:00)
- TSLA: $156.78 (updated at 2024-04-21 16:30:00)
- META: $123.45 (updated at 2024-04-21 16:30:00)"

Available tickers: GOOG, GOOGL, TSLA, META, AAPL, MSFT, AMZN`,
		Tools: []tool.Tool{getStockPriceTool},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create stock analyst agent: %w", err)
	}

	return stockAnalyst, nil
}
