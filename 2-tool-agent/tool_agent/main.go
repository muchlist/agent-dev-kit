// Package main provides a tool agent example using ADK with Google Search.
package main

import (
	"context"
	"log"
	"os"

	// "time"

	"github.com/joho/godotenv"
	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	// "google.golang.org/adk/tool/functiontool"
)

// Custom function tool example (commented out)
// Uncomment to use this instead of Google Search

// type getCurrentTimeArgs struct{}
//
// type getCurrentTimeResults struct {
// 	CurrentTime string `json:"current_time"`
// }
//
// func getCurrentTime(ctx tool.Context, input getCurrentTimeArgs) (getCurrentTimeResults, error) {
// 	currentTime := time.Now().Format("2006-01-02 15:04:05")
// 	return getCurrentTimeResults{CurrentTime: currentTime}, nil
// }

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

	// Option 1: Using built-in Google Search tool (default)
	tools := []tool.Tool{
		geminitool.GoogleSearch{},
	}

	// Option 2: Using custom function tool (commented out)
	// Uncomment the lines below and comment out the Google Search tool above to use the custom tool
	//
	// currentTimeTool, err := functiontool.New(
	// 	functiontool.Config{
	// 		Name:        "get_current_time",
	// 		Description: "Get the current time in the format YYYY-MM-DD HH:MM:SS",
	// 	},
	// 	getCurrentTime)
	// if err != nil {
	// 	log.Fatalf("Failed to create current time tool: %v", err)
	// }
	// tools = []tool.Tool{currentTimeTool}

	// IMPORTANT NOTE:
	// Currently, for each root agent or single agent, only ONE built-in tool is supported.
	// You CANNOT mix built-in tools (like GoogleSearch) with custom function tools in the same agent.
	// To use both types, you would need to use a multi-agent approach.
	//
	// This WILL NOT WORK:
	// tools = []tool.Tool{
	//     geminitool.GoogleSearch{},
	//     currentTimeTool,
	// }

	// Create the tool agent
	a, err := llmagent.New(llmagent.Config{
		Name:        "tool_agent",
		Model:       model,
		Description: "Tool agent",
		Instruction: `You are a helpful assistant that can use the following tools:
- google_search`,
		Tools: tools,
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
