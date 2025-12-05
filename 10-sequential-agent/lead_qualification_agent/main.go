// Package main implements a lead qualification sequential agent in Go.
// This example demonstrates how to create a SequentialAgent using Google's ADK framework.
//
// The lead qualification pipeline orchestrates three sub-agents in sequence:
// 1. Lead Validator Agent: Validates lead information completeness
// 2. Lead Scorer Agent: Scores the lead from 1-10 based on qualification criteria
// 3. Action Recommender Agent: Recommends next actions based on validation and scoring
//
// Each agent stores its output in session state using output keys, allowing the next
// agent in the sequence to access the results of previous agents.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"

	"github.com/muchlist/agent-dev-kit/10-sequential-agent/lead_qualification_agent/agents"
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

	// Create sub-agents for the sequential workflow
	validator, err := agents.NewLeadValidator(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create lead validator agent: %v", err)
	}

	scorer, err := agents.NewLeadScorer(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create lead scorer agent: %v", err)
	}

	recommender, err := agents.NewActionRecommender(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create action recommender agent: %v", err)
	}

	// Create the sequential agent using ADK SequentialAgent
	fmt.Println("🔗 Creating Sequential Agent...")
	sequentialAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "LeadQualificationPipeline",
			Description: "A sequential pipeline that validates, scores, and recommends actions for sales leads",
			SubAgents:   []agent.Agent{validator, scorer, recommender},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create lead qualification sequential agent: %v", err)
	}

	fmt.Println("\n🚀 Launching Lead Qualification Sequential Agent...")
	fmt.Println("========================================================")
	fmt.Println("Example qualified lead:")
	fmt.Println("Name: Sarah Johnson")
	fmt.Println("Email: sarah.j@techinnovate.com")
	fmt.Println("Phone: 555-123-4567")
	fmt.Println("Company: Tech Innovate Solutions")
	fmt.Println("Position: CTO")
	fmt.Println("Interest: Looking for an AI solution to automate customer support")
	fmt.Println("Budget: $50K-100K available")
	fmt.Println("Timeline: Next quarter")
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
