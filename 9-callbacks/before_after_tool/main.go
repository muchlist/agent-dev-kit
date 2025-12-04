// Package main demonstrates before and after tool callbacks in ADK.
// This example shows how to use tool callbacks to:
// - Modify tool arguments before execution (before_tool_callback)
// - Block certain tool calls completely (before_tool_callback)
// - Enhance tool responses with additional information (after_tool_callback)
// - Handle errors gracefully (after_tool_callback)
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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

// ===== Tool Structures =====

type getCapitalCityArgs struct {
	Country string `json:"country"`
}

type getCapitalCityResults struct {
	Result string `json:"result"`
}

// getCapitalCity retrieves the capital city of a given country
func getCapitalCity(ctx tool.Context, input getCapitalCityArgs) (getCapitalCityResults, error) {
	country := input.Country
	fmt.Printf("[TOOL] Executing get_capital_city tool with country: '%s'\n", country)
	fmt.Printf("[TOOL] Input struct: %+v\n", input)

	countryCapitals := map[string]string{
		"united states": "Washington, D.C.",
		"usa":           "Washington, D.C.",
		"canada":        "Ottawa",
		"france":        "Paris",
		"germany":       "Berlin",
		"japan":         "Tokyo",
		"brazil":        "Brasília",
		"australia":     "Canberra",
		"india":         "New Delhi",
	}

	// Use lowercase for comparison
	result, exists := countryCapitals[strings.ToLower(country)]
	if !exists {
		result = fmt.Sprintf("Capital not found for %s", country)
	}

	fmt.Printf("[TOOL] Result: %s\n", result)
	fmt.Printf("[TOOL] Returning: {Result: '%s'}\n", result)

	return getCapitalCityResults{Result: result}, nil
}

// beforeToolCallback runs before a tool is executed
// It can modify tool arguments or skip tool execution entirely
func beforeToolCallback(ctx tool.Context, tool tool.Tool, args map[string]any) (map[string]any, error) {
	toolName := tool.Name()
	fmt.Printf("[Callback] Before tool call for '%s'\n", toolName)
	fmt.Printf("[Callback] Original args: %v\n", args)
	fmt.Printf("[Callback] Args type: %T\n", args)
	if args != nil {
		fmt.Printf("[Callback] Args length: %d\n", len(args))
		for k, v := range args {
			fmt.Printf("[Callback]   %s: %v (type: %T)\n", k, v, v)
		}
	}

	// Get the country argument
	var country string
	if countryVal, ok := args["country"]; ok {
		if countryStr, ok := countryVal.(string); ok {
			country = countryStr
			fmt.Printf("[Callback] Found country argument: '%s'\n", country)
		} else {
			fmt.Printf("[Callback] Country argument exists but is not string: %T\n", countryVal)
		}
	} else {
		fmt.Printf("[Callback] No 'country' argument found in args\n")
	}

	// If someone asks about 'Merica, convert to United States
	if toolName == "get_capital_city" && strings.ToLower(country) == "merica" {
		fmt.Println("[Callback] Converting 'Merica to 'United States'")
		args["country"] = "United States"
		fmt.Printf("[Callback] Modified args: %v\n", args)
		// Return nil to proceed with modified arguments
		return nil, nil
	}

	// Skip the call completely for restricted countries
	if toolName == "get_capital_city" && strings.ToLower(country) == "restricted" {
		fmt.Println("[Callback] Blocking restricted country")
		return map[string]any{"result": "Access to this information has been restricted."}, nil
	}

	fmt.Println("[Callback] Proceeding with normal tool call")
	// Return nil to proceed with normal tool call
	return nil, nil
}

// afterToolCallback runs after a tool has executed
// It can modify the tool response or handle errors
func afterToolCallback(ctx tool.Context, tool tool.Tool, args, result map[string]any, err error) (map[string]any, error) {
	toolName := tool.Name()
	fmt.Printf("[Callback] After tool call for '%s'\n", toolName)
	fmt.Printf("[Callback] Args used: %v\n", args)
	fmt.Printf("[Callback] Original response: %v\n", result)
	if err != nil {
		fmt.Printf("[Callback] Error: %v\n", err)
	}

	// Detailed result debugging
	if result == nil {
		fmt.Printf("[Callback] Result is nil!\n")
	} else {
		fmt.Printf("[Callback] Result type: %T, length: %d\n", result, len(result))
		for k, v := range result {
			fmt.Printf("[Callback]   Result %s: %v (type: %T)\n", k, v, v)
		}
	}

	// Extract the result - try both "result" and "Result" keys
	var originalResult string
	if result != nil {
		// Try lowercase "result" first
		if resultVal, ok := result["result"]; ok {
			if resultStr, ok := resultVal.(string); ok {
				originalResult = resultStr
				fmt.Printf("[Callback] Found result using 'result' key: '%s'\n", originalResult)
			} else {
				fmt.Printf("[Callback] 'result' key exists but is not string: %T\n", resultVal)
			}
		} else if resultVal, ok := result["Result"]; ok {
			// Try uppercase "Result" (from struct)
			if resultStr, ok := resultVal.(string); ok {
				originalResult = resultStr
				fmt.Printf("[Callback] Found result using 'Result' key: '%s'\n", originalResult)
			} else {
				fmt.Printf("[Callback] 'Result' key exists but is not string: %T\n", resultVal)
			}
		} else {
			fmt.Printf("[Callback] No 'result' or 'Result' key found in result map\n")
		}
	}

	fmt.Printf("[Callback] Extracted result: '%s'\n", originalResult)

	// Add a note for any USA capital responses
	if toolName == "get_capital_city" && strings.Contains(strings.ToLower(originalResult), "washington") {
		fmt.Println("[Callback] DETECTED USA CAPITAL - adding patriotic note!")

		// Create a modified copy of the response - only modify the result field
		modifiedResponse := map[string]any{}
		for k, v := range result {
			modifiedResponse[k] = v
		}
		// Update both possible keys to be safe
		modifiedResponse["result"] = fmt.Sprintf("%s (Note: This is the capital of the USA. 🇺🇸)", originalResult)
		modifiedResponse["Result"] = fmt.Sprintf("%s (Note: This is the capital of the USA. 🇺🇸)", originalResult)

		fmt.Printf("[Callback] Modified response: %v\n", modifiedResponse)
		return modifiedResponse, nil
	}

	fmt.Println("[Callback] No modifications needed, returning original response")
	// Return nil to use original response
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

	// Create the tool from the function
	getCapitalCityTool, err := functiontool.New(
		functiontool.Config{
			Name:        "get_capital_city",
			Description: "Retrieves the capital city of a given country",
		},
		getCapitalCity)
	if err != nil {
		log.Fatalf("Failed to create get_capital_city tool: %v", err)
	}

	// Create the agent with before and after tool callbacks
	a, err := llmagent.New(llmagent.Config{
		Name:        "tool_callback_agent",
		Model:       model,
		Description: "An agent that demonstrates tool callbacks by looking up capital cities",
		Instruction: `You are a helpful geography assistant.

Your job is to:
- Find capital cities when asked using the get_capital_city tool
- Use the exact country name provided by the user
- ALWAYS return the EXACT result from the tool, without changing it
- When reporting a capital, display it EXACTLY as returned by the tool

Examples:
- "What is the capital of France?" → Use get_capital_city with country="France"
- "Tell me the capital city of Japan" → Use get_capital_city with country="Japan"`,
		Tools:                []tool.Tool{getCapitalCityTool},
		BeforeToolCallbacks:  []llmagent.BeforeToolCallback{beforeToolCallback},
		AfterToolCallbacks:   []llmagent.AfterToolCallback{afterToolCallback},
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