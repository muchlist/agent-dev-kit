package agents

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// ===== Funny Nerd Tool Structures =====

type getNerdJokeArgs struct {
	Topic string `json:"topic"`
}

type getNerdJokeResults struct {
	Status string `json:"status"`
	Joke   string `json:"joke"`
	Topic  string `json:"topic"`
}

// ===== Tool Implementation =====

// getNerdJoke returns a nerdy joke about a specific topic
func getNerdJoke(ctx tool.Context, input getNerdJokeArgs) (getNerdJokeResults, error) {
	fmt.Printf("--- Tool: get_nerd_joke called for topic: %s ---\n", input.Topic)

	// Collection of nerdy jokes by topic
	// In production, you might want to use a jokes API or larger database
	jokes := map[string]string{
		"python":      "Why don't Python programmers like to use inheritance? Because they don't like to inherit anything!",
		"javascript":  "Why did the JavaScript developer go broke? Because he used up all his cache!",
		"java":        "Why do Java developers wear glasses? Because they can't C#!",
		"go":          "Why do Go programmers prefer channels over callbacks? Because they don't want to get caught in callback hell!",
		"golang":      "What's a gopher's favorite type of code? Go code that's concurrent and simple!",
		"programming": "Why do programmers prefer dark mode? Because light attracts bugs!",
		"math":        "Why was the equal sign so humble? Because he knew he wasn't less than or greater than anyone else!",
		"physics":     "Why did the photon check into a hotel? Because it was travelling light!",
		"chemistry":   "Why did the acid go to the gym? To become a buffer solution!",
		"biology":     "Why did the cell go to therapy? Because it had too many issues!",
		"computer":    "Why did the computer keep freezing? It left its Windows open!",
		"database":    "Why did the DBA break up with their partner? Too many relationship conflicts!",
		"default":     "Why did the computer go to the doctor? Because it had a virus!",
	}

	// Find joke, use default if topic not found
	joke, exists := jokes[input.Topic]
	if !exists {
		joke = jokes["default"]
	}

	// Store last joke topic in session state
	state := ctx.State()
	state.Set("last_joke_topic", input.Topic)

	return getNerdJokeResults{
		Status: "success",
		Joke:   joke,
		Topic:  input.Topic,
	}, nil
}

// ===== Agent Creation =====

// NewFunnyNerd creates a specialized agent for telling nerdy jokes
func NewFunnyNerd(ctx context.Context, mdl model.LLM) (agent.Agent, error) {
	// Create get_nerd_joke tool
	getNerdJokeTool, err := functiontool.New(
		functiontool.Config{
			Name:        "get_nerd_joke",
			Description: "Get a nerdy joke about a specific topic",
		},
		getNerdJoke)
	if err != nil {
		return nil, fmt.Errorf("failed to create get_nerd_joke tool: %w", err)
	}

	// Create funny nerd agent
	funnyNerd, err := llmagent.New(llmagent.Config{
		Name:        "funny_nerd",
		Model:       mdl,
		Description: "An agent that tells nerdy jokes about various topics.",
		Instruction: `You are a funny nerd agent that tells nerdy jokes about various topics.

When asked to tell a joke:
1. Use the get_nerd_joke tool to fetch a joke about the requested topic
2. If no specific topic is mentioned, ask the user what kind of nerdy joke they'd like to hear
3. Format the response to include both the joke and a brief explanation if needed

Available topics include:
- python
- javascript
- java
- go / golang
- programming
- math
- physics
- chemistry
- biology
- computer
- database

Example response format:
"Here's a nerdy joke about <TOPIC>:

<JOKE>

😄 Explanation: {brief explanation if needed}"

If the user asks about anything else, you should delegate the task to the manager agent.`,
		Tools: []tool.Tool{getNerdJokeTool},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create funny nerd agent: %w", err)
	}

	return funnyNerd, nil
}
