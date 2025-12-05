// Package tools implements real system information gathering tools using gopsutil.
package tools

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// CPUInfoArgs represents the input arguments for CPU info gathering
type CPUInfoArgs struct{}

// CPUInfoResults represents the result from CPU info gathering
type CPUInfoResults struct {
	Result         CPUInfo       `json:"result"`
	Stats          CPUStats      `json:"stats"`
	AdditionalInfo AdditionalInfo `json:"additional_info"`
}

// CPUInfo contains detailed CPU information
type CPUInfo struct {
	PhysicalCores   int      `json:"physical_cores"`
	LogicalCores    int      `json:"logical_cores"`
	CPUUsagePerCore []string `json:"cpu_usage_per_core"`
	AvgCPUUsage     string   `json:"avg_cpu_usage"`
}

// CPUStats contains CPU statistics
type CPUStats struct {
	PhysicalCores      int     `json:"physical_cores"`
	LogicalCores       int     `json:"logical_cores"`
	AvgUsagePercentage float64 `json:"avg_usage_percentage"`
	HighUsageAlert     bool    `json:"high_usage_alert"`
}

// AdditionalInfo contains metadata about the data collection
type AdditionalInfo struct {
	DataFormat           string  `json:"data_format"`
	CollectionTimestamp  float64 `json:"collection_timestamp"`
	PerformanceConcern   *string `json:"performance_concern,omitempty"`
	SwapConcern          *string `json:"swap_concern,omitempty"`
	DiskSpaceConcern     *string `json:"disk_space_concern,omitempty"`
}

// NewGetCPUInfo creates a tool to gather real CPU information using gopsutil.
// This tool collects actual CPU metrics from the system.
func NewGetCPUInfo() (tool.Tool, error) {
	getCPUInfo := func(ctx tool.Context, input CPUInfoArgs) (CPUInfoResults, error) {
		fmt.Println("\n🔧 Tool: get_cpu_info called - gathering real CPU metrics")

		// Get CPU counts
		physicalCount, err := cpu.Counts(false)
		if err != nil {
			return CPUInfoResults{}, fmt.Errorf("failed to get physical CPU count: %w", err)
		}

		logicalCount, err := cpu.Counts(true)
		if err != nil {
			return CPUInfoResults{}, fmt.Errorf("failed to get logical CPU count: %w", err)
		}

		// Get per-core CPU usage (with 1 second interval for accuracy)
		perCPU, err := cpu.Percent(time.Second, true)
		if err != nil {
			return CPUInfoResults{}, fmt.Errorf("failed to get per-CPU usage: %w", err)
		}

		// Format per-core usage
		var cpuUsagePerCore []string
		var totalUsage float64
		for i, percentage := range perCPU {
			cpuUsagePerCore = append(cpuUsagePerCore, fmt.Sprintf("Core %d: %.1f%%", i, percentage))
			totalUsage += percentage
		}

		// Calculate average usage
		avgUsage := totalUsage / float64(len(perCPU))
		highUsage := avgUsage > 80

		// Performance concern
		var performanceConcern *string
		if highUsage {
			concern := "High CPU usage detected"
			performanceConcern = &concern
		}

		cpuInfo := CPUInfo{
			PhysicalCores:   physicalCount,
			LogicalCores:    logicalCount,
			CPUUsagePerCore: cpuUsagePerCore,
			AvgCPUUsage:     fmt.Sprintf("%.1f%%", avgUsage),
		}

		stats := CPUStats{
			PhysicalCores:      physicalCount,
			LogicalCores:       logicalCount,
			AvgUsagePercentage: avgUsage,
			HighUsageAlert:     highUsage,
		}

		additionalInfo := AdditionalInfo{
			DataFormat:          "dictionary",
			CollectionTimestamp: float64(time.Now().Unix()),
			PerformanceConcern:  performanceConcern,
		}

		fmt.Printf("   ✓ Collected: %d physical cores, %d logical cores, avg usage: %.1f%%\n",
			physicalCount, logicalCount, avgUsage)

		return CPUInfoResults{
			Result:         cpuInfo,
			Stats:          stats,
			AdditionalInfo: additionalInfo,
		}, nil
	}

	return functiontool.New(
		functiontool.Config{
			Name:        "get_cpu_info",
			Description: "Gather real CPU information including core count and usage statistics from the system",
		},
		getCPUInfo,
	)
}
