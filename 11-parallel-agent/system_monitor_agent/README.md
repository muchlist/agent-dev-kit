# Parallel Agents in ADK (Go Implementation)

This example demonstrates how to implement a Parallel Agent in the Agent Development Kit (ADK) using Go. The main agent in this example, `system_monitor_agent`, uses a hybrid workflow approach that combines parallel execution with sequential processing to create an efficient system monitoring application.

## What are Parallel Agents?

Parallel Agents are workflow agents in ADK that:

1. **Execute Concurrently**: Sub-agents run simultaneously rather than sequentially
2. **Operate Independently**: Each sub-agent works independently without sharing state during execution
3. **Improve Performance**: Dramatically speed up workflows where tasks can be performed in parallel

Use Parallel Agents when you need to execute multiple independent tasks efficiently and time is a critical factor.

## System Monitoring Example (Go)

This example creates a system monitoring application that demonstrates a hybrid workflow pattern:

### Hybrid Workflow Architecture

1. **Parallel Information Gathering**: Three sub-agents run concurrently to collect:
   - CPU usage and statistics
   - Memory utilization
   - Disk space and usage

2. **Sequential Report Synthesis**: After parallel data collection, a synthesizer agent combines all information into a comprehensive report

### Sub-Agents

1. **CPU Info Agent**: Collects and analyzes CPU information
   - Simulates CPU model and architecture analysis
   - Provides usage statistics and performance indicators
   - Identifies potential performance issues

2. **Memory Info Agent**: Gathers memory usage information
   - Analyzes memory utilization and pressure
   - Provides swap usage analysis
   - Identifies memory bottlenecks

3. **Disk Info Agent**: Analyzes disk space and usage
   - Reports on disk capacity and utilization
   - Identifies disks running low on space
   - Provides storage health indicators

4. **System Report Synthesizer**: Combines all gathered information into a comprehensive system health report
   - Creates an executive summary of system health
   - Organizes component-specific information into sections
   - Provides actionable recommendations

## Project Structure

```
11-parallel-agent/
└── system_monitor_agent/          # Main System Monitor Agent package
    ├── main.go                    # Hybrid workflow implementation
    ├── .env.example              # Environment variables template
    └── agents/                    # Sub-agents directory
        ├── cpu_info.go           # CPU information agent
        ├── memory_info.go        # Memory information agent
        ├── disk_info.go          # Disk information agent
        └── synthesizer.go        # Report synthesizing agent
```

## Getting Started

### Setup

1. Navigate to the agent directory:
```bash
cd 11-parallel-agent/system_monitor_agent
```

2. Copy the `.env.example` file to `.env` and add your Google API key:
```bash
cp .env.example .env
# Edit .env with your GOOGLE_API_KEY
```

3. Get Google API key from https://aistudio.google.com/apikey

### Running the Example

```bash
# From the parallel agent directory
go run main.go web api webui

# Or from the root directory using Makefile
make run/11
```

The web UI will launch at http://localhost:8080.

## Example Interactions

### 🎯 **Basic System Health Check:**
```
Check my system health
```

### 📊 **Comprehensive Report Request:**
```
Provide a comprehensive system report with recommendations
```

### 🔍 **Specific Issues Investigation:**
```
Is my system running out of memory or disk space?
```

### 📋 **Detailed Status Report:**
```
Generate a detailed system status report including all components
```

## How It Works

### Hybrid Workflow Architecture

This implementation demonstrates how to combine workflow agent types for optimal performance:

```go
// 1. Create Parallel Agent for concurrent data gathering
parallelInfoGatherer, _ := parallelagent.New(parallelagent.Config{
    AgentConfig: agent.Config{
        Name:        "system_info_gatherer",
        SubAgents:   []agent.Agent{cpuInfoAgent, memoryInfoAgent, diskInfoAgent},
    },
})

// 2. Create Sequential Agent for overall workflow
sequentialAgent, _ := sequentialagent.New(sequentialagent.Config{
    AgentConfig: agent.Config{
        Name:        "system_monitor_agent",
        SubAgents:   []agent.Agent{parallelInfoGatherer, reportSynthesizer},
    },
})
```

### Execution Flow

1. **Parallel Phase**: Three information agents run simultaneously
   - CPU Info Agent → `state["cpu_info_report"]`
   - Memory Info Agent → `state["memory_info_report"]`
   - Disk Info Agent → `state["disk_info_report"]`

2. **Sequential Phase**: Report synthesizer runs after parallel completion
   - Accesses all three reports from state
   - Creates comprehensive health report
   - Stores final result in `state["system_health_report"]`

### Performance Benefits

**Without Parallel (Sequential Only):**
```
CPU Agent (5s) → Memory Agent (5s) → Disk Agent (5s) → Synthesizer (3s)
Total Time: ~18 seconds
```

**With Parallel (This Implementation):**
```
CPU Agent (5s) ↘
Memory Agent (5s) → Parallel (5s) → Synthesizer (3s)
Disk Agent (5s) ↗
Total Time: ~8 seconds (55% faster!)
```

## Key Concepts: Independent Execution

One key aspect of Parallel Agents is that **sub-agents run independently without sharing state during execution**. In this example:

1. Each information gathering agent operates in isolation
2. The results from each agent are collected after parallel execution completes
3. The synthesizer agent then uses these collected results to create the final report

This approach is ideal for scenarios where tasks are completely independent and don't require interaction during execution.

## How Parallel Agents Compare to Other Workflow Agents

| Agent Type | Execution Order | Use Case | Performance |
|------------|----------------|----------|-------------|
| **Sequential Agents** | Strict order | Dependent steps | Baseline |
| **Loop Agents** | Repeated based on conditions | Iterative processing | Variable |
| **Parallel Agents** | Concurrent | Independent tasks | **High Performance** |

## Go vs Python Implementation

This Go version uses the `parallelagent` package from Google's ADK framework:

```go
// Go ParallelAgent creation
parallelAgent, err := parallelagent.New(parallelagent.Config{
    AgentConfig: agent.Config{
        Name:        "parallel_agent",
        Description: "A parallel agent that runs sub-agents",
        SubAgents:   []agent.Agent{agent1, agent2, agent3},
    },
})
```

The Go implementation provides the same concurrent execution capabilities as the Python version while leveraging Go's built-in concurrency support.

## Benefits of This Approach

### Performance Optimization
- **Concurrent Execution**: Multiple agents run simultaneously
- **Reduced Latency**: Independent tasks complete in parallel
- **Scalable Architecture**: Easy to add more parallel agents

### Clean Separation of Concerns
- **Independent Agents**: Each agent focuses on one system component
- **Modular Design**: Easy to modify or extend individual agents
- **Testable Components**: Each agent can be tested independently

### Efficient Data Flow
- **State Management**: Clean data sharing through session state
- **Result Aggregation**: Synthesizer combines all parallel results
- **Minimal Overhead**: No inter-agent communication during execution

## Additional Resources

- [ADK Parallel Agents Documentation](https://google.github.io/adk-docs/agents/workflow-agents/parallel-agents/)
- [Go ADK ParallelAgent Examples](https://github.com/google/adk-go/tree/main/examples/workflowagents/parallel)
- [Full Example: Parallel Web Research](https://google.github.io/adk-docs/agents/workflow-agents/parallel-agents/#full-example-parallel-web-research)