// Package main demonstrates before and after agent callbacks in ADK.
// This example shows how to use callbacks to:
// - Log when agent processing starts and ends
// - Track request counts across sessions
// - Measure request duration
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
)

// beforeAgentCallback runs when the agent starts processing a request
func beforeAgentCallback(ctx agent.CallbackContext) (*genai.Content, error) {
	// Get the session state
	state := ctx.State()

	// Set agent name if not present
	if _, err := state.Get("agent_name"); err != nil {
		if err := state.Set("agent_name", "SimpleChatBot"); err != nil {
			return nil, fmt.Errorf("failed to set agent_name: %w", err)
		}
	}

	// Initialize request counter
	var counter int64 = 1
	if val, err := state.Get("request_counter"); err == nil {
		if counterVal, ok := val.(int64); ok {
			counter = counterVal + 1
		}
	}
	if err := state.Set("request_counter", counter); err != nil {
		return nil, fmt.Errorf("failed to set request_counter: %w", err)
	}

	// Store start time for duration calculation
	startTime := time.Now()
	if err := state.Set("request_start_time", startTime); err != nil {
		return nil, fmt.Errorf("failed to set request_start_time: %w", err)
	}

	// Log the request
	fmt.Println("=== AGENT EXECUTION STARTED ===")
	fmt.Printf("Request #: %d\n", counter)
	fmt.Printf("Timestamp: %s\n", startTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("\n[BEFORE CALLBACK] Agent processing request #%d\n", counter)

	// Return nil to continue with normal agent processing
	return nil, nil
}

// afterAgentCallback runs when the agent finishes processing a request
func afterAgentCallback(ctx agent.CallbackContext) (*genai.Content, error) {
	// Get the session state
	state := ctx.State()

	// Calculate request duration if start time is available
	var duration float64
	timestamp := time.Now()
	if val, err := state.Get("request_start_time"); err == nil {
		if startTime, ok := val.(time.Time); ok {
			duration = timestamp.Sub(startTime).Seconds()
		}
	}

	// Get request counter
	var counter int64
	if val, err := state.Get("request_counter"); err == nil {
		if counterVal, ok := val.(int64); ok {
			counter = counterVal
		}
	}

	// Log the completion
	fmt.Println("=== AGENT EXECUTION COMPLETED ===")
	fmt.Printf("Request #: %d\n", counter)
	if duration > 0 {
		fmt.Printf("Duration: %.2f seconds\n", duration)
	}

	fmt.Printf("[AFTER CALLBACK] Agent completed request #%d\n", counter)
	if duration > 0 {
		fmt.Printf("[AFTER CALLBACK] Processing took %.2f seconds\n", duration)
	}

	// Return nil to continue with normal agent processing
	return nil, nil
}

func main() {
	godotenv.Load()
	ctx := context.Background()

	// Create the Gemini model with API key from environment
	model, err := gemini.NewModel(ctx, "gemini-2.0-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	// Create the agent with before and after callbacks
	a, err := llmagent.New(llmagent.Config{
		Name:        "before_after_agent",
		Model:       model,
		Description: "A basic agent that demonstrates before and after agent callbacks",
		Instruction: fmt.Sprintf(`You are a friendly greeting agent. Your name is {agent_name}.

Your job is to:
- Greet users politely
- Respond to basic questions
- Keep your responses friendly and concise

Current request counter: %d`, 1),
		BeforeAgentCallbacks: []agent.BeforeAgentCallback{beforeAgentCallback},
		AfterAgentCallbacks:  []agent.AfterAgentCallback{afterAgentCallback},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Configure and launch the agent
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(a),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
