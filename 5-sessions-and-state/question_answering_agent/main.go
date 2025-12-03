// Package main demonstrates sessions and state management in ADK.
// This example shows how to create sessions with initial state and use
// template variables to access that state in agent instructions.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
)

const (
	APP_NAME   = "Bot"
	USER_ID    = "muchlis"
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

	// Create the question answering agent with template variables
	// The {user_name} and {user_preferences} will be replaced with values from session state
	questionAnsweringAgent, err := llmagent.New(llmagent.Config{
		Name:        "question_answering_agent",
		Model:       model,
		Description: "Question answering agent",
		Instruction: `You are a helpful assistant that answers questions about the user's preferences.

Here is some information about the user:
Name:
{user_name}
Preferences:
{user_preferences}`,
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Create an in-memory session service
	sessionService := session.InMemoryService()

	// Define initial state with user information
	initialState := map[string]any{
		"user_name": "Muchlis",
		"user_preferences": `
        I like to play Pickleball, Disc Golf, and Tennis.
        My favorite food is Mexican.
        My favorite TV show is Game of Thrones.
        Loves it when people like and subscribe to his YouTube channel.
    `,
	}

	// Create a new session with the initial state
	SESSION_ID := uuid.New().String()
	_, err = sessionService.Create(ctx, &session.CreateRequest{
		AppName:   APP_NAME,
		UserID:    USER_ID,
		SessionID: SESSION_ID,
		State:     initialState,
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	fmt.Println("CREATED NEW SESSION:")
	fmt.Printf("\tSession ID: %s\n", SESSION_ID)
	fmt.Printf("\tApp Name: %s\n", APP_NAME)
	fmt.Printf("\tUser ID: %s\n", USER_ID)
	fmt.Println()

	// Create a runner with the agent and session service
	r, err := runner.New(runner.Config{
		AppName:        APP_NAME,
		Agent:          questionAnsweringAgent,
		SessionService: sessionService,
	})
	if err != nil {
		log.Fatalf("Failed to create runner: %v", err)
	}

	// Create a user message asking about stored preferences
	userMessage := &genai.Content{
		Role: "user",
		Parts: []*genai.Part{
			{Text: "What is Muchlis's favorite TV show?"},
		},
	}

	fmt.Println("User Question: What is Muchlis's favorite TV show?")
	fmt.Println()

	// Run the agent with the session context
	// The agent will have access to session state via template variables
	var finalResponse string
	for event, err := range r.Run(ctx, USER_ID, SESSION_ID, userMessage, agent.RunConfig{}) {
		if err != nil {
			log.Fatalf("Error during agent run: %v", err)
		}

		// Check if this is the final response
		if event.Content != nil && len(event.Content.Parts) > 0 {
			finalResponse = event.Content.Parts[0].Text
		}
	}

	fmt.Println("Final Response:", finalResponse)
	fmt.Println()

	// Retrieve and display the final session state
	fmt.Println("==== Session Event Exploration ====")
	getResp, err := sessionService.Get(ctx, &session.GetRequest{
		AppName:   APP_NAME,
		UserID:    USER_ID,
		SessionID: SESSION_ID,
	})
	if err != nil {
		log.Fatalf("Failed to get session: %v", err)
	}

	retrievedSession := getResp.Session
	fmt.Println("=== Final Session State ===")
	for key, value := range retrievedSession.State().All() {
		fmt.Printf("%s: %v\n", key, value)
	}

	// Display session history
	fmt.Println("\n=== Session Message History ===")
	events := retrievedSession.Events()
	count := 0
	for event := range events.All() {
		count++
		if event.Content != nil {
			role := event.Content.Role
			if len(event.Content.Parts) > 0 {
				text := event.Content.Parts[0].Text
				preview := text
				if len(text) > 100 {
					preview = text[:100] + "..."
				}
				fmt.Printf("[%d] %s: %s\n", count, role, preview)
			}
		}
	}

	fmt.Println("\nExample completed successfully!")
}
