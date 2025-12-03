# Stateful Multi-Agent Systems in ADK (Go)

This example demonstrates how to create a stateful multi-agent system in ADK using Go, combining the power of persistent state management with specialized agent delegation. This approach creates intelligent agent systems that remember user information across interactions while leveraging specialized domain expertise.

## What is a Stateful Multi-Agent System?

A Stateful Multi-Agent System combines two powerful patterns:

1. **State Management**: Persisting information about users and conversations across interactions
2. **Multi-Agent Architecture**: Distributing tasks among specialized agents based on their expertise

The result is a sophisticated agent ecosystem that can:
- Remember user information and interaction history
- Route queries to the most appropriate specialized agent
- Provide personalized responses based on past interactions
- Maintain context across multiple agent delegations

This example implements a customer service system for an online course platform, where specialized agents handle different aspects of customer support while sharing a common state.

## Project Structure

```
8-stateful-multi-agent/
└── customer_service_agent/
    ├── main.go                     # Entry point with session management
    ├── agents/                     # Modular specialized agents
    │   ├── sales_agent.go          # Course sales + purchase tool
    │   ├── policy_agent.go         # Policies and guidelines
    │   ├── course_support_agent.go # Course content help
    │   └── order_agent.go          # Order history + refund tool
    ├── utils/                      # State management utilities
    │   └── state.go                # Display and update helpers
    ├── .env.example
    └── .env
```

## Key Features

### 1. **Modular Agent Organization**
Following the pattern from Example 7, each agent is organized in its own file:
- **Sales Agent**: Handles purchases, updates state with courses
- **Policy Agent**: Provides policy information
- **Course Support Agent**: Helps with course content (access-controlled)
- **Order Agent**: Shows purchase history, processes refunds

### 2. **Session State Management**
The system maintains state across interactions:
```go
initialState := map[string]any{
    "user_name":           "Brandon Hancock",
    "purchased_courses":   []any{},
    "interaction_history": []any{},
}
```

### 3. **State Sharing Across Agents**
All agents access the same session state:
- Sales agent updates `purchased_courses` when user buys a course
- Course support agent checks if user owns course before helping
- Order agent can refund courses (removes from `purchased_courses`)
- All interactions tracked in `interaction_history`

### 4. **In-Memory Session Service**
Uses `session.InMemoryService()` for demonstration:
- Fast, no external dependencies
- Perfect for development and testing
- Can be replaced with `database.NewSessionService()` for production

## Key Components

### Session Management

```go
// Create in-memory session service
sessionService := session.InMemoryService()

// Create session with initial state
createResp, err := sessionService.Create(ctx, &session.CreateRequest{
    AppName: APP_NAME,
    UserID:  USER_ID,
    State:   initialState,
})
```

### State Updates in Tools

```go
// Get state from tool context
state := ctx.State()

// Read from state
var purchasedCourses []Course
if val, err := state.Get("purchased_courses"); err == nil {
    // Process courses...
}

// Update state
state.Set("purchased_courses", updatedCourses)
// Changes are automatically persisted!
```

### Runner Pattern

Unlike the launcher pattern in Example 7, this uses Runner for console interaction:

```go
r, err := runner.New(runner.Config{
    AppName:        APP_NAME,
    Agent:          customerServiceAgent,
    SessionService: sessionService,
})

// Run agent with message
for event, err := range r.Run(ctx, USER_ID, SESSION_ID, userMessage, agent.RunConfig{}) {
    // Process events...
}
```

## Getting Started

### Prerequisites

1. Go 1.21 or later
2. Google API key from https://aistudio.google.com/apikey

### Setup

1. Navigate to the customer_service_agent directory:
```bash
cd 8-stateful-multi-agent/customer_service_agent
```

2. Create your `.env` file:
```bash
cp .env.example .env
```

3. Edit `.env` and add your Google API key:
```env
GOOGLE_API_KEY=your_actual_api_key_here
```

## Running the Example

### Using Make (Recommended - from repository root)

```bash
make run/8
```

### Direct Execution

```bash
go run 8-stateful-multi-agent/customer_service_agent/main.go
```

## Example Conversation Flow

Try this conversation to see the stateful multi-agent system in action:

### 1. Start with General Inquiry
```
You: What courses do you offer?
```
*Manager routes to Sales Agent*

### 2. Purchase a Course
```
You: I want to buy the AI Marketing Platform course
```
*Sales Agent uses `purchase_course` tool, updates state*

### 3. Ask About Course Content
```
You: Can you tell me about section 10?
```
*Manager checks state, sees user owns course, routes to Course Support Agent*

### 4. Check Purchase History
```
You: What courses have I purchased?
```
*Manager routes to Order Agent, which reads from state*

### 5. Request a Refund
```
You: I'd like a refund for the course
```
*Order Agent uses `refund_course` tool, removes from state*

Notice how the system remembers your purchase across different agents!

## Agent Delegation Logic

### Customer Service Manager

Routes queries based on intent:
- **"purchase", "buy", "price"** → Sales Agent
- **"policy", "refund policy", "guidelines"** → Policy Agent
- **"course content", "section", "how do I"** → Course Support Agent (if owns course)
- **"order history", "refund", "purchased"** → Order Agent

### Access Control

Course Support Agent only helps if user owns the course:
```
Check if a course with id "ai_marketing_platform" exists in the purchased courses
before providing detailed help
```

## State Structure

### User Information
```go
"user_name": "Brandon Hancock"
```

### Purchased Courses
```go
"purchased_courses": [
    {
        "id": "ai_marketing_platform",
        "purchase_date": "2024-12-03 15:30:00"
    }
]
```

### Interaction History
```go
"interaction_history": [
    {
        "action": "user_query",
        "query": "What courses do you offer?",
        "timestamp": "2024-12-03 15:29:00"
    },
    {
        "action": "purchase_course",
        "course_id": "ai_marketing_platform",
        "timestamp": "2024-12-03 15:30:00"
    }
]
```

## Tools Overview

### Sales Agent Tools

**purchase_course**:
- Checks if user already owns course
- Adds course to `purchased_courses`
- Updates `interaction_history`
- Returns success/error status

### Order Agent Tools

**refund_course**:
- Verifies user owns the course
- Removes course from `purchased_courses`
- Updates `interaction_history`
- Returns success message

**get_current_time**:
- Returns current timestamp
- Used for order history queries

## Comparison with Python Version

| Feature | Python | Go (This Example) |
|---------|--------|-------------------|
| Structure | Multiple dirs/packages | Modular (agents/ + utils/) |
| Session Service | `InMemorySessionService()` | `session.InMemoryService()` |
| State Access | `tool_context.state["key"]` | `ctx.State().Set/Get()` |
| Runner | `Runner(agent, app_name, session_service)` | `runner.New(runner.Config{...})` |
| Execution | `runner.run_async()` | `r.Run()` with iterator |
| State Save | Manual via `create_session()` | Automatic on `Set()` |

## Extending the System

### Adding a New Agent

1. **Create agent file** in `agents/`:
```go
// agents/billing_agent.go
package agents

func NewBillingAgent(ctx context.Context, mdl model.LLM) (agent.Agent, error) {
    // Create agent with tools
}
```

2. **Add to main.go**:
```go
billingAgent, err := agents.NewBillingAgent(ctx, model)
// Add to customer service sub-agents
```

### Adding State Fields

Modify initial state in `main.go`:
```go
initialState := map[string]any{
    "user_name":           "Brandon Hancock",
    "purchased_courses":   []any{},
    "interaction_history": []any{},
    "subscription_tier":   "free",  // New field
    "last_login":          time.Now(),
}
```

### Using Persistent Storage

Replace `session.InMemoryService()` with database:

```go
import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "google.golang.org/adk/session/database"
)

// Create database session service
sessionService, err := database.NewSessionService(
    sqlite.Open("./customer_service.db"),
    &gorm.Config{PrepareStmt: true},
)

// Auto-migrate schema
if err := database.AutoMigrate(sessionService); err != nil {
    log.Fatalf("Failed to migrate: %v", err)
}
```

## Production Considerations

### 1. Persistent Storage
- Replace in-memory service with database
- Use PostgreSQL or SQLite
- Implement proper error handling

### 2. User Authentication
- Implement secure user identification
- Use proper USER_ID from auth system
- Validate session ownership

### 3. State Validation
- Validate state structure on read
- Handle missing/corrupted state gracefully
- Implement state migration for schema changes

### 4. Error Handling
- Handle agent failures gracefully
- Implement retry logic for transient errors
- Log errors for debugging

### 5. Monitoring
- Track agent selection patterns
- Monitor state size and growth
- Log interaction history for analytics

## Troubleshooting

### Common Issues

1. **"Failed to get session"**
   - Check if session was created successfully
   - Verify APP_NAME, USER_ID, and SESSION_ID match

2. **State not persisting**
   - With `session.InMemoryService()`, state only lasts while program runs
   - Use database service for true persistence

3. **Course support denies access**
   - Check if course exists in `purchased_courses`
   - Verify course ID is exactly `"ai_marketing_platform"`

4. **Agent not delegating correctly**
   - Review manager's `Instruction` field
   - Check agent `Description` fields
   - Ensure query matches delegation logic

### Debug Output

The system includes debug output:
```
--- Tool: purchase_course called ---
--- Tool: refund_course called ---
--- Tool: get_current_time called ---
```

State is displayed before and after each interaction:
```
---------- State BEFORE processing ----------
👤 User: Brandon Hancock
📚 Courses: None
📝 Interaction History: [...]
```

## Additional Resources

- [ADK Sessions Documentation](https://google.github.io/adk-docs/sessions/session/)
- [ADK Multi-Agent Systems](https://google.github.io/adk-docs/agents/multi-agent-systems/)
- [State Management in ADK](https://google.github.io/adk-docs/sessions/state/)
- [Go ADK API Reference](https://pkg.go.dev/google.golang.org/adk)

## Next Steps

After mastering stateful multi-agent systems, explore:
- **Example 9 - Callbacks**: Monitor agent events and add custom hooks
- **Example 10 - Sequential Agent**: Create pipeline workflows
- **Example 11 - Parallel Agent**: Concurrent agent operations
