// Package main demonstrates before and after model callbacks in ADK.
// This example shows how to use model callbacks to:
// - Filter inappropriate content before it reaches the model (before_model_callback)
// - Replace negative words with positive alternatives in responses (after_model_callback)
// - Log model interactions and measure processing time
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
)

// beforeModelCallback runs before the model processes a request
// It filters inappropriate content and logs request info
func beforeModelCallback(ctx agent.CallbackContext, llmRequest *model.LLMRequest) (*model.LLMResponse, error) {
	// Get the state and agent name
	state := ctx.State()
	agentName := ctx.AgentName()

	// Extract the last user message
	var lastUserMessage string
	if len(llmRequest.Contents) > 0 {
		for i := len(llmRequest.Contents) - 1; i >= 0; i-- {
			content := llmRequest.Contents[i]
			if content.Role == "user" && len(content.Parts) > 0 {
				lastUserMessage = content.Parts[0].Text
				break
			}
		}
	}

	// Log the request
	fmt.Println("=== MODEL REQUEST STARTED ===")
	fmt.Printf("Agent: %s\n", agentName)
	if lastUserMessage != "" {
		truncated := lastUserMessage
		if len(truncated) > 100 {
			truncated = truncated[:100] + "..."
		}
		fmt.Printf("User message: %s\n", truncated)
		// Store for later use
		if err := state.Set("last_user_message", lastUserMessage); err != nil {
			return nil, fmt.Errorf("failed to set last_user_message: %w", err)
		}
	} else {
		fmt.Println("User message: <empty>")
	}

	fmt.Printf("Timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	// Check for inappropriate content
	if lastUserMessage != "" && strings.Contains(strings.ToLower(lastUserMessage), "suck") {
		fmt.Println("=== INAPPROPRIATE CONTENT BLOCKED ===")
		fmt.Println("Blocked text containing prohibited word: 'suck'")
		fmt.Println("[BEFORE MODEL] ⚠️ Request blocked due to inappropriate content")

		// Return a response to skip the model call
		content := &genai.Content{
			Role: "model",
			Parts: []*genai.Part{
				{
					Text: "I cannot respond to messages containing inappropriate language. " +
						"Please rephrase your request without using words like 'sucks'.",
				},
			},
		}

		return &model.LLMResponse{
			Content: content,
		}, nil
	}

	// Record start time for duration calculation
	if err := state.Set("model_start_time", time.Now()); err != nil {
		return nil, fmt.Errorf("failed to set model_start_time: %w", err)
	}

	fmt.Println("[BEFORE MODEL] ✓ Request approved for processing")

	// Return nil to proceed with normal model request
	return nil, nil
}

// afterModelCallback runs after the model returns a response
// It modifies response text to replace negative words with positive alternatives
func afterModelCallback(ctx agent.CallbackContext, llmResponse *model.LLMResponse, llmResponseError error) (*model.LLMResponse, error) {
	fmt.Println("[AFTER MODEL] Processing response")

	// Skip processing if there's an error or response is empty
	if llmResponseError != nil {
		return nil, nil
	}

	// Skip processing if response is empty or has no text content
	if llmResponse == nil || llmResponse.Content == nil || llmResponse.Content.Parts == nil || len(llmResponse.Content.Parts) == 0 {
		return nil, nil
	}

	// Extract text from the response
	var responseText string
	for _, part := range llmResponse.Content.Parts {
		responseText += part.Text
	}

	if responseText == "" {
		return nil, nil
	}

	// Simple word replacements (case-insensitive)
	replacements := map[string]string{
		"problem":   "challenge",
		"difficult": "complex",
		"hard":      "challenging",
		"bad":       "suboptimal",
		"terrible":  "problematic",
		"awful":     "suboptimal",
		"hate":      "dislike",
	}

	// Perform replacements
	modifiedText := responseText
	modified := false

	for original, replacement := range replacements {
		if strings.Contains(strings.ToLower(modifiedText), original) {
			// Replace with proper case handling
			modifiedText = replaceCaseInsensitive(modifiedText, original, replacement)
			modified = true
		}
	}

	// Return modified response if changes were made
	if modified {
		fmt.Println("[AFTER MODEL] ↺ Modified response text")

		// Create a copy of the response with modified text
		modifiedResponse := &model.LLMResponse{
			Content: &genai.Content{
				Role:  llmResponse.Content.Role,
				Parts: []*genai.Part{},
			},
		}

		// Copy all parts, but modify the text parts
		for _, part := range llmResponse.Content.Parts {
			newPart := &genai.Part{}
			*newPart = *part // Copy all fields
			newPart.Text = modifiedText
			modifiedResponse.Content.Parts = append(modifiedResponse.Content.Parts, newPart)
		}

		// Copy other fields from original response
		modifiedResponse.CitationMetadata = llmResponse.CitationMetadata
		modifiedResponse.GroundingMetadata = llmResponse.GroundingMetadata
		modifiedResponse.UsageMetadata = llmResponse.UsageMetadata
		modifiedResponse.CustomMetadata = llmResponse.CustomMetadata
		modifiedResponse.LogprobsResult = llmResponse.LogprobsResult
		modifiedResponse.Partial = llmResponse.Partial
		modifiedResponse.TurnComplete = llmResponse.TurnComplete
		modifiedResponse.Interrupted = llmResponse.Interrupted
		modifiedResponse.ErrorCode = llmResponse.ErrorCode
		modifiedResponse.ErrorMessage = llmResponse.ErrorMessage
		modifiedResponse.FinishReason = llmResponse.FinishReason
		modifiedResponse.AvgLogprobs = llmResponse.AvgLogprobs

		return modifiedResponse, nil
	}

	// Return nil to use the original response
	return nil, nil
}

// replaceCaseInsensitive performs case-insensitive string replacement
func replaceCaseInsensitive(text, old, new string) string {
	lowerOld := strings.ToLower(old)

	var result string
	i := 0
	for i < len(text) {
		if i+len(old) <= len(text) && strings.ToLower(text[i:i+len(old)]) == lowerOld {
			// Preserve original case
			segment := text[i : i+len(old)]
			if segment == strings.ToUpper(segment) {
				result += strings.ToUpper(new)
			} else if segment == strings.ToLower(segment) {
				result += strings.ToLower(new)
			} else {
				// Title case or mixed case
				result += new
			}
			i += len(old)
		} else {
			result += string(text[i])
			i++
		}
	}
	return result
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

	// Create the agent with before and after model callbacks
	a, err := llmagent.New(llmagent.Config{
		Name:        "content_filter_agent",
		Model:       model,
		Description: "An agent that demonstrates model callbacks for content filtering and logging",
		Instruction: `You are a helpful assistant.

Your job is to:
- Answer user questions concisely
- Provide factual information
- Be friendly and respectful`,
		BeforeModelCallbacks: []llmagent.BeforeModelCallback{beforeModelCallback},
		AfterModelCallbacks:  []llmagent.AfterModelCallback{afterModelCallback},
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
