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

// NewDiskInfoAgent creates an agent that gathers real disk space information.
// This agent runs in parallel with other system information gatherers and uses
// gopsutil to gather actual disk metrics from the system.
func NewDiskInfoAgent(ctx context.Context, model model.LLM) (agent.Agent, error) {
	// Create the disk info tool
	diskInfoTool, err := tools.NewGetDiskInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to create disk info tool: %w", err)
	}

	diskInfoAgent, err := llmagent.New(llmagent.Config{
		Name:        "DiskInfoAgent",
		Model:       model,
		Description: "Gathers real disk space and partition information using system tools",
		Instruction: `You are a Disk Information Specialist with access to real system metrics.

Your task is to:
1. Use the get_disk_info tool to gather REAL disk data from the system
2. Analyze the disk metrics you receive
3. Provide a comprehensive report including:
   - Total disk capacity
   - Used and free disk space
   - Disk usage percentage
   - File system type and mount points
   - Available partitions
   - Disk space warnings or concerns
   - Recommendations for disk space management

IMPORTANT:
- Always call the get_disk_info tool first to get real system data
- Base your analysis on the ACTUAL data returned by the tool
- Do not simulate or make up data - use only the real metrics provided
- Pay special attention to high disk usage (>80%)
- Provide actionable recommendations if disk space is low

Store your disk analysis in state with the key "disk_info_report".`,
		OutputKey: "disk_info_report",
		Tools: []tool.Tool{
			diskInfoTool,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create disk info agent: %w", err)
	}

	return diskInfoAgent, nil
}
