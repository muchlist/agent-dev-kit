// Package main demonstrates persistent storage with database session management in ADK.
// This example creates a reminder agent that remembers user information across sessions.
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/adk/session/database"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

const (
	APP_NAME   = "Memory Agent"
	MODEL_NAME = "gemini-2.0-flash"
	DB_FILE    = "./my_agent_data.db"
)

// ===== Tool Argument and Result Structures =====

type addReminderArgs struct {
	Reminder string `json:"reminder"`
}

type addReminderResults struct {
	Action   string `json:"action"`
	Reminder string `json:"reminder"`
	Message  string `json:"message"`
}

type viewRemindersArgs struct{}

type viewRemindersResults struct {
	Action    string   `json:"action"`
	Reminders []string `json:"reminders"`
	Count     int      `json:"count"`
}

type updateReminderArgs struct {
	Index       int    `json:"index"`
	UpdatedText string `json:"updated_text"`
}

type updateReminderResults struct {
	Action      string `json:"action"`
	Status      string `json:"status,omitempty"`
	Index       int    `json:"index,omitempty"`
	OldText     string `json:"old_text,omitempty"`
	UpdatedText string `json:"updated_text,omitempty"`
	Message     string `json:"message"`
}

type deleteReminderArgs struct {
	Index int `json:"index"`
}

type deleteReminderResults struct {
	Action          string `json:"action"`
	Status          string `json:"status,omitempty"`
	Index           int    `json:"index,omitempty"`
	DeletedReminder string `json:"deleted_reminder,omitempty"`
	Message         string `json:"message"`
}

type updateUserNameArgs struct {
	Name string `json:"name"`
}

type updateUserNameResults struct {
	Action  string `json:"action"`
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
	Message string `json:"message"`
}

// ===== Tool Implementations =====

// Note: Go ADK tools access session state using ctx.State(), similar to Python's tool_context.state

func addReminder(ctx tool.Context, input addReminderArgs) (addReminderResults, error) {
	fmt.Printf("--- Tool: add_reminder called for '%s' ---\n", input.Reminder)

	// Access session state using ctx.State()
	state := ctx.State()

	// Get current reminders from state using the proper Get() method
	reminders := getRemindersList(state)

	// Add new reminder
	reminders = append(reminders, input.Reminder)

	// Update state using Set() method - changes are persisted automatically
	state.Set("reminders", reminders)

	return addReminderResults{
		Action:   "add_reminder",
		Reminder: input.Reminder,
		Message:  fmt.Sprintf("Added reminder: %s", input.Reminder),
	}, nil
}

func viewReminders(ctx tool.Context, input viewRemindersArgs) (viewRemindersResults, error) {
	fmt.Println("--- Tool: view_reminders called ---")

	// Access session state using ctx.State()
	state := ctx.State()

	// Get reminders from state using the proper Get() method
	reminders := getRemindersList(state)
	count := len(reminders)

	return viewRemindersResults{
		Action:    "view_reminders",
		Reminders: reminders,
		Count:     count,
	}, nil
}

func updateReminder(ctx tool.Context, input updateReminderArgs) (updateReminderResults, error) {
	fmt.Printf("--- Tool: update_reminder called for index %d with '%s' ---\n", input.Index, input.UpdatedText)

	// Access session state using ctx.State()
	state := ctx.State()

	// Get current reminders from state using the proper Get() method
	reminders := getRemindersList(state)

	// Check if index is valid and update reminder
	if input.Index >= 1 && input.Index <= len(reminders) {
		oldReminder := reminders[input.Index-1]
		reminders[input.Index-1] = input.UpdatedText

		// Update state using Set() method - changes are persisted automatically
		state.Set("reminders", reminders)

		return updateReminderResults{
			Action:      "update_reminder",
			Index:       input.Index,
			OldText:     oldReminder,
			UpdatedText: input.UpdatedText,
			Message:     fmt.Sprintf("Updated reminder %d from '%s' to '%s'", input.Index, oldReminder, input.UpdatedText),
		}, nil
	}

	return updateReminderResults{
		Action:      "update_reminder",
		Index:       input.Index,
		UpdatedText: input.UpdatedText,
		Message:     fmt.Sprintf("Could not find reminder at position %d. Currently there are %d reminders.", input.Index, len(reminders)),
	}, nil
}

func deleteReminder(ctx tool.Context, input deleteReminderArgs) (deleteReminderResults, error) {
	fmt.Printf("--- Tool: delete_reminder called for index %d ---\n", input.Index)

	// Access session state using ctx.State()
	state := ctx.State()

	// Get current reminders from state using the proper Get() method
	reminders := getRemindersList(state)

	// Check if index is valid and delete reminder
	if input.Index >= 1 && input.Index <= len(reminders) {
		deletedReminder := reminders[input.Index-1]

		// Remove the reminder
		reminders = append(reminders[:input.Index-1], reminders[input.Index:]...)

		// Update state using Set() method - changes are persisted automatically
		state.Set("reminders", reminders)

		return deleteReminderResults{
			Action:          "delete_reminder",
			Index:           input.Index,
			DeletedReminder: deletedReminder,
			Message:         fmt.Sprintf("Deleted reminder %d: '%s'", input.Index, deletedReminder),
		}, nil
	}

	return deleteReminderResults{
		Action:  "delete_reminder",
		Index:   input.Index,
		Message: fmt.Sprintf("Could not find reminder at position %d. Currently there are %d reminders.", input.Index, len(reminders)),
	}, nil
}

func updateUserName(ctx tool.Context, input updateUserNameArgs) (updateUserNameResults, error) {
	fmt.Printf("--- Tool: update_user_name called with '%s' ---\n", input.Name)

	// Access session state using ctx.State()
	state := ctx.State()

	// Get current name from state using the proper Get() method
	var oldName string
	if val, err := state.Get("user_name"); err == nil {
		if str, ok := val.(string); ok {
			oldName = str
		}
	}

	// Update state using Set() method - changes are persisted automatically
	state.Set("user_name", input.Name)

	return updateUserNameResults{
		Action:  "update_user_name",
		OldName: oldName,
		NewName: input.Name,
		Message: fmt.Sprintf("Updated your name from '%s' to: %s", oldName, input.Name),
	}, nil
}

// ===== Utility Functions =====

func getRemindersList(state session.ReadonlyState) []string {
	reminders := []string{}
	if val, err := state.Get("reminders"); err == nil {
		if remindersList, ok := val.([]interface{}); ok {
			for _, r := range remindersList {
				if str, ok := r.(string); ok {
					reminders = append(reminders, str)
				}
			}
		}
	}
	return reminders
}

func displayState(sessionService session.Service, appName, userID, sessionID, label string) {
	ctx := context.Background()
	getResp, err := sessionService.Get(ctx, &session.GetRequest{
		AppName:   appName,
		UserID:    userID,
		SessionID: sessionID,
	})
	if err != nil {
		fmt.Printf("Error displaying state: %v\n", err)
		return
	}

	sess := getResp.Session
	state := sess.State()

	fmt.Printf("\n---------- %s ----------\n", label)

	// Display user name
	userName := "Unknown"
	if val, err := state.Get("user_name"); err == nil {
		if str, ok := val.(string); ok {
			userName = str
		}
	}
	fmt.Printf("👤 User: %s\n", userName)

	// Display reminders
	reminders := []string{}
	if val, err := state.Get("reminders"); err == nil {
		if remindersList, ok := val.([]interface{}); ok {
			for _, r := range remindersList {
				if str, ok := r.(string); ok {
					reminders = append(reminders, str)
				}
			}
		}
	}

	if len(reminders) > 0 {
		fmt.Println("📝 Reminders:")
		for idx, reminder := range reminders {
			fmt.Printf("  %d. %s\n", idx+1, reminder)
		}
	} else {
		fmt.Println("📝 Reminders: None")
	}

	fmt.Printf("--%s--\n", strings.Repeat("-", len(label)+20))
}

// ===== Main Function =====

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

	// Create database session service with SQLite
	sessionService, err := database.NewSessionService(
		sqlite.Open(DB_FILE),
		&gorm.Config{
			PrepareStmt: true,
			Logger:      logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		log.Fatalf("Failed to create database session service: %v", err)
	}

	// Initialize database schema
	if err := database.AutoMigrate(sessionService); err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	fmt.Println("✅ Connected to database:", DB_FILE)

	// Create reminder management tools
	addReminderTool, err := functiontool.New(
		functiontool.Config{
			Name:        "add_reminder",
			Description: "Add a new reminder to the user's reminder list",
		},
		addReminder)
	if err != nil {
		log.Fatalf("Failed to create add_reminder tool: %v", err)
	}

	viewRemindersTool, err := functiontool.New(
		functiontool.Config{
			Name:        "view_reminders",
			Description: "View all current reminders",
		},
		viewReminders)
	if err != nil {
		log.Fatalf("Failed to create view_reminders tool: %v", err)
	}

	updateReminderTool, err := functiontool.New(
		functiontool.Config{
			Name:        "update_reminder",
			Description: "Update an existing reminder",
		},
		updateReminder)
	if err != nil {
		log.Fatalf("Failed to create update_reminder tool: %v", err)
	}

	deleteReminderTool, err := functiontool.New(
		functiontool.Config{
			Name:        "delete_reminder",
			Description: "Delete a reminder",
		},
		deleteReminder)
	if err != nil {
		log.Fatalf("Failed to create delete_reminder tool: %v", err)
	}

	updateUserNameTool, err := functiontool.New(
		functiontool.Config{
			Name:        "update_user_name",
			Description: "Update the user's name",
		},
		updateUserName)
	if err != nil {
		log.Fatalf("Failed to create update_user_name tool: %v", err)
	}

	// Create the memory agent
	memoryAgent, err := llmagent.New(llmagent.Config{
		Name:        "memory_agent",
		Model:       model,
		Description: "A smart reminder agent with persistent memory",
		Instruction: `You are a friendly reminder assistant that remembers users across conversations.

You have access to tools to manage reminders and user information.

You can help users manage their reminders with the following capabilities:
1. Add new reminders
2. View existing reminders
3. Update reminders
4. Delete reminders
5. Update the user's name

Always be friendly and address the user by name. If you don't know their name yet,
use the update_user_name tool to store it when they introduce themselves.

**REMINDER MANAGEMENT GUIDELINES:**

When dealing with reminders, you need to be smart about finding the right reminder:

1. When the user asks to update or delete a reminder but doesn't provide an index:
   - If they mention the content of the reminder (e.g., "delete my meeting reminder"),
     look through the reminders to find a match
   - If you find an exact or close match, use that index
   - Never ask for clarification, just use the first match
   - If no match is found, list all reminders and ask the user to specify

2. When the user mentions a number or position:
   - Use that as the index (e.g., "delete reminder 2" means index=2)
   - Remember that indexing starts at 1 for the user

3. For relative positions:
   - Handle "first", "last", "second", etc. appropriately
   - "First reminder" = index 1
   - "Last reminder" = the highest index
   - "Second reminder" = index 2, and so on

4. For viewing:
   - Always use the view_reminders tool when the user asks to see their reminders
   - IMPORTANT: The tool result may not contain the actual reminder data
   - Use the current session state information that is displayed before/after processing
   - Format the response in a numbered list for clarity
   - If there are no reminders, suggest adding some

5. For addition:
   - Extract the actual reminder text from the user's request
   - Remove phrases like "add a reminder to" or "remind me to"
   - Focus on the task itself (e.g., "add a reminder to buy milk" → add_reminder("buy milk"))

6. For updates:
   - Identify both which reminder to update and what the new text should be
   - For example, "change my second reminder to pick up groceries" → update_reminder(2, "pick up groceries")

7. For deletions:
   - Confirm deletion when complete and mention which reminder was removed
   - For example, "I've deleted your reminder to 'buy milk'"

Remember to explain that you can remember their information across conversations.

IMPORTANT:
- Use your best judgement to determine which reminder the user is referring to
- You don't have to be 100% correct, but try to be as close as possible
- Never ask the user to clarify which reminder they are referring to`,
		Tools: []tool.Tool{
			addReminderTool,
			viewRemindersTool,
			updateReminderTool,
			deleteReminderTool,
			updateUserNameTool,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Setup user and check for existing sessions
	USER_ID := "user_" + os.Getenv("USER")
	if USER_ID == "user_" {
		USER_ID = "default_user"
	}

	// Check for existing sessions
	listResp, err := sessionService.List(ctx, &session.ListRequest{
		AppName: APP_NAME,
		UserID:  USER_ID,
	})
	if err != nil {
		log.Fatalf("Failed to list sessions: %v", err)
	}

	var SESSION_ID string
	if len(listResp.Sessions) > 0 {
		// Use the most recent session
		SESSION_ID = listResp.Sessions[0].ID()
		fmt.Printf("🔄 Continuing existing session: %s\n", SESSION_ID)
	} else {
		// Create a new session with initial state
		initialState := map[string]any{
			"user_name": "User",
			"reminders": []string{},
		}
		createResp, err := sessionService.Create(ctx, &session.CreateRequest{
			AppName: APP_NAME,
			UserID:  USER_ID,
			State:   initialState,
		})
		if err != nil {
			log.Fatalf("Failed to create session: %v", err)
		}
		SESSION_ID = createResp.Session.ID()
		fmt.Printf("✨ Created new session: %s\n", SESSION_ID)
	}

	// Create runner with the memory agent
	r, err := runner.New(runner.Config{
		AppName:        APP_NAME,
		Agent:          memoryAgent,
		SessionService: sessionService,
	})
	if err != nil {
		log.Fatalf("Failed to create runner: %v", err)
	}

	// Interactive conversation loop
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Welcome to Memory Agent Chat!")
	fmt.Println("Your reminders will be remembered across conversations.")
	fmt.Println("Type 'exit' or 'quit' to end the conversation.")
	fmt.Println(strings.Repeat("=", 60) + "\n")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}

		userInput := strings.TrimSpace(scanner.Text())

		if userInput == "" {
			continue
		}

		// Check if user wants to exit
		if strings.ToLower(userInput) == "exit" || strings.ToLower(userInput) == "quit" {
			fmt.Println("\nEnding conversation. Your data has been saved to the database.")
			break
		}

		// Display state before processing
		displayState(sessionService, APP_NAME, USER_ID, SESSION_ID, "State BEFORE processing")

		// Create user message
		userMessage := &genai.Content{
			Role: "user",
			Parts: []*genai.Part{
				{Text: userInput},
			},
		}

		// Run the agent
		fmt.Printf("\n--- Running Query: %s ---\n", userInput)
		var finalResponse string

		for event, err := range r.Run(ctx, USER_ID, SESSION_ID, userMessage, agent.RunConfig{}) {
			if err != nil {
				fmt.Printf("Error during agent run: %v\n", err)
				break
			}

			// Capture final response
			if event.Content != nil && len(event.Content.Parts) > 0 && event.Content.Parts[0].Text != "" {
				finalResponse = event.Content.Parts[0].Text
			}
		}

		// Display agent response
		if finalResponse != "" {
			fmt.Println("\n╔══ AGENT RESPONSE ══════════════════════════════════════")
			fmt.Println(finalResponse)
			fmt.Println("╚════════════════════════════════════════════════════════")
		}

		// Display state after processing
		displayState(sessionService, APP_NAME, USER_ID, SESSION_ID, "State AFTER processing")
		fmt.Println()
	}
}
