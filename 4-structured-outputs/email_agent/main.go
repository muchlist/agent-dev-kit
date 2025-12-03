// Package main provides an email generator agent example using ADK with structured outputs.
package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
)

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

	// Define the output schema for structured email content
	// This ensures the LLM response is in a specific JSON format
	emailSchema := &genai.Schema{
		Type: "OBJECT",
		Properties: map[string]*genai.Schema{
			"subject": {
				Type:        "STRING",
				Description: "The subject line of the email. Should be concise and descriptive.",
			},
			"body": {
				Type:        "STRING",
				Description: "The main content of the email. Should be well-formatted with proper greeting, paragraphs, and signature.",
			},
		},
		Required: []string{"subject", "body"},
	}

	// Create the email generator agent with structured output
	a, err := llmagent.New(llmagent.Config{
		Name:        "email_agent",
		Model:       model,
		Description: "Generates professional emails with structured subject and body",
		Instruction: `You are an Email Generation Assistant.
Your task is to generate a professional email based on the user's request.

GUIDELINES:
- Create an appropriate subject line (concise and relevant)
- Write a well-structured email body with:
    * Professional greeting
    * Clear and concise main content
    * Appropriate closing
    * Your name as signature
- Email tone should match the purpose (formal for business, friendly for colleagues)
- Keep emails concise but complete

IMPORTANT: Your response MUST be valid JSON matching this structure:
{
    "subject": "Subject line here",
    "body": "Email body here with proper paragraphs and formatting"
}

DO NOT include any explanations or additional text outside the JSON response.`,
		OutputSchema: emailSchema,
		OutputKey:    "email",
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
