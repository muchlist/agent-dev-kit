# Agent Development Kit (ADK) Crash Course - Go

A comprehensive hands-on guide to building LLM-powered agents using Google's Agent Development Kit in Go. This repository contains progressive examples that teach you how to create intelligent agents from basic interactions to complex multi-agent systems.

## What is ADK?

The Agent Development Kit (ADK) is Google's framework for building LLM-powered agents. Unlike traditional deterministic workflows, ADK agents use Large Language Models to:

- Understand natural language
- Reason about complex tasks
- Make dynamic decisions
- Interact with external tools
- Maintain conversation state

## Prerequisites

- Go 1.25 or higher
- Google API key for Gemini ([Get one here](https://aistudio.google.com/apikey))

## Quick Start

### 1. Clone the Repository

```bash
git clone [<repository-url>](https://github.com/muchlist/agent-dev-kit)
cd agent-dev-kit
```

### 2. Set Up Environment

```bash
# Copy the example environment file
cp .env.example .env

# Edit .env and add your Google API key
# GOOGLE_API_KEY=your_api_key_here
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run Your First Agent

```bash
# Using Makefile
make run/1

# Or directly
cd 1-basic-agent/greeting_agent
go run main.go web api webui
```

Open your browser to `http://localhost:8080` and start chatting with your agent!

## Examples Overview

### 1. Basic Agent
**Directory:** `1-basic-agent/greeting_agent`

Learn the fundamentals of ADK agents. This example introduces:
- Agent identity (name and description)
- Model selection
- Instructions
- Running agents in different modes (CLI, web, API)

**Try it:**
```bash
make run/1
# or
cd 1-basic-agent/greeting_agent && go run main.go web api webui
```

**Example prompts:**
- "Hello, what's your name?"
- "My name is Alice, can you greet me?"

---

### 2. Tool Agent
**Directory:** `2-tool-agent/tool_agent`

Enhance agents with tools to interact with the outside world. Learn about:
- Built-in tools (Google Search)
- Custom function tools
- Tool limitations in single agents

**Try it:**
```bash
make run/2
# or
cd 2-tool-agent/tool_agent && go run main.go web api webui
```

**Example prompts:**
- "What's the weather in Tokyo today?"
- "Find the latest news about artificial intelligence"

**Important:** Single agents can only use ONE built-in tool OR multiple custom function tools, but cannot mix them.

---

### 3. LiteLLM Agent
**Directory:** `3-litellm-agent/dad_joke_agent`

Use LiteLLM for provider abstraction. Learn how to:
- Abstract away LLM provider details
- Switch between different models easily
- Use Gemini through LiteLLM

**Try it:**
```bash
make run/3
# or
cd 3-litellm-agent/dad_joke_agent && go run main.go web api webui
```

**Example prompts:**
- "Tell me a dad joke about programming"
- "Give me a joke about coffee"

---

### 4. Structured Outputs
**Directory:** `4-structured-outputs/email_agent`

Ensure consistent, structured responses using schemas. Learn about:
- Defining output schemas with `genai.Schema`
- Guaranteeing JSON format responses
- Using `OutputKey` for structured data

**Try it:**
```bash
make run/4
# or
cd 4-structured-outputs/email_agent && go run main.go web api webui
```

**Example prompts:**
- "Write a professional email to schedule a meeting"
- "Create an email asking for project feedback"

---

### 5. Sessions and State
**Directory:** `5-sessions-and-state/question_answering_agent`

Maintain context across multiple interactions. Learn about:
- Session management
- State persistence
- In-memory session service
- Accessing state in tools using `ctx.State()`

**Try it:**
```bash
make run/5
# or
cd 5-sessions-and-state/question_answering_agent && go run main.go
```

**Example interaction:**
- "My name is John"
- (later) "What's my name?" → Agent remembers "John"

---

### 6. Persistent Storage
**Directory:** `6-persistent-storage/memory_agent`

Store agent data persistently across application restarts. Learn about:
- Database-backed sessions (SQLite with GORM)
- State serialization
- Database migrations
- Building a reminder management system

**Try it:**
```bash
make run/6
# or
cd 6-persistent-storage/memory_agent && go run main.go
```

**Example prompts:**
- "Add a reminder to buy groceries"
- "Show me my reminders"
- "Update my second reminder to call mom"
- "Delete the first reminder"

The agent remembers your data even after restarting!

---

## Running Modes

Each agent supports multiple running modes:

### Web UI (Recommended for learning)
```bash
go run main.go web api webui
```
Opens a browser interface at `http://localhost:8080`

### Command Line Interface
```bash
go run main.go run
```
Interactive chat in your terminal

### API Server
```bash
go run main.go api
```
REST API server at `http://localhost:8080`

### Help
```bash
go run main.go help
```
Show all available commands

## Project Structure

```
agent-development-kit-crash-course/
├── 1-basic-agent/
│   └── greeting_agent/
│       ├── main.go
│       └── .env.example
├── 2-tool-agent/
│   └── tool_agent/
│       ├── main.go
│       └── .env.example
├── 3-litellm-agent/
│   └── dad_joke_agent/
│       ├── main.go
│       └── .env.example
├── 4-structured-outputs/
│   └── email_agent/
│       ├── main.go
│       └── .env.example
├── 5-sessions-and-state/
│   └── question_answering_agent/
│       ├── main.go
│       └── .env.example
├── 6-persistent-storage/
│   └── memory_agent/
│       ├── main.go
│       └── .env.example
├── go.mod
├── go.sum
├── Makefile
├── .env.example
└── README.md
```

## Key Concepts

### Agent Anatomy

Every ADK agent consists of:

1. **Name** - Unique identifier for your agent
2. **Model** - LLM to power the agent (e.g., "gemini-2.0-flash")
3. **Instruction** - Defines behavior, personality, and constraints
4. **Description** - Summary of agent capabilities
5. **Tools** (Optional) - Functions the agent can call
6. **OutputSchema** (Optional) - Structure for consistent responses

### Tools in ADK

Tools allow agents to perform actions beyond text generation:

- **Built-in Tools**: Google Search, Code Execution (via `geminitool`)
- **Custom Function Tools**: Your own Go functions wrapped with `functiontool.New()`
- **Agent Tools**: Other agents wrapped as tools (multi-agent systems)

**Important Limitation:** A single agent can use only ONE built-in tool. To use multiple built-in tools, use a multi-agent architecture.

### Session State Management

State is accessed through the `tool.Context`:

```go
func myTool(ctx tool.Context, input Args) (Results, error) {
    state := ctx.State()

    // Get value
    val, err := state.Get("key")

    // Set value
    state.Set("key", value)

    return results, nil
}
```

## Common Issues & Solutions

### "Failed to create model: invalid API key"
- Make sure you've created a `.env` file in the agent directory
- Verify your API key is correct
- Ensure the API key is connected to a billing account

### "Cannot use built-in tool with custom tools"
- Single agents can only use ONE built-in tool OR multiple custom tools
- Solution: Use multi-agent architecture to combine different tool types

### Database locked errors
- Close any other instances accessing the database
- The SQLite database file is created at `my_agent_data.db`

## Dependencies

Main Go packages used:

```go
google.golang.org/adk          // Agent Development Kit
google.golang.org/genai        // Gemini API client
github.com/joho/godotenv       // Environment variables
gorm.io/gorm                   // ORM for database
gorm.io/driver/sqlite          // SQLite driver
```

## Learning Path

1. **Start with 1-basic-agent** - Understand agent fundamentals
2. **Move to 2-tool-agent** - Learn how agents interact with tools
3. **Try 3-litellm-agent** - See provider abstraction in action
4. **Explore 4-structured-outputs** - Get consistent JSON responses
5. **Practice 5-sessions-and-state** - Build stateful conversations
6. **Master 6-persistent-storage** - Create production-ready agents

## Resources

- [ADK Go Documentation](https://google.github.io/adk-docs/get-started/go/)
- [ADK Go Package Docs](https://pkg.go.dev/google.golang.org/adk)
- [Gemini API Models](https://ai.google.dev/gemini-api/docs/models)
- [Google AI Studio](https://aistudio.google.com/)

## Troubleshooting

### Environment Variables Not Loading
Make sure you're running from the correct directory or using absolute paths for `.env` files.

### Port Already in Use
The default port is 8080. You can change it in the agent configuration or kill the process using the port:
```bash
# Find process on port 8080
lsof -ti:8080

# Kill process
kill -9 <PID>
```

## Next Steps

After completing these examples, you're ready to:

- Build custom agents for your specific use cases
- Integrate agents into existing applications
- Create multi-agent systems with specialized agents
- Deploy agents to production with persistent storage

## Contributing

Feel free to submit issues, fork the repository, and create pull requests for any improvements.

## License

This project is for educational purposes.
