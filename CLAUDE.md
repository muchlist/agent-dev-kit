# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This is a learning repository for Google's Agent Development Kit (ADK), containing parallel implementations in both **Go** and **Python**. The examples progress from basic agents to complex multi-agent systems with persistent storage.

## Key Architecture Concepts

### Language Differences

**Python Implementation:**
- Uses package structure with `__init__.py` and `agent.py`
- Defines `root_agent` variable that ADK discovers automatically
- Runs with `adk web` command from parent directory
- Session state accessed via `tool_context.state`

**Go Implementation:**
- Uses single `main.go` file with `main()` function
- Creates agent directly in main function
- Runs with `go run main.go [command]` from within the agent directory
- Uses ADK launcher pattern (CLI, web, API modes)
- Session state accessed via `ctx.State()`

### Agent Components

All ADK agents (regardless of language) consist of:
1. **Name** - unique identifier
2. **Model** - LLM to use (e.g., "gemini-2.0-flash")
3. **Instruction** - defines behavior, personality, constraints
4. **Description** - summary of capabilities
5. **Tools** (optional) - functions the agent can call
6. **OutputSchema** (optional) - for structured outputs

### Tool Limitations (Go)

In Go ADK, single agents can only use ONE built-in tool (like GoogleSearch). You CANNOT mix built-in tools with custom function tools in the same agent. To use both types, implement a multi-agent architecture.

### Multi-Agent Architecture

Multi-agent systems use two patterns:
- **Sub-agents**: Directly delegated to by parent agent
- **Agent Tools**: Wrapped agents that appear as tools (for tool-like behavior)

See `python/7-multi-agent/manager/agent.py` for the pattern:
```python
root_agent = Agent(
    sub_agents=[agent1, agent2],  # Direct delegation
    tools=[AgentTool(agent3)]     # Tool-like access
)
```

### Sessions and State

ADK sessions persist conversation state:
- **In-memory**: `session.NewMemoryService()` - development only
- **Database**: `database.NewSessionService()` - production (SQLite/PostgreSQL)

State management:
- Python: `tool_context.state["key"] = value`
- Go: `ctx.State().Set("key", value)` and `ctx.State().Get("key")`

## Common Commands

### Python

```bash
# Setup (run once from root)
python -m venv .venv
source .venv/bin/activate  # macOS/Linux
.venv\Scripts\activate.bat # Windows CMD
pip install -r python/requirements.txt

# Run any example (from root directory)
cd python/<example-dir>/<agent-name>
adk web  # Start web UI on http://localhost:8080
adk run  # CLI mode
adk api  # API server
```

### Go

```bash
# Setup environment (from agent directory)
cd <example-dir>/<agent-name>
cp .env.example .env
# Edit .env with your GOOGLE_API_KEY

# Run examples
go run main.go web api webui  # Web UI + API
go run main.go run            # CLI mode
go run main.go api            # API only
go run main.go help           # Show all commands

# Or use Makefile shortcuts (from root)
make run/1  # Basic agent
make run/2  # Tool agent
make run/3  # LiteLLM agent
make run/4  # Structured outputs
make run/5  # Sessions and state
make run/6  # Persistent storage
```

## Environment Setup

1. Get Google API key from https://aistudio.google.com/apikey
2. Copy `.env.example` to `.env` in the agent directory
3. Add your key: `GOOGLE_API_KEY=your_api_key_here`

## Example Progression

1. **basic-agent**: Simple greeting agent - foundation concepts
2. **tool-agent**: Adds GoogleSearch or custom tools
3. **litellm-agent**: Provider abstraction with LiteLLM
4. **structured-outputs**: Pydantic models (Python) or genai.Schema (Go) for JSON responses
5. **sessions-and-state**: Maintain conversation context
6. **persistent-storage**: Database-backed sessions (SQLite/GORM in Go)
7. **multi-agent**: Orchestrate specialized agents with sub-agents and AgentTools
8. **stateful-multi-agent**: Multi-agent with shared state
9. **callbacks**: Event monitoring and hooks
10. **sequential-agent**: Pipeline workflows
11. **parallel-agent**: Concurrent operations
12. **loop-agent**: Iterative refinement

## Database Usage

The Go examples use GORM with SQLite (file: `my_agent_data.db`). For persistent storage:
- Database is auto-migrated on startup
- Sessions persist across application restarts
- State is serialized as JSON in the database

## Module and Dependencies

- Go module: `github.com/muchlist/agent-dev-kit`
- Main dependencies:
  - `google.golang.org/adk` - Agent Development Kit
  - `google.golang.org/genai` - Gemini API client
  - `gorm.io/gorm` - ORM for database
  - `github.com/joho/godotenv` - Environment variables

## Testing Examples

When modifying agents, test with these interaction patterns:
- Basic info requests ("What can you do?")
- Multi-turn conversations
- Tool invocations
- State persistence (for sessions/storage examples)
- Error cases (invalid inputs, missing data)
