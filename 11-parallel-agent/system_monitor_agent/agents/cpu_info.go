// Package agents implements the sub-agents for the system monitor parallel workflow.
package agents

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
)

// NewCPUInfoAgent creates an agent that collects and analyzes CPU information.
// This agent runs in parallel with other system information gatherers.
func NewCPUInfoAgent(ctx context.Context, model model.LLM) (agent.Agent, error) {
	cpuInfoAgent, err := llmagent.New(llmagent.Config{
		Name:        "CPUInfoAgent",
		Model:       model,
		Description: "Collects and analyzes CPU information and performance metrics",
		Instruction: `You are a CPU Information Specialist.

Analyze the system's CPU and provide a comprehensive report including:
- CPU model and architecture
- Number of cores/threads
- Current usage statistics
- Performance indicators
- Any potential issues (high usage, overheating, etc.)
- Recommendations for optimization

Since you don't have access to real system metrics, provide a realistic simulation based on the user's request context. Focus on providing useful insights about CPU performance and health.

Store your CPU analysis in state with the key "cpu_info_report".`,
		OutputKey: "cpu_info_report",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create CPU info agent: %w", err)
	}

	return cpuInfoAgent, nil
}