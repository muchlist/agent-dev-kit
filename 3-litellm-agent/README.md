# Multi-Model Agent Example (Go)

This example demonstrates using custom function tools with ADK agents. The Python version uses LiteLLM to access non-Google models like OpenAI through OpenRouter, but the Go version currently uses Gemini due to Go ADK's current model support.

## Important: Model Provider Differences

### Python Version (Reference)
The Python version in `python/3-litellm-agent` demonstrates:
- ✓ LiteLLM integration for accessing 100+ different LLM providers
- ✓ Using OpenAI GPT-4 through OpenRouter
- ✓ Anthropic Claude through OpenRouter
- ✓ Easy switching between model providers

### Go Version (This Example)
**Current Limitation**: Go ADK does not yet have native LiteLLM integration like Python ADK.

**Model Support Status**:
- ✓ **Gemini** - Full native support via `google.golang.org/adk/model/gemini`
- ✗ **OpenAI** - No native package available yet
- ✗ **Anthropic** - No native package available yet
- ✗ **LiteLLM** - No integration like Python ADK has

While some sources claim Go ADK is "model-agnostic" and supports multiple providers, concrete implementation packages for OpenAI/Anthropic are not currently available in the official Go ADK package (as of December 2025).

## What This Example Demonstrates

Despite the model provider limitation, this example still demonstrates the core concepts:

1. **Custom Function Tools** - Creating tools with struct-based arguments and returns
2. **Tool Integration** - Registering and using tools with an agent
3. **Agent Configuration** - Setting up an agent with specific instructions

The functionality (telling dad jokes via a custom tool) is identical to the Python version, just powered by Gemini instead of OpenAI/OpenRouter.

## Future Outlook

The Go ADK team may add support for:
- Direct OpenAI integration (similar to Python's approach)
- Anthropic Claude support
- A Go-native LiteLLM-like abstraction layer
- Community-contributed model adapters

For now, developers using Go ADK should use Gemini models, which offer excellent performance and capabilities.

## Getting Started

### Prerequisites

1. Go 1.25 or higher installed
2. Google API key for Gemini

### Setup

1. Set up your API key:
   - Copy `.env.example` to `.env` in the dad_joke_agent folder
   ```bash
   cd 3-litellm-agent/dad_joke_agent
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
cd 3-litellm-agent/dad_joke_agent
go run main.go web api webui
```

Then open your browser to `http://localhost:8080`

### Method 2: Command Line Interface

Run the agent directly in your terminal for an interactive CLI session:

```bash
cd 3-litellm-agent/dad_joke_agent
go run main.go run
```

### Method 3: API Server

Start a REST API server for your agent:

```bash
cd 3-litellm-agent/dad_joke_agent
go run main.go api
```

The API will be available at `http://localhost:8080`

### Method 4: Using Make (from root directory)

If you're at the repository root:

```bash
make run/3
```

### Getting Help

To see all available commands and options:

```bash
cd 3-litellm-agent/dad_joke_agent
go run main.go help
```

## Example Prompts to Try

- "Tell me a dad joke"
- "Can you tell me another joke?"
- "Give me your best dad joke"

You can exit the CLI conversation by typing `exit` or pressing `Ctrl+C`.

## Understanding the Code

### Custom Function Tool Structure

The Go version demonstrates the proper way to create custom tools:

```go
// 1. Define argument struct (empty in this case)
type getDadJokeArgs struct{}

// 2. Define result struct with JSON tags
type getDadJokeResults struct {
    Joke string `json:"joke"`
}

// 3. Implement the function
func getDadJoke(ctx tool.Context, input getDadJokeArgs) (getDadJokeResults, error) {
    jokes := []string{
        "Why did the chicken cross the road? To get to the other side!",
        "What do you call a belt made of watches? A waist of time.",
        // ... more jokes
    }
    randomJoke := jokes[rand.Intn(len(jokes))]
    return getDadJokeResults{Joke: randomJoke}, nil
}

// 4. Register the tool
dadJokeTool, err := functiontool.New(
    functiontool.Config{
        Name:        "get_dad_joke",
        Description: "Returns a random dad joke",
    },
    getDadJoke)
```

## Comparison: Python vs Go

| Aspect | Python | Go |
|--------|--------|-----|
| **Model Provider** | OpenAI via LiteLLM/OpenRouter | Gemini (native) |
| **API Key** | `OPENROUTER_API_KEY` | `GOOGLE_API_KEY` |
| **Model Import** | `from google.adk.models.lite_llm import LiteLlm` | `import "google.golang.org/adk/model/gemini"` |
| **Model Creation** | `LiteLlm(model="openrouter/openai/gpt-4.1", api_key=...)` | `gemini.NewModel(ctx, "gemini-2.0-flash", ...)` |
| **Function Tool Args** | Function parameters | Struct with tags |
| **Function Tool Return** | Python dict or primitive | Struct with JSON tags + error |
| **Multi-Provider** | ✓ Easy via LiteLLM | ✗ Gemini only currently |

## Implementing Custom Model Providers (Advanced)

If you need to use non-Gemini models in Go ADK, you can implement the `model.LLM` interface:

```go
type LLM interface {
    Name() string
    GenerateContent(ctx context.Context, req *LLMRequest, stream bool) iter.Seq2[*LLMResponse, error]
}
```

However, this requires:
- Implementing the complete interface
- Handling request/response translation
- Managing tool call serialization
- Supporting streaming responses

This is significantly more complex than Python's LiteLLM approach and is only recommended for advanced users with specific requirements.

## Learn More

**Go ADK Documentation:**
- [Go ADK Getting Started](https://google.github.io/adk-docs/get-started/go/)
- [ADK Models & Authentication](https://google.github.io/adk-docs/agents/models/)
- [Go ADK Package Documentation](https://pkg.go.dev/google.golang.org/adk)
- [Go ADK GitHub Repository](https://github.com/google/adk-go)

**Python ADK with LiteLLM (Reference):**
- [Google ADK LiteLLM Integration](https://docs.litellm.ai/docs/tutorials/google_adk)
- [LiteLLM Documentation](https://docs.litellm.ai/docs/)
- [OpenRouter Documentation](https://openrouter.ai/docs)

**Community Resources:**
- [Google ADK Masterclass Part 3: Using Different Models](https://saptak.in/writing/2025/05/10/google-adk-masterclass-part3)
- [Building AI Agents with ADK-Go](https://byteiota.com/google-adk-go-tutorial-build-ai-agents-in-go-2025/)
