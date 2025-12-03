## What is an ADK Agent?

The ADK agent is a core component that acts as the "thinking" part of your application. It leverages the power of a Large Language Model (LLM) for:
- Reasoning
- Understanding natural language
- Making decisions
- Generating responses
- Interacting with tools

Unlike deterministic workflow agents that follow predefined paths, an ADK agent's behavior is non-deterministic. It uses the LLM to interpret instructions and context, deciding dynamically how to proceed.

## Key Components

### 1. Identity (`Name` and `Description`)
- **Name** (Required): A unique string identifier for your agent
- **Description** (Optional, but recommended): A concise summary of the agent's capabilities

### 2. Model (`Model`)
- Specifies which LLM powers the agent (e.g., "gemini-2.0-flash")
- Affects the agent's capabilities, cost, and performance

### 3. Instructions (`Instruction`)
The most critical parameter for shaping your agent's behavior. It defines:
- Core task or goal
- Personality or persona
- Behavioral constraints
- How to use available tools
- Desired output format

### 4. Tools (`Tools`)
Optional capabilities beyond the LLM's built-in knowledge, allowing the agent to:
- Interact with external systems
- Perform calculations
- Fetch real-time data
- Execute specific actions

## Getting Started

### Prerequisites

1. Go 1.25 or higher installed
2. Google API key for Gemini

### Setup

1. Set up your API key:
   - Copy `.env.example` to `.env` in the greeting_agent folder
   ```bash
   cd greeting_agent
   cp .env.example .env
   ```
   - Add your Google API key to the `GOOGLE_API_KEY` variable in the `.env` file

2. Initialize Go modules (if not already done):
   ```bash
   go mod init github.com/muchlist/agent-dev-kit
   go get google.golang.org/adk
   ```

3. Load environment variables:
   ```bash
   # macOS/Linux:
   export $(cat greeting_agent/.env | xargs)

   # Windows PowerShell:
   Get-Content greeting_agent\.env | ForEach-Object {
       $name, $value = $_.split('=')
       Set-Item -Path env:$name -Value $value
   }
   ```

## Running the Example

### Method 1: Web Interface

Start the web UI to interact with your agent through a browser:

```bash
cd greeting_agent
go run main.go web api webui
```

Then open your browser to `http://localhost:8080`

### Method 2: Command Line Interface

Run the agent directly in your terminal for an interactive CLI session:

```bash
cd greeting_agent
go run main.go run
```

### Method 3: API Server

Start a REST API server for your agent:

```bash
cd greeting_agent
go run main.go api
```

The API will be available at `http://localhost:8080`

### Getting Help

To see all available commands and options:

```bash
cd greeting_agent
go run main.go help
```

## Example Prompts to Try

- "Hello, what's your name?"
- "My name is Alice, can you greet me?"
- "What's a formal way to introduce myself?"

You can exit the CLI conversation by typing `exit` or pressing `Ctrl+C`.

## Differences from Python Version

While the functionality is the same as the Python version, there are some structural differences:

**Python Version:**
- Uses a package structure with `__init__.py` and `agent.py`
- Runs with `adk web` command from parent directory
- Defines `root_agent` variable

**Go Version:**
- Uses a single `main.go` file with a `main()` function
- Runs with `go run main.go [command]` from within the agent directory
- Creates agent directly in the main function
- Uses the ADK launcher pattern for CLI, web, and API modes

Both versions accomplish the same goal: creating a simple greeting agent that asks for the user's name and greets them.

## Learn More

- [Go ADK Documentation](https://google.github.io/adk-docs/get-started/go/)
- [Go ADK GitHub Repository](https://github.com/google/adk-go)
- [ADK Go Package Documentation](https://pkg.go.dev/google.golang.org/adk)
