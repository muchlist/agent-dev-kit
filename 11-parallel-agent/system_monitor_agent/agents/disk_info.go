// Package agents implements the sub-agents for the system monitor parallel workflow.
package agents

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
)

// NewDiskInfoAgent creates an agent that analyzes disk space and usage.
// This agent runs in parallel with other system information gatherers.
func NewDiskInfoAgent(ctx context.Context, model model.LLM) (agent.Agent, error) {
	diskInfoAgent, err := llmagent.New(llmagent.Config{
		Name:        "DiskInfoAgent",
		Model:       model,
		Description: "Analyzes disk space, usage, and storage health",
		Instruction: `You are a Disk Information Specialist.

Analyze the system's disk storage and provide a comprehensive report including:
- Total disk capacity
- Used disk space
- Available disk space
- Disk usage percentage
- File system health indicators
- Disks running low on space
- Potential storage issues or warnings
- Recommendations for storage optimization

Since you don't have access to real system metrics, provide a realistic simulation based on the user's request context. Focus on providing useful insights about disk utilization and storage health.

Store your disk analysis in state with the key "disk_info_report".`,
		OutputKey: "disk_info_report",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create disk info agent: %w", err)
	}

	return diskInfoAgent, nil
}