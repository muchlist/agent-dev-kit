# Sessions and State Management in ADK (Go)

This example demonstrates how to create and manage stateful sessions in the Agent Development Kit (ADK) using Go, enabling your agents to maintain context and remember user information across interactions.

## What Are Sessions in ADK?

Sessions in ADK provide a way to:

1. **Maintain State**: Store and access user data, preferences, and other information between interactions
2. **Track Conversation History**: Automatically record and retrieve message history
3. **Personalize Responses**: Use stored information to create more contextual and personalized agent experiences

Unlike simple conversational agents that forget previous interactions, stateful agents can build relationships with users over time by remembering important details and preferences.

## Example Overview

This example demonstrates:

- Creating a session with user preferences
- Using template variables to access session state in agent instructions
- Running the agent with a session to maintain context using the Go ADK Runner
- Retrieving session state and message history

The example uses a simple question-answering agent that responds based on stored user information in the session state.

## Project Structure

```
5-sessions-and-state/
│
└── question_answering_agent/      # Agent implementation
    ├── main.go                    # Main example with session management
    └── .env.example               # Environment template
```

## Getting Started

### Prerequisites

1. Go 1.25 or higher installed
2. Google API key for Gemini

### Setup

1. Set up your API key:
   - Copy `.env.example` to `.env` in the question_answering_agent folder
   ```bash
   cd 5-sessions-and-state/question_answering_agent
   cp .env.example .env
   ```
   - Add your Google API key to the `GOOGLE_API_KEY` variable in the `.env` file

2. Load environment variables:
   ```bash
   # macOS/Linux:
   export $(cat .env | xargs)

   # Windows PowerShell:
   Get-Content .env | ForEach-Object {
       $name, $value = $_.split('=')
       Set-Item -Path env:$name -Value $value
   }
   ```

## Running the Example

### Method 1: Direct Execution

Run the example to see a stateful session in action:

```bash
cd 5-sessions-and-state/question_answering_agent
go run main.go
```

This will:
1. Create a new session with user information (name and preferences)
2. Initialize the agent with access to that session via template variables
3. Process a user query about the stored preferences
4. Display the agent's response based on the session data
5. Show the final session state and message history

### Method 2: Using Make (from root directory)

```bash
make run/5
```

## Key Components

### 1. Session Service

The example uses `session.InMemoryService()` which stores sessions in memory:

```go
sessionService := session.InMemoryService()
```

### 2. Initial State

Sessions are created with an initial state containing user information:

```go
initialState := map[string]any{
    "user_name": "Muchlis",
    "user_preferences": `
        I like to play Pickleball, Disc Golf, and Tennis.
        My favorite food is Mexican.
        My favorite TV show is Game of Thrones.
        Loves it when people like and subscribe to his YouTube channel.
    `,
}
```

### 3. Creating a Session

The example creates a session with a unique identifier:

```go
SESSION_ID := uuid.New().String()
createResp, err := sessionService.Create(ctx, &session.CreateRequest{
    AppName:   APP_NAME,
    UserID:    USER_ID,
    SessionID: SESSION_ID,
    State:     initialState,
})
```

### 4. Accessing State in Agent Instructions

The agent accesses session state using template variables in its instructions:

```go
Instruction: `You are a helpful assistant that answers questions about the user's preferences.

Here is some information about the user:
Name:
{user_name}
Preferences:
{user_preferences}`,
```

**Template Variable Syntax:**
- `{variable_name}` - Required variable, error if missing
- `{variable_name?}` - Optional variable, empty string if missing
- `{user:variable}` - User-scoped state (shared across sessions)
- `{app:variable}` - App-scoped state (shared across all users)
- `{temp:variable}` - Temporary state (not persisted)

### 5. Running with Sessions

Sessions are integrated with the `Runner` to maintain state between interactions:

```go
r, err := runner.New(runner.Config{
    AppName:        APP_NAME,
    Agent:          questionAnsweringAgent,
    SessionService: sessionService,
})

// Run the agent with session context
for event, err := range r.Run(ctx, USER_ID, SESSION_ID, userMessage, agent.RunConfig{}) {
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    // Process events
}
```

### 6. Retrieving Session State

After agent execution, you can retrieve the updated session:

```go
getResp, err := sessionService.Get(ctx, &session.GetRequest{
    AppName:   APP_NAME,
    UserID:    USER_ID,
    SessionID: SESSION_ID,
})

session := getResp.Session

// Access state
for key, value := range session.State().All() {
    fmt.Printf("%s: %v\n", key, value)
}

// Access message history
for event := range session.Events().All() {
    // Process each event
}
```

## State Scopes in Go ADK

Go ADK supports multiple state scopes using key prefixes:

### 1. Session-Specific State (Default)

State without prefixes is specific to each session:

```go
State: map[string]any{
    "user_name": "Alice",  // Only in this session
}
```

### 2. User-Scoped State (`user:` prefix)

State shared across all sessions for the same user:

```go
State: map[string]any{
    "user:preferences": "dark_mode",  // Shared for this user
}
```

### 3. App-Scoped State (`app:` prefix)

State shared across all users and sessions:

```go
State: map[string]any{
    "app:version": "1.0",  // Shared globally
}
```

### 4. Temporary State (`temp:` prefix)

State that is not persisted to storage (cleared after invocation):

```go
State: map[string]any{
    "temp:working_data": "...",  // Not saved
}
```

### Example with All Scopes:

```go
initialState := map[string]any{
    "session_var":      "session-specific",
    "user:theme":       "dark",
    "app:api_version":  "v2",
    "temp:cache_key":   "tmp123",
}
```

In agent instructions, access any scope:
```
Session: {session_var}
User Theme: {user:theme}
API Version: {app:api_version}
Temp Data: {temp:cache_key}
```

## Understanding the Runner

The `Runner` is the core component for executing agents with session context:

### Runner Configuration:

```go
type Config struct {
    AppName         string              // Application identifier
    Agent           agent.Agent         // The root agent
    SessionService  session.Service     // Session storage service
    ArtifactService artifact.Service    // Optional: artifact storage
    MemoryService   memory.Service      // Optional: RAG memory
}
```

### Runner Workflow:

1. **Retrieves** existing session from SessionService
2. **Finds** appropriate agent to run based on session history
3. **Creates** invocation context with session state
4. **Appends** user message as an event to session
5. **Executes** agent with full session context
6. **Persists** agent responses and state changes back to session

### Event Streaming:

The `Run` method returns an iterator that streams events:

```go
for event, err := range r.Run(ctx, userID, sessionID, message, config) {
    if err != nil {
        // Handle error
    }

    // Event types:
    // - User message events
    // - Agent response events (partial and final)
    // - Tool execution events
    // - State update events
}
```

## Example Output

```
CREATED NEW SESSION:
	Session ID: 550e8400-e29b-41d4-a716-446655440000
	App Name: Bot
	User ID: muchlis

User Question: What is Muchlis's favorite TV show?

Final Response: Based on the information provided, Muchlis's favorite TV show is Game of Thrones.

==== Session Event Exploration ====
=== Final Session State ===
user_name: muchlis
user_preferences:
        I like to play Pickleball, Disc Golf, and Tennis.
        My favorite food is Mexican.
        My favorite TV show is Game of Thrones.
        Loves it when people like and subscribe to his YouTube channel.


=== Session Message History ===
[1] user: What is Muchlis's favorite TV show?
[2] model: Based on the information provided, Muchlis's favorite TV show is Game of Thrones.

Example completed successfully!
```

## Comparison: Python vs Go

| Aspect | Python | Go |
|--------|--------|-----|
| **Session Service** | `InMemorySessionService()` | `session.InMemoryService()` |
| **Create Session** | `create_session(app_name, user_id, session_id, state)` | `Create(ctx, &session.CreateRequest{...})` |
| **Runner** | `Runner(agent, app_name, session_service)` | `runner.New(runner.Config{...})` |
| **Run Agent** | `for event in runner.run(user_id, session_id, message)` | `for event, err := range r.Run(ctx, userID, sessionID, msg, cfg)` |
| **Template Variables** | `{user_name}` | `{user_name}` (same syntax) |
| **Access State** | `session.state['key']` | `session.State().Get("key")` |
| **State Scopes** | `app:`, `user:`, `temp:` | `app:`, `user:`, `temp:` (same) |

### Python Example:
```python
session_service = InMemorySessionService()

session = session_service.create_session(
    app_name="my_app",
    user_id="user1",
    session_id="session1",
    state={"user_name": "Alice"}
)

runner = Runner(
    agent=my_agent,
    app_name="my_app",
    session_service=session_service
)

for event in runner.run(user_id, session_id, message):
    # Process events
```

### Go Example:
```go
sessionService := session.InMemoryService()

createResp, err := sessionService.Create(ctx, &session.CreateRequest{
    AppName:   "my_app",
    UserID:    "user1",
    SessionID: "session1",
    State:     map[string]any{"user_name": "Alice"},
})

r, err := runner.New(runner.Config{
    AppName:        "my_app",
    Agent:          myAgent,
    SessionService: sessionService,
})

for event, err := range r.Run(ctx, userID, sessionID, message, config) {
    // Process events
}
```

## Advanced: Modifying State During Execution

You can modify session state during agent execution by returning state deltas in events:

```go
// In a custom tool or callback
event := session.NewEvent(invocationID)
event.Actions.StateDelta = map[string]any{
    "last_query": "updated value",
    "user:visit_count": visitCount + 1,
}
```

When the event is appended to the session, the state is automatically merged.

## Session Persistence

By default, `InMemoryService()` stores sessions in memory. For production use:

1. **Implement `session.Service` interface** with persistent storage (Redis, PostgreSQL, etc.)
2. **Use session IDs consistently** across requests
3. **Handle session expiration** with TTLs
4. **Consider storage costs** for long conversation histories

## Use Cases

1. **Personalized Chatbots**: Remember user preferences, context, and history
2. **Multi-Turn Conversations**: Maintain context across multiple interactions
3. **User Profiles**: Store and retrieve user-specific settings
4. **Stateful Workflows**: Track progress through multi-step processes
5. **Collaborative Agents**: Share state between multiple agents in a workflow

## Additional Resources

**Go ADK Documentation:**
- [Session Package Documentation](https://pkg.go.dev/google.golang.org/adk/session)
- [Runner Package Documentation](https://pkg.go.dev/google.golang.org/adk/runner)
- [Go ADK GitHub Repository](https://github.com/google/adk-go)

**Python ADK (Reference):**
- [Google ADK Sessions Documentation](https://google.github.io/adk-docs/sessions/session/)
- [State Management in ADK](https://google.github.io/adk-docs/sessions/state/)

## Next Steps

Try modifying the example to:
- Add more user preferences to the initial state
- Create multiple sessions for different users
- Update state during the conversation
- Implement user-scoped and app-scoped state
- Build a multi-turn conversation flow
