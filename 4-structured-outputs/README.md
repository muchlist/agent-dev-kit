# Structured Outputs in ADK (Go)

This example demonstrates how to implement structured outputs in the Agent Development Kit (ADK) using Go. The main agent in this example, `email_agent`, uses the `OutputSchema` parameter to ensure its responses conform to a specific structured format.

## What are Structured Outputs?

ADK allows you to define structured data formats for agent outputs using schema definitions:

1. **Controlled Output Format**: Using `OutputSchema` ensures the LLM produces responses in a consistent JSON structure
2. **Data Validation**: The schema validates that all required fields are present and correctly formatted
3. **Improved Downstream Processing**: Structured outputs are easier to handle in downstream applications or by other agents
4. **Agent Coordination**: Using `OutputKey` stores the result in session state for use by other agents

Use structured outputs when you need guaranteed format consistency for integration with other systems or agents.

## Email Generator Example

In this example, we've created an email generator agent that produces structured output with:

1. **Email Subject**: A concise, relevant subject line
2. **Email Body**: Well-formatted email content with greeting, paragraphs, and signature

The agent uses a `genai.Schema` definition to ensure every response follows the same format.

### Output Schema Definition

The Go version uses `genai.Schema` from the `google.golang.org/genai` package:

```go
emailSchema := &genai.Schema{
    Type: "OBJECT",
    Properties: map[string]*genai.Schema{
        "subject": {
            Type:        "STRING",
            Description: "The subject line of the email. Should be concise and descriptive.",
        },
        "body": {
            Type:        "STRING",
            Description: "The main content of the email. Should be well-formatted with proper greeting, paragraphs, and signature.",
        },
    },
    Required: []string{"subject", "body"},
}
```

### How It Works

1. The user provides a description of the email they need
2. The LLM agent processes this request and generates both a subject and body
3. The agent formats its response as a JSON object matching the schema
4. ADK validates the response against the schema before returning it
5. The structured output is stored in the session state under the specified `OutputKey`

### Agent Configuration

```go
a, err := llmagent.New(llmagent.Config{
    Name:         "email_agent",
    Model:        model,
    Description:  "Generates professional emails with structured subject and body",
    Instruction:  `You are an Email Generation Assistant...`,
    OutputSchema: emailSchema,  // Defines the required output structure
    OutputKey:    "email",       // Stores result in session state["email"]
})
```

## Important Limitations

When using `OutputSchema`:

1. **No Tool Usage**: Agents with an output schema cannot use tools during their execution
2. **Direct JSON Response**: The LLM must produce a JSON response matching the schema as its final output
3. **Clear Instructions**: The agent's instructions must explicitly guide the LLM to produce properly formatted JSON

This limitation is documented in the Go ADK source code:

```go
// NOTE: when this is set, agent can only reply and cannot use any tools,
// such as function tools, RAGs, agent transfer, etc.
```

## Getting Started

### Prerequisites

1. Go 1.25 or higher installed
2. Google API key for Gemini

### Setup

1. Set up your API key:
   - Copy `.env.example` to `.env` in the email_agent folder
   ```bash
   cd 4-structured-outputs/email_agent
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
cd 4-structured-outputs/email_agent
go run main.go web api webui
```

Then open your browser to `http://localhost:8080`

### Method 2: Command Line Interface

Run the agent directly in your terminal for an interactive CLI session:

```bash
cd 4-structured-outputs/email_agent
go run main.go run
```

### Method 3: API Server

Start a REST API server for your agent:

```bash
cd 4-structured-outputs/email_agent
go run main.go api
```

The API will be available at `http://localhost:8080`

### Method 4: Using Make (from root directory)

If you're at the repository root:

```bash
make run/4
```

### Getting Help

To see all available commands and options:

```bash
cd 4-structured-outputs/email_agent
go run main.go help
```

## Example Prompts to Try

```
Write a professional email to my team about the upcoming project deadline that has been extended by two weeks.
```

```
Draft an email to a client explaining that we need additional information before we can proceed with their order.
```

```
Create an email to schedule a meeting with the marketing department to discuss the new product launch strategy.
```

You can exit the CLI conversation by typing `exit` or pressing `Ctrl+C`.

## Understanding the Code

### genai.Schema Structure

The `genai.Schema` type supports OpenAPI 3.0 schema definitions with the following key fields:

```go
type Schema struct {
    Type        string                  // "OBJECT", "STRING", "NUMBER", "INTEGER", "BOOLEAN", "ARRAY", "NULL"
    Properties  map[string]*Schema      // For OBJECT types - defines nested fields
    Required    []string                // Required field names
    Description string                  // Field description
    Items       *Schema                 // For ARRAY types - defines array item schema
    Format      string                  // Additional format details (e.g., "date-time")
    Pattern     string                  // Regex pattern for STRING types
    // ... and other OpenAPI schema properties
}
```

### Example: Nested Object Schema

For more complex structures:

```go
complexSchema := &genai.Schema{
    Type: "OBJECT",
    Properties: map[string]*genai.Schema{
        "subject": {Type: "STRING"},
        "body": {Type: "STRING"},
        "attachments": {
            Type: "ARRAY",
            Items: &genai.Schema{
                Type: "OBJECT",
                Properties: map[string]*genai.Schema{
                    "filename": {Type: "STRING"},
                    "type": {Type: "STRING"},
                },
                Required: []string{"filename", "type"},
            },
        },
    },
    Required: []string{"subject", "body"},
}
```

## Key Concepts: Structured Data Exchange

Structured outputs are part of ADK's broader support for structured data exchange, which includes:

1. **InputSchema**: Define expected input format (not used in this example)
2. **OutputSchema**: Define required output format (used in this example)
3. **OutputKey**: Store the result in session state for use by other agents (used in this example)

### OutputKey Usage Pattern

When an agent has `OutputKey` set, its response is automatically stored in the session state:

```go
// Agent 1: Generates code
codeWriter, err := llmagent.New(llmagent.Config{
    Name:      "code_writer",
    OutputKey: "generated_code",  // Stores output in state["generated_code"]
    // ...
})

// Agent 2: Reviews the generated code
codeReviewer, err := llmagent.New(llmagent.Config{
    Name: "code_reviewer",
    Instruction: `Review this code:
{generated_code}  // References state["generated_code"] from Agent 1
...`,
    OutputKey: "review_comments",
})
```

This pattern enables reliable data passing between agents in multi-agent workflows.

## Comparison: Python vs Go

| Aspect | Python | Go |
|--------|--------|-----|
| **Schema Definition** | Pydantic `BaseModel` | `genai.Schema` struct |
| **Import** | `from pydantic import BaseModel, Field` | `import "google.golang.org/genai"` |
| **Field Types** | Python type hints (`str`, `int`, etc.) | Schema types (`"STRING"`, `"INTEGER"`, etc.) |
| **Schema Config** | `output_schema=EmailContent` | `OutputSchema: emailSchema` |
| **Output Key** | `output_key="email"` | `OutputKey: "email"` |
| **Validation** | Pydantic validation | JSON schema validation |
| **Tool Limitation** | ✓ No tools allowed | ✓ No tools allowed |

### Python vs Go Examples Side by Side

**Python:**
```python
from pydantic import BaseModel, Field

class EmailContent(BaseModel):
    subject: str = Field(description="...")
    body: str = Field(description="...")

root_agent = LlmAgent(
    name="email_agent",
    output_schema=EmailContent,
    output_key="email",
)
```

**Go:**
```go
import "google.golang.org/genai"

emailSchema := &genai.Schema{
    Type: "OBJECT",
    Properties: map[string]*genai.Schema{
        "subject": {Type: "STRING", Description: "..."},
        "body": {Type: "STRING", Description: "..."},
    },
    Required: []string{"subject", "body"},
}

a, err := llmagent.New(llmagent.Config{
    Name:         "email_agent",
    OutputSchema: emailSchema,
    OutputKey:    "email",
})
```

## Advanced: Schema Types Reference

Common schema types you can use:

- `"STRING"` - Text data
- `"INTEGER"` - Whole numbers
- `"NUMBER"` - Floating-point numbers
- `"BOOLEAN"` - true/false values
- `"OBJECT"` - Nested structures (use with `Properties`)
- `"ARRAY"` - Lists (use with `Items`)
- `"NULL"` - Null values

## Real-World Use Cases

1. **Data Extraction**: Extract structured information from unstructured text
2. **Form Generation**: Generate form data in a specific format
3. **API Integration**: Produce outputs that match external API requirements
4. **Multi-Agent Workflows**: Pass structured data between agents in a pipeline
5. **Validation**: Ensure LLM responses meet strict format requirements

## Learn More

**Go ADK Documentation:**
- [ADK Structured Data Documentation](https://google.github.io/adk-docs/agents/llm-agents/#structuring-data-input_schema-output_schema-output_key)
- [Go ADK GitHub Repository](https://github.com/google/adk-go)
- [Go ADK Package Documentation](https://pkg.go.dev/google.golang.org/adk)
- [Go genai Package](https://pkg.go.dev/google.golang.org/genai)

**Reference Implementation:**
- [Sequential Code Pipeline Example](https://github.com/google/adk-go/tree/main/examples/workflowagents/sequentialCode) - demonstrates OutputKey usage

**Python ADK (Reference):**
- [Python ADK Structured Outputs](https://google.github.io/adk-docs/agents/llm-agents/)
- [Pydantic Documentation](https://docs.pydantic.dev/latest/)
