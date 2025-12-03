# Persistent Storage in ADK (Go)

This example demonstrates how to implement persistent storage for your ADK agents using database-backed session management, allowing them to remember information and maintain conversation history across multiple sessions, application restarts, and even server deployments.

## What is Persistent Storage in ADK?

In previous examples, we used `session.InMemoryService()` which stores session data only in memory - this data is lost when the application stops. For real-world applications, you'll often need your agents to remember user information and conversation history long-term. This is where persistent storage comes in.

Go ADK provides the `database.NewSessionService()` that allows you to store session data in SQL databases, ensuring:

1. **Long-term Memory**: Information persists across application restarts
2. **Consistent User Experiences**: Users can continue conversations where they left off
3. **Multi-user Support**: Different users' data remains separate and secure
4. **Scalability**: Works with production databases for high-scale deployments

This example shows how to implement a reminder agent that remembers your name and todos across different conversations using an SQLite database.

## Project Structure

```
6-persistent-storage/
│
└── memory_agent/               # Agent package
    ├── main.go                 # Application with database session setup
    ├── .env.example            # Environment template
    └── my_agent_data.db        # SQLite database file (created on first run)
```

## Key Components

### 1. DatabaseSessionService

The core component that provides persistence is `database.NewSessionService()`, which is initialized with a GORM dialector:

#### Understanding "Record Not Found" Messages
When running the agent, you may see debug output like:
```
2025/12/03 17:43:18 /Users/muchlis-/go/pkg/mod/google.golang.org/adk@v0.2.0/session/database/service.go:433 record not found
[0.049ms] [rows:0] SELECT * FROM `app_states` WHERE app_name = "Memory Agent" ORDER BY `app_states`.`app_name` LIMIT 1

2025/12/03 17:43:18 /Users/muchlis-/go/pkg/mod/google.golang.org/adk@v0.2.0/session/database/service.go:445 record not found
[0.052ms] [rows:0] SELECT * FROM `user_states` WHERE app_name = "Memory Agent" AND user_id = "user_muchlis-" ORDER BY `user_states`.`app_name` LIMIT 1
```

**These messages are completely normal and harmless!** They indicate that:

- The session service is checking for app-level state cache entries in `app_states` table
- No records were found, which is expected for new installations
- The service then falls back to reading from the main `sessions` table where actual data is stored

This is **efficient caching behavior** by Google ADK:
1. **Fast cache lookup**: First checks `app_states` table for quick app-level state access
2. **Graceful fallback**: If cache miss, reads from main `sessions` table with complete data
3. **Performance optimization**: Reduces database load by avoiding complex joins on every request

**Your data is always safe** in the main `sessions` table, even when cache tables are empty!

```go
import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "google.golang.org/adk/session/database"
)

sessionService, err := database.NewSessionService(
    sqlite.Open("./my_agent_data.db"),
    &gorm.Config{PrepareStmt: true},
)
```

This service allows ADK to:
- Store session data in a SQL database
- Retrieve previous sessions for a user
- Automatically manage database schemas

### 2. Database Schema Initialization

After creating the session service, you must initialize the database schema:

```go
err := database.AutoMigrate(sessionService)
if err != nil {
    log.Fatalf("Failed to auto-migrate database: %v", err)
}
```

This creates four tables:
- `sessions` - Session-specific state and metadata
- `events` - Conversation events and message history
- `app_states` - App-level state (shared across all users)
- `user_states` - User-level state (shared across user's sessions)

### 3. Session Management

The example demonstrates proper session management:

```go
// Check for existing sessions for this user
listResp, err := sessionService.List(ctx, &session.ListRequest{
    AppName: APP_NAME,
    UserID:  USER_ID,
})

// If there's an existing session, use it, otherwise create a new one
if len(listResp.Sessions) > 0 {
    // Use the most recent session
    SESSION_ID = listResp.Sessions[0].ID()
    fmt.Printf("Continuing existing session: %s\n", SESSION_ID)
} else {
    // Create a new session with initial state
    createResp, err := sessionService.Create(ctx, &session.CreateRequest{
        AppName: APP_NAME,
        UserID:  USER_ID,
        State: map[string]any{
            "user_name": "User",
            "reminders": []string{},
        },
    })
    SESSION_ID = createResp.Session.ID()
    fmt.Printf("Created new session: %s\n", SESSION_ID)
}
```

### 4. State Management with Tools

The agent includes tools that update the persistent state:

```go
func addReminder(ctx tool.Context, input addReminderArgs) (addReminderResults, error) {
    // Get current reminders from state
    reminders := []string{}
    if val, err := ctx.Session().State().Get("reminders"); err == nil {
        // Parse existing reminders
        if remindersList, ok := val.([]interface{}); ok {
            for _, r := range remindersList {
                if str, ok := r.(string); ok {
                    reminders = append(reminders, str)
                }
            }
        }
    }

    // Add the new reminder
    reminders = append(reminders, input.Reminder)

    // Update state - automatically persisted to database
    ctx.Session().State().Set("reminders", reminders)

    return addReminderResults{
        Action:   "add_reminder",
        Reminder: input.Reminder,
        Message:  fmt.Sprintf("Added reminder: %s", input.Reminder),
    }, nil
}
```

Each change to `ctx.Session().State()` is automatically saved to the database when events are appended.

## Getting Started

### Prerequisites

1. Go 1.25 or higher installed
2. Google API key for Gemini
3. SQLite support (included with GORM driver)

### Dependencies

The example requires these additional Go modules:

```bash
go get gorm.io/driver/sqlite
go get gorm.io/gorm
```

These are used for database connectivity and ORM operations.

### Setup

1. Set up your API key:
   - Copy `.env.example` to `.env` in the memory_agent folder
   ```bash
   cd 6-persistent-storage/memory_agent
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

Run the example to start an interactive conversation:

```bash
cd 6-persistent-storage/memory_agent
go run main.go
```

This will:
1. Connect to the SQLite database (or create it if it doesn't exist)
2. Check for previous sessions for the user
3. Start an interactive conversation with the memory agent
4. Save all interactions to the database

### Method 2: Using Make (from root directory)

```bash
make run/6
```

### Getting Help

The agent responds to natural language queries about reminders:

Try these interactions to test the agent's persistent memory:

1. **First run:**
   ```
   You: What's my name?
   You: My name is John
   You: Add a reminder to buy groceries
   You: Add another reminder to finish the report
   You: What are my reminders?
   You: exit
   ```

2. **Second run (restart the program):**
   ```
   You: What's my name?
   You: What reminders do I have?
   You: Update my second reminder to submit the report by Friday
   You: Delete the first reminder
   You: exit
   ```

The agent will remember your name and reminders between runs!

## Database Tables Created

When you run the example, `database.AutoMigrate()` creates these tables in `my_agent_data.db`:

### 1. `sessions` Table
Stores session-specific data:
- `app_name`, `user_id`, `id` (composite primary key)
- `state` (JSON) - Session-specific state
- `create_time`, `update_time`

### 2. `events` Table
Stores conversation history:
- `id`, `app_name`, `user_id`, `session_id` (composite primary key)
- `invocation_id`, `author`, `timestamp`
- `content` (JSON) - Message content
- `actions` (JSON) - State changes and tool calls
- Metadata fields for grounding, usage, citations, etc.

### 3. `app_states` Table
Stores app-level state (shared across all users):
- `app_name` (primary key)
- `state` (JSON) - App-wide state with `app:` prefix
- `update_time`

### 4. `user_states` Table
Stores user-level state (shared across user's sessions):
- `app_name`, `user_id` (composite primary key)
- `state` (JSON) - User-specific state with `user:` prefix
- `update_time`

## State Scopes in Database Storage

The database session service supports multiple state scopes:

### 1. Session-Specific State (Default)
```go
State: map[string]any{
    "conversation_count": 5,
    "last_topic": "reminders",
}
```
Stored in the `sessions` table.

### 2. User-Scoped State (`user:` prefix)
```go
StateDelta: map[string]any{
    "user:theme": "dark_mode",
    "user:language": "en",
}
```
Stored in the `user_states` table, shared across all user sessions.

### 3. App-Scoped State (`app:` prefix)
```go
StateDelta: map[string]any{
    "app:total_users": 100,
    "app:version": "v2",
}
```
Stored in the `app_states` table, shared across all users.

### 4. Temporary State (`temp:` prefix)
```go
StateDelta: map[string]any{
    "temp:cache_key": "temp123",
}
```
Not persisted to database, discarded after invocation.

## Supported Database Backends

While this example uses SQLite for simplicity, the `database.NewSessionService()` supports various database backends through GORM:

### SQLite (Development)
```go
import "gorm.io/driver/sqlite"

sessionService, err := database.NewSessionService(
    sqlite.Open("./my_agent_data.db"),
    &gorm.Config{},
)
```

### PostgreSQL (Production)
```go
import "gorm.io/driver/postgres"

dsn := "host=localhost user=postgres password=secret dbname=adk_db port=5432 sslmode=disable"
sessionService, err := database.NewSessionService(
    postgres.Open(dsn),
    &gorm.Config{},
)
```

### MySQL (Production)
```go
import "gorm.io/driver/mysql"

dsn := "user:password@tcp(localhost:3306)/adk_db?charset=utf8mb4&parseTime=True"
sessionService, err := database.NewSessionService(
    mysql.Open(dsn),
    &gorm.Config{},
)
```

### Google Cloud Spanner (Cloud-Native)
```go
import "gorm.io/driver/spanner"

dsn := "projects/PROJECT_ID/instances/INSTANCE_ID/databases/DATABASE_ID"
sessionService, err := database.NewSessionService(
    spanner.New(spanner.Config{DSN: dsn}),
    &gorm.Config{},
)
```

## Example Output

```
✅ Connected to database: ./my_agent_data.db
✨ Created new session: a1b2c3d4-e5f6-7890-abcd-ef1234567890

============================================================
Welcome to Memory Agent Chat!
Your reminders will be remembered across conversations.
Type 'exit' or 'quit' to end the conversation.
============================================================

You: My name is John
--- Running Query: My name is John ---

---------- State BEFORE processing ----------
👤 User: User
📝 Reminders: None
--------------------------------------------

--- Tool: update_user_name called with 'John' ---

╔══ AGENT RESPONSE ══════════════════════════════════════
Nice to meet you, John! I've updated your name. How can I help you today?
╚════════════════════════════════════════════════════════

---------- State AFTER processing ----------
👤 User: John
📝 Reminders: None
--------------------------------------------

You: Add a reminder to buy groceries
--- Running Query: Add a reminder to buy groceries ---

---------- State BEFORE processing ----------
👤 User: John
📝 Reminders: None
--------------------------------------------

--- Tool: add_reminder called for 'buy groceries' ---

╔══ AGENT RESPONSE ══════════════════════════════════════
I've added the reminder "buy groceries" for you, John!
╚════════════════════════════════════════════════════════

---------- State AFTER processing ----------
👤 User: John
📝 Reminders:
  1. buy groceries
--------------------------------------------

You: What are my reminders?
--- Running Query: What are my reminders? ---

--- Tool: view_reminders called ---

╔══ AGENT RESPONSE ══════════════════════════════════════
You have 1 reminder, John:
1. buy groceries
╚════════════════════════════════════════════════════════

You: exit

Ending conversation. Your data has been saved to the database.
```

On the **second run**, the agent remembers everything:

```
✅ Connected to database: ./my_agent_data.db
🔄 Continuing existing session: a1b2c3d4-e5f6-7890-abcd-ef1234567890

You: What's my name?

╔══ AGENT RESPONSE ══════════════════════════════════════
Your name is John!
╚════════════════════════════════════════════════════════

You: What are my reminders?

--- Tool: view_reminders called ---

╔══ AGENT RESPONSE ══════════════════════════════════════
You have 1 reminder, John:
1. buy groceries
╚════════════════════════════════════════════════════════
```

## Comparison: Python vs Go

| Aspect | Python | Go |
|--------|--------|-----|
| **Database Service** | `DatabaseSessionService(db_url="...")` | `database.NewSessionService(dialector, config)` |
| **Database URL** | SQLAlchemy string format | GORM dialectors |
| **Schema Creation** | Automatic on first use | Explicit `database.AutoMigrate()` required |
| **List Sessions** | `list_sessions(app_name, user_id)` | `List(ctx, &session.ListRequest{...})` |
| **Create Session** | `create_session(app_name, user_id, state=...)` | `Create(ctx, &session.CreateRequest{...})` |
| **State Access** | `tool_context.state["key"] = value` | `ctx.Session().State().Set("key", value)` |
| **State Persistence** | Automatic | Automatic on event append |
| **Database Driver** | SQLAlchemy | GORM with database-specific drivers |

### Python Example:
```python
from google.adk.sessions import DatabaseSessionService

session_service = DatabaseSessionService(db_url="sqlite:///./my_agent_data.db")

sessions = session_service.list_sessions(
    app_name="my_app",
    user_id="user123"
)

if sessions.sessions:
    session_id = sessions.sessions[0].id
else:
    new_session = session_service.create_session(
        app_name="my_app",
        user_id="user123",
        state={"key": "value"}
    )
    session_id = new_session.id
```

### Go Example:
```go
import (
    "gorm.io/driver/sqlite"
    "google.golang.org/adk/session/database"
)

sessionService, err := database.NewSessionService(
    sqlite.Open("./my_agent_data.db"),
    &gorm.Config{},
)

database.AutoMigrate(sessionService)

listResp, err := sessionService.List(ctx, &session.ListRequest{
    AppName: "my_app",
    UserID:  "user123",
})

if len(listResp.Sessions) > 0 {
    sessionID = listResp.Sessions[0].ID()
} else {
    createResp, err := sessionService.Create(ctx, &session.CreateRequest{
        AppName: "my_app",
        UserID:  "user123",
        State:   map[string]any{"key": "value"},
    })
    sessionID = createResp.Session.ID()
}
```

## Production Considerations

For production use:

1. **Choose the Right Database**:
   - SQLite for single-instance applications
   - PostgreSQL/MySQL for multi-instance applications
   - Cloud Spanner for globally distributed applications

2. **Connection Pooling**:
   ```go
   db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
   sqlDB, _ := db.DB()
   sqlDB.SetMaxIdleConns(10)
   sqlDB.SetMaxOpenConns(100)
   sqlDB.SetConnMaxLifetime(time.Hour)
   ```

3. **Security**:
   - Use environment variables for database credentials
   - Enable SSL/TLS for database connections
   - Implement proper authentication and authorization
   - Use database-level encryption for sensitive data

4. **Backup and Recovery**:
   - Implement regular database backups
   - Test restore procedures
   - Consider point-in-time recovery for critical data

5. **Monitoring**:
   - Track database query performance
   - Monitor connection pool usage
   - Set up alerts for errors and slow queries

## Additional Resources

**Go ADK Documentation:**
- [Session Package Documentation](https://pkg.go.dev/google.golang.org/adk/session)
- [Database Session Service](https://pkg.go.dev/google.golang.org/adk/session/database)
- [Go ADK GitHub Repository](https://github.com/google/adk-go)

**GORM Documentation:**
- [GORM Official Documentation](https://gorm.io/docs/)
- [GORM Database Drivers](https://gorm.io/docs/connecting_to_the_database.html)

**Python ADK (Reference):**
- [ADK Sessions Documentation](https://google.github.io/adk-docs/sessions/session/)
- [Session Service Implementations](https://google.github.io/adk-docs/sessions/session/#sessionservice-implementations)

## Next Steps

Try modifying the example to:
- Switch from SQLite to PostgreSQL or MySQL
- Add more custom tools with state persistence
- Implement user-scoped and app-scoped state
- Build a multi-user application with isolated data
- Add session expiration policies
- Implement data export/import functionality
