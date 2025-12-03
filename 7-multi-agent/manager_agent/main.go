// Package main demonstrates multi-agent systems in ADK with modular organization.
// This example creates a manager agent that delegates tasks to specialized agents:
// - Stock Analyst: Provides stock market information (agents/stock_analyst.go)
// - Funny Nerd: Tells nerdy jokes about technical topics (agents/funny_nerd.go)
// - News Analyst: Provides current technology news (agents/news_analyst.go)
//
// Each agent is organized in its own file for better maintainability.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"

	"github.com/muchlist/agent-dev-kit/7-multi-agent/manager_agent/agents"
	"github.com/muchlist/agent-dev-kit/7-multi-agent/manager_agent/tools"
)

const (
	MODEL_NAME = "gemini-2.0-flash"
)

// ===== Manager Agent Creation =====

// createManagerAgent creates the root manager agent that coordinates other agents
func createManagerAgent(_ context.Context, mdl model.LLM, stockAnalyst, funnyNerd, newsAnalyst agent.Agent) (agent.Agent, error) {
	// Create get_current_time tool from tools package
	getCurrentTimeTool, err := tools.NewGetCurrentTimeTool()
	if err != nil {
		return nil, fmt.Errorf("failed to create get_current_time tool: %w", err)
	}

	// Wrap news_analyst as an AgentTool
	// This allows the manager to use it like a tool while maintaining control
	// Note: In Go ADK, agents with built-in tools should be wrapped as AgentTools
	newsAnalystTool := agenttool.New(newsAnalyst, &agenttool.Config{})

	// Create manager agent with sub-agents and tools
	manager, err := llmagent.New(llmagent.Config{
		Name:        "manager",
		Model:       mdl,
		Description: "Manager agent that coordinates specialized agents",
		Instruction: `You are a manager agent that is responsible for overseeing the work of the other agents.

Always delegate the task to the appropriate agent. Use your best judgement
to determine which agent to delegate to.

You are responsible for delegating tasks to the following agents:
- stock_analyst: Use this agent for questions about stock prices, market data, or financial information
- funny_nerd: Use this agent when users want to hear nerdy jokes about technical topics

You also have access to the following tools:
- news_analyst: Use this tool to search and analyze current news (especially tech news)
- get_current_time: Use this tool to get the current date and time

When a user asks a question:
1. Determine if it's about stocks (→ delegate to stock_analyst)
2. Determine if it's about nerdy jokes (→ delegate to funny_nerd)
3. Determine if it's about news (→ use news_analyst tool)
4. Determine if it's about current time (→ use get_current_time tool)
5. For general questions, you can answer directly

Be friendly and helpful in your responses!`,
		SubAgents: []agent.Agent{stockAnalyst, funnyNerd},
		Tools:     []tool.Tool{newsAnalystTool, getCurrentTimeTool},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create manager agent: %w", err)
	}

	return manager, nil
}

// ===== Main Function =====

func main() {
	godotenv.Load()
	ctx := context.Background()

	// Create the Gemini model
	model, err := gemini.NewModel(ctx, MODEL_NAME, &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	fmt.Println("🤖 Creating specialized agents...")

	// Create specialized agents using modular agent constructors
	stockAnalyst, err := agents.NewStockAnalyst(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create stock analyst agent: %v", err)
	}
	fmt.Println("  ✓ Stock Analyst agent created")

	funnyNerd, err := agents.NewFunnyNerd(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create funny nerd agent: %v", err)
	}
	fmt.Println("  ✓ Funny Nerd agent created")

	newsAnalyst, err := agents.NewNewsAnalyst(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create news analyst agent: %v", err)
	}
	fmt.Println("  ✓ News Analyst agent created")

	// Create manager agent that coordinates all specialized agents
	fmt.Println("🎯 Creating manager agent...")
	managerAgent, err := createManagerAgent(ctx, model, stockAnalyst, funnyNerd, newsAnalyst)
	if err != nil {
		log.Fatalf("Failed to create manager agent: %v", err)
	}
	fmt.Println("  ✓ Manager agent created")

	fmt.Println("\n🚀 Launching Multi-Agent System...")
	fmt.Println("========================================")

	// Configure and launch the agent
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(managerAgent),
	}

	l := full.NewLauncher()
	if err := l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
