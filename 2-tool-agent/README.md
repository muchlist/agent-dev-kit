## What is a Tool Agent?

A Tool Agent extends the basic ADK agent by incorporating tools that allow the agent to perform actions beyond just generating text responses. Tools enable agents to interact with external systems, retrieve information, and perform specific functions to accomplish tasks more effectively.

In this example, we demonstrate how to build an agent that can use built-in tools (like Google Search) and custom function tools to enhance its capabilities.

## Key Components

### 1. Built-in Tools

ADK provides several built-in tools that you can use with your agents:

- **Google Search** (`geminitool.GoogleSearch{}`): Allows your agent to search the web for information
- **Code Execution**: Enables your agent to run code snippets
- **Vertex AI Search**: Lets your agent search through your own data

**Important Note**: Currently, for each root agent or single agent, only one built-in tool is supported.

### 2. Custom Function Tools

You can create your own tools by defining Go functions with specific signatures. These custom tools extend your agent's capabilities to perform specific tasks.

#### Creating Custom Function Tools in Go:

```go
// 1. Define argument struct with JSON tags
type getCurrentTimeArgs struct{}

// 2. Define result struct with JSON tags
type getCurrentTimeResults struct {
    CurrentTime string `json:"current_time"`
}

// 3. Implement the function with tool.Context and your structs
func getCurrentTime(ctx tool.Context, input getCurrentTimeArgs) (getCurrentTimeResults, error) {
    currentTime := time.Now().Format("2006-01-02 15:04:05")
    return getCurrentTimeResults{CurrentTime: currentTime}, nil
}

// 4. Register the tool with functiontool.New()
currentTimeTool, err := functiontool.New(
    functiontool.Config{
        Name:        "get_current_time",
        Description: "Get the current time in the format YYYY-MM-DD HH:MM:SS",
    },
    getCurrentTime)
```

#### Best Practices for Custom Function Tools:

- **Type Safety**: Use structs with `json` tags for parameters and return values
- **Required vs Optional**: Use `omitempty` in JSON tags for optional fields
- **Documentation**: Use `jsonschema` tags to document parameters for the LLM
- **Error Handling**: Return errors for invalid inputs or failures
- **Return Type**: Use structs for clear, structured responses

## Limitations

When working with built-in tools in ADK, there are several important limitations to be aware of:

### Single Built-in Tool Restriction

**Currently, for each root agent or single agent, only one built-in tool is supported.**

For example, this approach using two built-in tools within a single agent is **not** currently supported:

```go
tools := []tool.Tool{
    geminitool.GoogleSearch{},
    geminitool.CodeExecution{},  // NOT SUPPORTED with GoogleSearch
}
```

### Built-in Tools vs. Custom Tools

**You cannot mix built-in tools with custom function tools in the same agent.**

For example, this approach is **not** currently supported:

```go
tools := []tool.Tool{
    geminitool.GoogleSearch{},  // Built-in tool
    currentTimeTool,            // Custom tool - NOT SUPPORTED together
}
```

To use both types of tools, you would need to use the multi-agent approach described in later examples.

## Implementation Example

### Understanding the Code

The main.go file defines a tool agent that can use Google Search to find information on the web. The agent is configured with:

1. A name and description
2. The Gemini model to use
3. Instructions that tell the agent how to behave and what tools it can use
4. The tools it can access (in this case, `geminitool.GoogleSearch{}`)

The file also includes a commented-out example of a custom function tool `getCurrentTime()` that you can uncomment to explore custom tool functionality.

## Getting Started

### Prerequisites

1. Go 1.25 or higher installed
2. Google API key for Gemini

### Setup

1. Set up your API key:
   - Copy `.env.example` to `.env` in the tool_agent folder
   ```bash
   cd 2-tool-agent/tool_agent
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

### Method 1: Web Interface

Start the web UI to interact with your agent through a browser:

```bash
cd 2-tool-agent/tool_agent
go run main.go web api webui
```

Then open your browser to `http://localhost:8080`

### Method 2: Command Line Interface

Run the agent directly in your terminal for an interactive CLI session:

```bash
cd 2-tool-agent/tool_agent
go run main.go run
```

### Method 3: API Server

Start a REST API server for your agent:

```bash
cd 2-tool-agent/tool_agent
go run main.go api
```

The API will be available at `http://localhost:8080`

### Method 4: Using Make (from root directory)

If you're at the repository root:

```bash
make run/2
```

### Getting Help

To see all available commands and options:

```bash
cd 2-tool-agent/tool_agent
go run main.go help
```

## Example Prompts to Try

With Google Search enabled:
- "Search for recent news about artificial intelligence"
- "Find information about Google's Agent Development Kit"
- "What are the latest advancements in quantum computing?"

With custom time tool enabled (uncomment the code):
- "What time is it?"
- "Show me the current date and time"

You can exit the CLI conversation by typing `exit` or pressing `Ctrl+C`.

## Switching Between Tools

To switch from Google Search to the custom time tool:

1. Open `main.go`
2. Comment out the Google Search tool lines:
   ```go
   // tools := []tool.Tool{
   //     geminitool.GoogleSearch{},
   // }
   ```
3. Uncomment the custom function tool lines (around line 35-50)
4. Update the agent instruction to reflect the available tool
5. Run the agent again

Remember: You can only use one tool at a time in a single agent!

## Differences from Python Version

| Python | Go |
|--------|-----|
| `from google.adk.tools import google_search` | `import "google.golang.org/adk/tool/geminitool"` |
| `tools=[google_search]` | `tools := []tool.Tool{geminitool.GoogleSearch{}}` |
| Function returns dict | Function returns struct with JSON tags |
| Docstrings for descriptions | Config struct with Name/Description |
| `def get_current_time() -> dict:` | `func getCurrentTime(ctx tool.Context, input Args) (Results, error)` |

Both versions accomplish the same goal but leverage their language's idioms:
- Python uses dynamic typing and dictionaries
- Go uses static typing and structs with compile-time safety

## Learn More

- [ADK Function Tools Documentation](https://google.github.io/adk-docs/tools-custom/function-tools/)
- [ADK Built-in Tools Documentation](https://google.github.io/adk-docs/tools/built-in-tools/)
- [Go ADK GitHub Repository](https://github.com/google/adk-go)
- [ADK Go Package Documentation](https://pkg.go.dev/google.golang.org/adk)
