// Package agents implements the sub-agents for the system monitor parallel workflow.
package agents

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
)

// NewMemoryInfoAgent creates an agent that gathers memory usage information.
// This agent runs in parallel with other system information gatherers.
func NewMemoryInfoAgent(ctx context.Context, model model.LLM) (agent.Agent, error) {
	memoryInfoAgent, err := llmagent.New(llmagent.Config{
		Name:        "MemoryInfoAgent",
		Model:       model,
		Description: "Gathers memory usage information and analyzes memory pressure",
		Instruction: `You are a Memory Information Specialist.

Analyze the system's memory and provide a detailed report including:
- Total RAM capacity
- Current memory usage
- Available memory
- Memory usage percentage
- Swap usage (if applicable)
- Memory pressure indicators
- Potential memory issues or bottlenecks
- Recommendations for memory optimization

Since you don't have access to real system metrics, provide a realistic simulation based on the user's request context. Focus on providing useful insights about memory utilization and health.

Store your memory analysis in state with the key "memory_info_report".`,
		OutputKey: "memory_info_report",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create memory info agent: %w", err)
	}

	return memoryInfoAgent, nil
}