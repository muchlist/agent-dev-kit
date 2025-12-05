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

// NewCPUInfoAgent creates an agent that collects and analyzes real CPU information.
// This agent runs in parallel with other system information gatherers and uses
// gopsutil to gather actual CPU metrics from the system.
func NewCPUInfoAgent(ctx context.Context, model model.LLM) (agent.Agent, error) {
	// Create the CPU info tool
	cpuInfoTool, err := tools.NewGetCPUInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to create CPU info tool: %w", err)
	}

	cpuInfoAgent, err := llmagent.New(llmagent.Config{
		Name:        "CPUInfoAgent",
		Model:       model,
		Description: "Collects and analyzes real CPU information and performance metrics using system tools",
		Instruction: `You are a CPU Information Specialist with access to real system metrics.

Your task is to:
1. Use the get_cpu_info tool to gather REAL CPU data from the system
2. Analyze the CPU metrics you receive
3. Provide a comprehensive report including:
   - CPU model and architecture (from tool results)
   - Number of physical and logical cores
   - Current usage statistics per core
   - Average CPU usage
   - Performance indicators and trends
   - Any potential issues (high usage, bottlenecks, etc.)
   - Recommendations for optimization if needed

IMPORTANT:
- Always call the get_cpu_info tool first to get real system data
- Base your analysis on the ACTUAL data returned by the tool
- Do not simulate or make up data - use only the real metrics provided

Store your CPU analysis in state with the key "cpu_info_report".`,
		OutputKey: "cpu_info_report",
		Tools: []tool.Tool{
			cpuInfoTool,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create CPU info agent: %w", err)
	}

	return cpuInfoAgent, nil
}
