# Multi-Agent Systems in ADK (Go)

This example demonstrates how to create a multi-agent system in ADK using Go with **modular organization**, where specialized agents collaborate to handle complex tasks, each focusing on their area of expertise.

## Project Structure

This Go implementation uses a **modular approach** with separate files for each agent and shared tools:

```
7-multi-agent/
└── manager_agent/
    ├── main.go                 # Entry point and manager agent
    ├── agents/                 # Specialized agent modules
    │   ├── stock_analyst.go    # Stock market analysis agent
    │   ├── funny_nerd.go       # Nerdy jokes agent
    │   └── news_analyst.go     # News search agent
    ├── tools/                  # Shared utility tools
    │   └── time.go             # Current time tool
    ├── .env.example
    └── .env
```

## Key Features of This Implementation

### 1. **Modular Organization**
Unlike the typical Go ADK pattern of putting everything in `main.go`, this example organizes code into separate packages:

- **`agents/` package**: Each agent has its own file with tools and creation logic
- **`tools/` package**: Shared tools used by multiple agents
- **`main.go`**: Entry point that assembles everything

### 2. **Multi-Agent Architecture**
The system uses both delegation patterns supported by ADK:

- **Sub-Agents** (Direct delegation): `stock_analyst` and `funny_nerd`
- **Agent Tools** (Tool-like usage): `news_analyst` wrapped as a tool

### 3. **Built-in Tool Handling**
Demonstrates the Go ADK limitation workaround:
- Single agents can only use ONE built-in tool
- `news_analyst` uses `GoogleSearch` (built-in) and is wrapped as `AgentTool`
- Manager can then combine it with custom tools

## What is a Multi-Agent System?

A Multi-Agent System allows multiple specialized agents to work together:
- Each agent focuses on a specific domain
- Agents collaborate through delegation and communication
- Solves complex problems that would be difficult for a single agent

## Multi-Agent Architecture Options

### 1. Sub-Agent Delegation Model

Using the `SubAgents` parameter:

```go
managerAgent, err := llmagent.New(llmagent.Config{
    Name:        "manager",
    Model:       model,
    SubAgents:   []agent.Agent{stockAnalyst, funnyNerd},
})
```

**Characteristics:**
- Complete delegation - sub-agent takes full control
- Sub-agent handles the entire response
- Manager acts as a "router"

### 2. Agent-as-a-Tool Model

Using `agenttool.New()`:

```go
newsAnalystTool := agenttool.New(newsAnalyst, &agenttool.Config{})

managerAgent, err := llmagent.New(llmagent.Config{
    Name:  "manager",
    Tools: []tool.Tool{newsAnalystTool, getCurrentTimeTool},
})
```

**Characteristics:**
- Sub-agent returns results to manager
- Manager maintains control
- Can combine results from multiple agents

## Important: Go ADK Tool Limitations

### The One Built-in Tool Rule

**In Go ADK, a single agent can only use ONE built-in tool** (like `GoogleSearch`). You cannot mix built-in tools with custom function tools.

❌ **This will NOT work:**
```go
agent, err := llmagent.New(llmagent.Config{
    Name: "hybrid_agent",
    Tools: []tool.Tool{
        geminitool.GoogleSearch{},  // Built-in
        myCustomTool,                // Custom - ERROR!
    },
})
```

✅ **Workaround using Multi-Agent:**
```go
// Create agent with built-in tool
searchAgent, err := llmagent.New(llmagent.Config{
    Name:  "search_agent",
    Tools: []tool.Tool{geminitool.GoogleSearch{}},
})

// Wrap as AgentTool
searchTool := agenttool.New(searchAgent, &agenttool.Config{})

// Manager can now use both
manager, err := llmagent.New(llmagent.Config{
    Name: "manager",
    Tools: []tool.Tool{
        searchTool,      // Wrapped built-in
        myCustomTool,    // Custom - OK!
    },
})
```

This is exactly how `news_analyst` is handled in this example.

## Our Multi-Agent System

This example implements three specialized agents coordinated by a manager:

### 1. **Stock Analyst** (Sub-agent)
- **File**: `agents/stock_analyst.go`
- **Tool**: `get_stock_price` - retrieves mock stock prices
- **Purpose**: Provides stock market information
- **Available tickers**: GOOG, GOOGL, TSLA, META, AAPL, MSFT, AMZN

### 2. **Funny Nerd** (Sub-agent)
- **File**: `agents/funny_nerd.go`
- **Tool**: `get_nerd_joke` - returns topic-specific jokes
- **Purpose**: Tells nerdy jokes about technical topics
- **Features**: Uses state to store last joke topic
- **Topics**: python, javascript, java, go, programming, math, physics, chemistry, biology, computer, database

### 3. **News Analyst** (Agent Tool)
- **File**: `agents/news_analyst.go`
- **Tool**: `GoogleSearch` (built-in)
- **Purpose**: Searches and summarizes current news
- **Note**: Wrapped as AgentTool due to built-in tool limitation

### 4. **Manager Agent**
- **File**: `main.go`
- **Sub-agents**: stock_analyst, funny_nerd
- **Tools**: news_analyst (as AgentTool), get_current_time
- **Purpose**: Routes queries to appropriate specialists

## Getting Started

### Prerequisites

1. Go 1.21 or later
2. Google API key from https://aistudio.google.com/apikey

### Setup

1. Navigate to the manager_agent directory:
```bash
cd 7-multi-agent/manager_agent
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
make run/7
```

### Direct Execution

#### Web UI Mode (with API backend)
```bash
go run 7-multi-agent/manager_agent/main.go web api webui
```
Then open `http://localhost:8080` in your browser.

#### CLI Console Mode
```bash
go run 7-multi-agent/manager_agent/main.go console
```

#### API Server Only
```bash
go run 7-multi-agent/manager_agent/main.go web api
```

## Example Prompts to Try

### Test Stock Analyst
- "What's the current price of GOOG?"
- "Can you check the prices for TSLA and META?"
- "Show me Apple and Microsoft stock prices"

### Test Funny Nerd
- "Tell me a joke about Python"
- "I want to hear something funny about programming"
- "Do you know any physics jokes?"
- "Tell me a Go programming joke"

### Test News Analyst
- "What's the latest tech news?"
- "Search for news about artificial intelligence"
- "Find recent news about Google"

### Test Manager's Tools
- "What time is it?"
- "What's the current date and time?"

### Test Multi-Agent Routing
- "Tell me a joke about JavaScript and then check MSFT stock price"
- "What's the tech news today and what time is it?"

## Code Organization

### Agent Module Pattern

Each agent follows this pattern in `agents/` package:

```go
// 1. Tool argument and result structures
type getToolArgs struct { ... }
type getToolResults struct { ... }

// 2. Tool implementation function
func toolFunction(ctx tool.Context, input Args) (Results, error) {
    // Tool logic
    // Access state: ctx.State().Set/Get
}

// 3. Agent creation function
func NewAgentName(ctx context.Context, mdl model.LLM) (agent.Agent, error) {
    // Create tools
    // Create and return agent
}
```

### Shared Tools Pattern

Shared tools in `tools/` package:

```go
// NewToolName creates and returns the tool
func NewToolName() (tool.Tool, error) {
    return functiontool.New(
        functiontool.Config{...},
        toolImplementation,
    )
}
```

### Main Entry Point

`main.go` assembles everything:

```go
func main() {
    // 1. Create model
    // 2. Create specialized agents using agents.NewXxx()
    // 3. Create manager with createManagerAgent()
    // 4. Launch with launcher
}
```

## Extending the System

### Adding a New Sub-Agent

1. **Create agent file** in `agents/`:
```go
// agents/my_agent.go
package agents

func NewMyAgent(ctx context.Context, mdl model.LLM) (agent.Agent, error) {
    // Define tools and create agent
}
```

2. **Import and use in** `main.go`:
```go
myAgent, err := agents.NewMyAgent(ctx, model)
// Add to manager's SubAgents or wrap as AgentTool
```

### Adding a Shared Tool

1. **Create tool file** in `tools/`:
```go
// tools/mytool.go
package tools

func NewMyTool() (tool.Tool, error) {
    // Create and return tool
}
```

2. **Use in any agent or manager**:
```go
import "github.com/muchlist/agent-dev-kit/7-multi-agent/manager_agent/tools"

myTool, err := tools.NewMyTool()
```

### Using Real Stock API

Replace mock data in `agents/stock_analyst.go`:

```go
// Example with Alpha Vantage
apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
url := fmt.Sprintf(
    "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
    input.Ticker, apiKey)

resp, err := http.Get(url)
// Parse JSON response
```

## Modular vs Monolithic Comparison

| Aspect | This Example (Modular) | Typical Go ADK (Monolithic) |
|--------|----------------------|--------------------------|
| Organization | Separate files per agent | Single main.go |
| Maintainability | High - easy to find/modify | Low - everything mixed |
| Reusability | High - import packages | Low - copy/paste |
| Testing | Easy - test per module | Hard - must test all |
| Collaboration | Easy - parallel development | Hard - merge conflicts |
| Learning | Clear separation of concerns | Can be overwhelming |

## Comparison with Python Version

| Feature | Python | Go (This Example) |
|---------|--------|------------------|
| Structure | Multiple packages/files | Modular (agents/ + tools/) |
| Agent Import | `from .sub_agents...` | `agents.NewXxx()` |
| Tool Context | `tool_context.state["key"]` | `ctx.State().Set/Get()` |
| Running | `adk web` from parent | `go run ...` or `make run/7` |
| Built-in Tools | Can mix with custom | One per agent, use AgentTool |
| Type Safety | Runtime (duck typing) | Compile-time (static) |

## Troubleshooting

### Common Issues

1. **"no required module provides package"**
   - Make sure you're running from repository root
   - The parent `go.mod` must exist
   - Solution: `go run 7-multi-agent/manager_agent/main.go ...`

2. **"Failed to create google search tool"**
   - Check `GOOGLE_API_KEY` in `.env`
   - Verify API key has necessary permissions

3. **"undefined: gemini.Model"**
   - Use `model.LLM` interface, not `*gemini.Model`
   - Import: `"google.golang.org/adk/model"`

4. **Agent not delegating correctly**
   - Review manager's `Instruction` field
   - Check agent `Description` fields
   - Ensure instructions clearly define delegation rules

### Debug Output

The example includes debug prints for tool calls:

```
--- Tool: get_stock_price called for GOOG ---
--- Tool: get_nerd_joke called for topic: python ---
--- Tool: get_current_time called ---
```

These help trace which agents and tools are being invoked.

## Additional Resources

- [ADK Multi-Agent Documentation](https://google.github.io/adk-docs/agents/multi-agent-systems/)
- [Agent Tools Documentation](https://google.github.io/adk-docs/tools/function-tools/#3-agent-as-a-tool)
- [Go ADK API Reference](https://pkg.go.dev/google.golang.org/adk)
- [Gemini API Documentation](https://ai.google.dev/docs)

## Next Steps

After mastering multi-agent systems, explore:
- **Example 8 - Stateful Multi-Agent**: Add persistent database storage
- **Example 9 - Callbacks**: Monitor agent events and add hooks
- **Example 10 - Sequential Agent**: Pipeline workflows
- **Example 11 - Parallel Agent**: Concurrent operations
