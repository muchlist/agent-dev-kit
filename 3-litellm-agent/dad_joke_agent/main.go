// Package main provides a dad joke agent example using ADK.
// Note: This Go version uses Gemini instead of LiteLLM/OpenRouter due to current limitations.
package main

import (
	"context"
	"log"
	"math/rand"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// getDadJokeArgs defines the input parameters for the dad joke tool (none in this case)
type getDadJokeArgs struct{}

// getDadJokeResults defines the return structure for the dad joke tool
type getDadJokeResults struct {
	Joke string `json:"joke"`
}

// getDadJoke returns a random dad joke from a predefined list
func getDadJoke(ctx tool.Context, input getDadJokeArgs) (getDadJokeResults, error) {
	jokes := []string{
		"Why did the chicken cross the road? To get to the other side!",
		"What do you call a belt made of watches? A waist of time.",
		"What do you call fake spaghetti? An impasta!",
		"Why did the scarecrow win an award? Because he was outstanding in his field!",
	}

	// Select a random joke
	randomJoke := jokes[rand.Intn(len(jokes))]

	return getDadJokeResults{Joke: randomJoke}, nil
}

func main() {
	godotenv.Load()
	ctx := context.Background()

	// IMPORTANT NOTE:
	// The Python version of this example uses LiteLLM to connect to OpenAI/OpenRouter models.
	// However, Go ADK currently does not have native LiteLLM integration like Python ADK does.
	// Therefore, this Go version uses Gemini instead.
	//
	// Go ADK Model Support Status:
	// ✓ Gemini (native support via google.golang.org/adk/model/gemini)
	// ✗ OpenAI (no native support yet)
	// ✗ Anthropic (no native support yet)
	// ✗ LiteLLM (no integration like Python ADK has)
	//
	// While some sources claim ADK-Go is "model-agnostic" and supports OpenAI/Anthropic,
	// concrete implementation examples and packages are not currently available in the
	// official Go ADK package (as of 2025).

	// Create the Gemini model with API key from environment
	model, err := gemini.NewModel(ctx, "gemini-2.0-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	// Create the dad joke tool
	dadJokeTool, err := functiontool.New(
		functiontool.Config{
			Name:        "get_dad_joke",
			Description: "Returns a random dad joke",
		},
		getDadJoke)
	if err != nil {
		log.Fatalf("Failed to create dad joke tool: %v", err)
	}

	// Create the dad joke agent
	a, err := llmagent.New(llmagent.Config{
		Name:        "dad_joke_agent",
		Model:       model,
		Description: "Dad joke agent",
		Instruction: `You are a helpful assistant that can tell dad jokes.
Only use the tool 'get_dad_joke' to tell jokes.`,
		Tools: []tool.Tool{dadJokeTool},
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
