# Sequential Agents in ADK

This example demonstrates how to implement a Sequential Agent in the Agent Development Kit (ADK) using Go. The main agent in this example, `lead_qualification_agent`, is a Sequential Agent that executes sub-agents in a predefined order, with each agent's output feeding into the next agent in the sequence.

## What are Sequential Agents?

Sequential Agents are workflow agents in ADK that:

1. **Execute in a Fixed Order**: Sub-agents run one after another in the exact sequence they are specified
2. **Pass Data Between Agents**: Using state management to pass information from one sub-agent to the next
3. **Create Processing Pipelines**: Perfect for scenarios where each step depends on the previous step's output

Use Sequential Agents when you need a deterministic, step-by-step workflow where the execution order matters.

## Lead Qualification Pipeline Example (Go)

In this example, we've created `lead_qualification_agent` as a Sequential Agent that implements a lead qualification pipeline for sales teams. This Sequential Agent orchestrates three specialized sub-agents:

1. **Lead Validator Agent**: Checks if the lead information is complete enough for qualification
   - Validates for required information like contact details and interest
   - Outputs a simple "valid" or "invalid" with a reason

2. **Lead Scorer Agent**: Scores valid leads on a scale of 1-10
   - Analyzes factors like urgency, decision-making authority, budget, and timeline
   - Provides a numeric score with a brief justification

3. **Action Recommender Agent**: Suggests next steps based on the validation and score
   - For invalid leads: Recommends what information to gather
   - For low-scoring leads (1-3): Suggests nurturing actions
   - For medium-scoring leads (4-7): Suggests qualifying actions
   - For high-scoring leads (8-10): Suggests sales actions

### How It Works

The `lead_qualification_agent` Sequential Agent orchestrates this process by:

1. Running the Validator first to determine if the lead is complete
2. Running the Scorer next (which can access validation results via state)
3. Running the Recommender last (which can access both validation and scoring results)

The output of each sub-agent is stored in the session state using the `output_key` parameter:
- `validation_status`
- `lead_score`
- `action_recommendation`

## Project Structure

```
10-sequential-agent/
└── lead_qualification_agent/       # Main Sequential Agent package
    ├── main.go                     # Sequential Agent definition and main function
    ├── .env.example                # Environment variables example
    └── agents/                     # Sub-agents directory
        ├── validator.go            # Lead validation agent
        ├── scorer.go               # Lead scoring agent
        └── recommender.go          # Action recommendation agent
```

## Getting Started

### Setup

1. Navigate to the agent directory:
```bash
cd 10-sequential-agent/lead_qualification_agent
```

2. Copy the `.env.example` file to `.env` and add your Google API key:
```bash
cp .env.example .env
# Edit .env with your GOOGLE_API_KEY
```

3. Get Google API key from https://aistudio.google.com/apikey

### Running the Example

```bash
# From the sequential agent directory
go run main.go web api webui

# Or from the root directory using Makefile
make run/10
```

The web UI will launch at http://localhost:8080. Select "LeadQualificationPipeline" from the dropdown menu.

## Example Initial Chat Messages

### 🎯 Qualified Lead Example (Copy and paste this to start):

```
I need to qualify this sales lead:

Name: Sarah Johnson
Email: sarah.j@techinnovate.com
Phone: 555-123-4567
Company: Tech Innovate Solutions
Position: CTO
Interest: Looking for an AI solution to automate customer support
Budget: $50K-100K available
Timeline: Hoping to implement within next quarter
Notes: Currently using a competitor's product but unhappy with performance
```

### 📝 Unqualified Lead Example:

```
Please qualify this lead:

Name: John Doe
Email: john@gmail.com
Interest: Something with AI maybe
Notes: Met at conference, seemed interested but was vague about needs
```

### 💬 Conversational Format:

```
Hi, I have a new lead for you. Sarah Johnson is the CTO at Tech Innovate Solutions. You can reach her at sarah.j@techinnovate.com or 555-123-4567. She's looking for an AI solution to automate customer support, has a budget of $50K-100K, and wants to implement within the next quarter. Currently using a competitor's product but is unhappy with their performance.
```

### 📊 Medium-Quality Lead Example:

```
Lead Information:
Name: Mike Chen
Email: mike.chen@startup.com
Company: Tech Startup Inc
Interest: ML tools for data analysis
Budget: Limited but growing
Timeline: Maybe next year
Notes: Technical person but decision-making authority unclear
```

## What to Expect

When you provide lead information, the Sequential Agent will:

1. **Step 1 - Validation**: Check if the lead has sufficient information
2. **Step 2 - Scoring**: Assign a qualification score (1-10) with justification
3. **Step 3 - Recommendation**: Provide actionable next steps

**Example Output for Qualified Lead:**
```
✅ VALIDATION: valid

📊 SCORING: 8: Decision maker with clear budget and immediate need

🎯 RECOMMENDATION:
- Schedule a demo within 48 hours
- Prepare a technical proposal focused on customer support automation
- Address competitor comparison in the demo
- Follow up with detailed pricing options
```

## How Sequential Agents Compare to Other Workflow Agents

ADK offers different types of workflow agents for different needs:

- **Sequential Agents**: For strict, ordered execution (like this example)
- **Loop Agents**: For repeated execution of sub-agents based on conditions
- **Parallel Agents**: For concurrent execution of independent sub-agents

## Go vs Python Implementation

This Go version uses the `sequentialagent` package from Google's ADK framework:

```go
// Go SequentialAgent creation
sequentialAgent, err := sequentialagent.New(sequentialagent.Config{
    AgentConfig: agent.Config{
        Name:        "LeadQualificationPipeline",
        Description: "A sequential pipeline that validates, scores, and recommends actions for sales leads",
        SubAgents:   []agent.Agent{validator, scorer, recommender},
    },
})
```

The Go implementation provides the same sequential execution capabilities as the Python version while leveraging Go's strong typing and performance characteristics.

## Additional Resources

- [ADK Sequential Agents Documentation](https://google.github.io/adk-docs/agents/workflow-agents/sequential-agents/)
- [Go ADK SequentialAgent Examples](https://github.com/google/adk-go/tree/main/examples/workflowagents)
- [Full Code Development Pipeline Example](https://google.github.io/adk-docs/agents/workflow-agents/sequential-agents/#full-example-code-development-pipeline)