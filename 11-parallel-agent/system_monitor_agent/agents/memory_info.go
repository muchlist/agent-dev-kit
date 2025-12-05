// Package agents implements the sub-agents for the system monitor parallel workflow.
package agents

import (
	"context"
	"fmt"

	"github.com/muchlist/agent-dev-kit/11-parallel-agent/system_monitor_agent/tools"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
)

// NewMemoryInfoAgent creates an agent that gathers real memory usage information.
// This agent runs in parallel with other system information gatherers and uses
// gopsutil to gather actual memory metrics from the system.
func NewMemoryInfoAgent(ctx context.Context, model model.LLM) (agent.Agent, error) {
	// Create the memory info tool
	memoryInfoTool, err := tools.NewGetMemoryInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to create memory info tool: %w", err)
	}

	memoryInfoAgent, err := llmagent.New(llmagent.Config{
		Name:        "MemoryInfoAgent",
		Model:       model,
		Description: "Gathers real memory usage information and analyzes memory pressure using system tools",
		Instruction: `You are a Memory Information Specialist with access to real system metrics.

Your task is to:
1. Use the get_memory_info tool to gather REAL memory data from the system
2. Analyze the memory metrics you receive
3. Provide a detailed report including:
   - Total RAM capacity
   - Current memory usage and available memory
   - Memory usage percentage
   - Swap usage and swap percentage
   - Memory pressure indicators
   - Potential memory issues or bottlenecks
   - Recommendations for memory optimization

IMPORTANT:
- Always call the get_memory_info tool first to get real system data
- Base your analysis on the ACTUAL data returned by the tool
- Do not simulate or make up data - use only the real metrics provided
- Pay special attention to high memory usage (>80%) or swap usage

Store your memory analysis in state with the key "memory_info_report".`,
		OutputKey: "memory_info_report",
		Tools: []tool.Tool{
			memoryInfoTool,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create memory info agent: %w", err)
	}

	return memoryInfoAgent, nil
}
