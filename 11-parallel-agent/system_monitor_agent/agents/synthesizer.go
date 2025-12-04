// Package agents implements the sub-agents for the system monitor parallel workflow.
package agents

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
)

// NewSystemReportSynthesizer creates an agent that combines all gathered information into a comprehensive report.
// This agent runs after the parallel information gathering is complete.
func NewSystemReportSynthesizer(ctx context.Context, model model.LLM) (agent.Agent, error) {
	reportSynthesizer, err := llmagent.New(llmagent.Config{
		Name:        "SystemReportSynthesizer",
		Model:       model,
		Description: "Combines parallel system information into a comprehensive health report",
		Instruction: `You are a System Report Synthesizer.

Combine the system information gathered by the parallel agents into a comprehensive system health report. You have access to:

CPU Information: {state.cpu_info_report}
Memory Information: {state.memory_info_report}
Disk Information: {state.disk_info_report}

Create a well-structured report that includes:

EXECUTIVE SUMMARY:
- Overall system health status
- Key metrics and their implications
- Critical issues requiring immediate attention

DETAILED ANALYSIS:
- CPU performance and utilization
- Memory usage and pressure indicators
- Disk space and storage health
- Performance bottlenecks or concerns

RECOMMENDATIONS:
- Immediate actions needed
- Optimization suggestions
- Preventive maintenance recommendations
- Future upgrade considerations

Format the report professionally with clear sections and actionable insights. Make it easy to understand for both technical and non-technical users.

Store your comprehensive report in state with the key "system_health_report".`,
		OutputKey: "system_health_report",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create system report synthesizer agent: %w", err)
	}

	return reportSynthesizer, nil
}