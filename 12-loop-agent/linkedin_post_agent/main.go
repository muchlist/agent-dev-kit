// Package main implements a LinkedIn post generator with iterative refinement using Loop Agent in Go.
// This example demonstrates how to create a hybrid workflow that combines Sequential and Loop agents
// for iterative content improvement until quality requirements are met.
//
// The LinkedIn Post Generator workflow:
// 1. Initial Post Generation: Creates first draft of LinkedIn post
// 2. Refinement Loop: Iteratively reviews and refines until quality criteria met
//
// Key patterns demonstrated:
// - Sequential pipeline with initial generation followed by iterative refinement
// - Loop agent with max iterations and exit conditions
// - Quality-driven loop termination using exit tools
// - Feedback-based improvement process
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/loopagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"

	"github.com/muchlist/agent-dev-kit/12-loop-agent/linkedin_post_agent/agents"
)

const (
	MODEL_NAME = "gemini-2.0-flash"
)

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

	fmt.Println("📝 Creating LinkedIn Post Generator with Iterative Refinement...")

	// Create sub-agents for the refinement loop
	postReviewer, err := agents.NewPostReviewer(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create post reviewer agent: %v", err)
	}
	fmt.Println("  ✓ Post Reviewer agent created")

	postRefiner, err := agents.NewPostRefiner(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create post refiner agent: %v", err)
	}
	fmt.Println("  ✓ Post Refiner agent created")

	// Create initial post generator
	initialPostGenerator, err := agents.NewInitialPostGenerator(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create initial post generator agent: %v", err)
	}
	fmt.Println("  ✓ Initial Post Generator agent created")

	// Create Loop Agent for iterative refinement
	fmt.Println("🔄 Creating Refinement Loop...")
	refinementLoop, err := loopagent.New(loopagent.Config{
		MaxIterations: 8,
		AgentConfig: agent.Config{
			Name:        "PostRefinementLoop",
			Description: "Iteratively reviews and refines LinkedIn post until quality requirements are met",
			SubAgents:   []agent.Agent{postReviewer, postRefiner},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create refinement loop agent: %v", err)
	}
	fmt.Println("  ✓ Refinement Loop created (max 10 iterations)")

	// Create Sequential Agent for overall pipeline
	fmt.Println("🔗 Creating Sequential Pipeline...")
	sequentialAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "LinkedInPostGenerationPipeline",
			Description: "Generates and refines LinkedIn post through iterative review process",
			SubAgents:   []agent.Agent{initialPostGenerator, refinementLoop},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create LinkedIn post generation pipeline: %v", err)
	}
	fmt.Println("  ✓ LinkedIn Post Generation Pipeline created")

	fmt.Println("\n🚀 Launching LinkedIn Post Generator with Loop Agent...")
	fmt.Println("========================================================")
	fmt.Println("Example prompt to try:")
	fmt.Println("Generate a LinkedIn post about what I've learned from Agent Development Kit tutorial.")
	fmt.Println("========================================================")

	// Configure and launch the agent
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(sequentialAgent),
	}

	l := full.NewLauncher()
	if err := l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
