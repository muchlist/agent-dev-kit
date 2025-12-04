# LinkedIn Post Generator Loop Agent (Go Implementation)

This example demonstrates how to implement a Loop Agent in the Agent Development Kit (ADK) using Go. The main agent in this example, `linkedin_post_generator`, uses a hybrid Sequential-Loop pattern to generate and iteratively refine a LinkedIn post until quality requirements are met.

## What are Loop Agents?

Loop Agents are workflow agents in ADK that:

1. **Execute Iteratively**: Run sub-agents repeatedly based on conditions
2. **Automatic Termination**: Stop when exit conditions are met or max iterations reached
3. **Quality-Driven Refinement**: Continue improving content until quality criteria satisfied
4. **Feedback-Based Improvement**: Each iteration uses feedback from previous attempts

Use Loop Agents when you need to iteratively refine content or perform repeated tasks until satisfactory results are achieved.

## LinkedIn Post Generator Example (Go)

This example creates a LinkedIn post generator that demonstrates iterative content refinement:

### Hybrid Workflow Architecture

1. **Initial Post Generation**: Creates first draft of LinkedIn post
2. **Refinement Loop**: Iteratively reviews and refines until quality requirements met

### Loop Control with Exit Tool

A key design pattern is the use of an `exit_loop` tool to control when the loop terminates. The Post Reviewer has two responsibilities:

1. **Quality Evaluation**: Checks if the post meets all requirements
2. **Loop Control**: Calls the exit_loop tool when the post passes all quality checks

When the exit_loop tool is called:
1. It sets `ctx.Actions().Escalate = True`
2. This signals to the LoopAgent that it should stop iterating

### Workflow Components

#### Root Sequential Agent

`LinkedInPostGenerationPipeline` - A SequentialAgent that orchestrates the overall process:
1. First runs the initial post generator
2. Then executes the refinement loop

#### Initial Post Generator

`InitialPostGenerator` - Creates the first draft of the LinkedIn post with no prior context.

#### Refinement Loop

`PostRefinementLoop` - A LoopAgent that executes a two-stage refinement process:
1. First runs the reviewer to evaluate the post and possibly exit the loop
2. Then runs the refiner to improve the post if the loop continues

#### Sub-Agents Inside the Refinement Loop

1. **Post Reviewer** (`PostReviewer`) - Reviews posts for quality and provides feedback or exits the loop if requirements are met
2. **Post Refiner** (`PostRefiner`) - Refines the post based on feedback to improve quality

#### Tools

1. **Character Counter** - Validates post length against requirements (used by the Reviewer)
2. **Exit Loop** - Terminates the loop when all quality criteria are satisfied (used by the Reviewer)

## Project Structure

```
12-loop-agent/
└── linkedin_post_agent/             # Main LinkedIn Post Generator package
    ├── main.go                     # Hybrid workflow implementation
    ├── .env.example               # Environment variables template
    ├── README.md                  # This documentation
    ├── agents/                    # Sub-agents directory
    │   ├── post_generator.go      # Initial post generation agent
    │   ├── post_reviewer.go       # Post quality review agent
    │   └── post_refiner.go        # Post refinement agent
    └── tools/                     # Tools directory
        ├── character_counter.go   # Post length validation tool
        └── exit_loop.go          # Loop termination tool
```

## Getting Started

### Setup

1. Navigate to the agent directory:
```bash
cd 12-loop-agent/linkedin_post_agent
```

2. Copy the `.env.example` file to `.env` and add your Google API key:
```bash
cp .env.example .env
# Edit .env with your GOOGLE_API_KEY
```

3. Get Google API key from https://aistudio.google.com/apikey

### Running the Example

```bash
# From the loop agent directory
go run main.go web api webui

# Or from the root directory using Makefile
make run/12
```

The web UI will launch at http://localhost:8080.

## Example Interactions

### 🎯 **Standard Prompt:**
```
Generate a LinkedIn post about what I've learned from Agent Development Kit tutorial.
```

### 📝 **More Specific Request:**
```
Create a LinkedIn post about the practical applications of Google's ADK that I learned from Agent Development Kit tutorials, focusing on how it helps developers build AI agents.
```

### 🚀 **Technical Focus:**
```
Write a LinkedIn post about the Go implementation of Google's Agent Development Kit, mentioning the different workflow agents (Sequential, Parallel, Loop).
```

## How It Works

### Iterative Refinement Process

The system will:

1. **Generate Initial Post**: Creates first draft based on your request
2. **Review Quality**: Evaluates against quality criteria
3. **Check Requirements**:
   - Mentions @kalseldev ✓
   - Lists 4+ ADK capabilities ✓
   - Has clear call-to-action ✓
   - 1000-1500 characters ✓
   - Professional tone (no emojis/hashtags) ✓
4. **If All Requirements Met**: Calls exit_loop and finishes
5. **If Requirements Not Met**: Provides specific feedback and refines
6. **Continue Loop**: Repeat until max 10 iterations or requirements met

### Quality Requirements

#### **Content Requirements:**
- ✅ Mentions @kalseldev
- ✅ Lists multiple ADK capabilities (at least 4)
- ✅ Has a clear call-to-action
- ✅ Includes practical applications
- ✅ Shows genuine enthusiasm

#### **Style Requirements:**
- ✅ NO emojis
- ✅ NO hashtags
- ✅ Professional tone
- ✅ Conversational style
- ✅ Clear and concise writing

#### **Length Requirements:**
- ✅ 1000-1500 characters
- ✅ Substantial content for LinkedIn

### Loop Termination

The loop terminates in one of two ways:
1. **Quality Success**: When the post meets all requirements (reviewer calls the exit_loop tool)
2. **Max Iterations**: After reaching the maximum number of iterations (10)

## Technical Implementation

### Hybrid Workflow Pattern

```go
// 1. Create Loop Agent for iterative refinement
refinementLoop, err := loopagent.New(loopagent.Config{
    MaxIterations: 10,
    AgentConfig: agent.Config{
        Name:        "PostRefinementLoop",
        SubAgents:   []agent.Agent{postReviewer, postRefiner},
    },
})

// 2. Create Sequential Agent for overall pipeline
sequentialAgent, err := sequentialagent.New(sequentialagent.Config{
    AgentConfig: agent.Config{
        Name:        "LinkedInPostGenerationPipeline",
        SubAgents:   []agent.Agent{initialPostGenerator, refinementLoop},
    },
})
```

### Exit Loop Tool Implementation

```go
// When quality requirements are met
exitLoop := func(ctx tool.Context, args ExitLoopArgs) (ExitLoopResult, error) {
    // Signal to the LoopAgent that we should stop iterating
    ctx.Actions().Escalate = true
    return ExitLoopResult{Success: true}, nil
}
```

### State Management

The workflow uses session state to maintain data between iterations:

- `state["current_post"]` - Current version of the LinkedIn post
- `state["review_feedback"]` - Feedback from the reviewer
- `state["review_status"]` - Pass/fail status from character counter

## Benefits of This Approach

### Quality Assurance
- **Automated Quality Checks**: Systematic evaluation against criteria
- **Iterative Improvement**: Content gets better with each iteration
- **Consistent Standards**: Same quality requirements applied every time

### Controlled Termination
- **Smart Exit**: Stops when quality achieved (not after fixed iterations)
- **Max Iteration Safety**: Prevents infinite loops
- **Tool-Based Control**: Clean separation of logic and flow control

### Modular Design
- **Specialized Agents**: Each agent has a single responsibility
- **Reusable Components**: Tools and agents can be used in other workflows
- **Easy Testing**: Each component can be tested independently

## Performance Considerations

- **Initial Generation**: Fast first draft creation
- **Quality Loop**: 2-4 iterations typically needed for quality posts
- **Tool Overhead**: Character counting and exit checks are lightweight
- **Total Time**: Usually completes within 30-60 seconds

## Extension Ideas

### Additional Quality Checks
- Sentiment analysis
- Grammar and spelling checks
- Professional terminology validation
- LinkedIn engagement optimization

### Alternative Refinement Strategies
- Different reviewer agents for different aspects
- Multiple refinement agents in parallel
- User feedback integration
- A/B testing capabilities

## How Loop Agents Compare to Other Workflow Agents

| Agent Type | Execution Pattern | Use Case | Termination |
|------------|-------------------|----------|-------------|
| **Sequential Agents** | Fixed order | Dependent steps | After final agent |
| **Parallel Agents** | Concurrent | Independent tasks | After all complete |
| **Loop Agents** | Repeated cycles | Iterative refinement | Quality or max iterations |

## Go vs Python Implementation

This Go version uses the `loopagent` package from Google's ADK framework:

```go
// Go LoopAgent creation
loopAgent, err := loopagent.New(loopagent.Config{
    MaxIterations: 10,
    AgentConfig: agent.Config{
        Name:        "loop_agent",
        SubAgents:   []agent.Agent{agent1, agent2},
    },
})
```

The Go implementation provides the same iterative execution capabilities as the Python version while leveraging Go's strong typing and performance characteristics.

## Additional Resources

- [ADK Loop Agents Documentation](https://google.github.io/adk-docs/agents/workflow-agents/loop-agents/)
- [Go ADK LoopAgent Examples](https://github.com/google/adk-go/tree/main/examples/workflowagents/loop)
- [Loop Control Best Practices](https://google.github.io/adk-docs/agents/workflow-agents/loop-agents/#loop-control-best-practices)