// Package tools implements real system information gathering tools using gopsutil.
package tools

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// MemoryInfoArgs represents the input arguments for memory info gathering
type MemoryInfoArgs struct{}

// MemoryInfoResults represents the result from memory info gathering
type MemoryInfoResults struct {
	Result         MemoryInfo     `json:"result"`
	Stats          MemoryStats    `json:"stats"`
	AdditionalInfo AdditionalInfo `json:"additional_info"`
}

// MemoryInfo contains detailed memory information
type MemoryInfo struct {
	TotalMemory      string `json:"total_memory"`
	AvailableMemory  string `json:"available_memory"`
	UsedMemory       string `json:"used_memory"`
	MemoryPercentage string `json:"memory_percentage"`
	SwapTotal        string `json:"swap_total"`
	SwapUsed         string `json:"swap_used"`
	SwapPercentage   string `json:"swap_percentage"`
}

// MemoryStats contains memory statistics
type MemoryStats struct {
	MemoryUsagePercentage float64 `json:"memory_usage_percentage"`
	SwapUsagePercentage   float64 `json:"swap_usage_percentage"`
	TotalMemoryGB         float64 `json:"total_memory_gb"`
	AvailableMemoryGB     float64 `json:"available_memory_gb"`
}

// NewGetMemoryInfo creates a tool to gather real memory information using gopsutil.
// This tool collects actual RAM and swap usage from the system.
func NewGetMemoryInfo() (tool.Tool, error) {
	getMemoryInfo := func(ctx tool.Context, input MemoryInfoArgs) (MemoryInfoResults, error) {
		fmt.Println("\n🔧 Tool: get_memory_info called - gathering real memory metrics")

		// Get virtual memory information
		vmStat, err := mem.VirtualMemory()
		if err != nil {
			return MemoryInfoResults{}, fmt.Errorf("failed to get virtual memory stats: %w", err)
		}

		// Get swap memory information
		swapStat, err := mem.SwapMemory()
		if err != nil {
			return MemoryInfoResults{}, fmt.Errorf("failed to get swap memory stats: %w", err)
		}

		// Convert bytes to GB
		totalGB := float64(vmStat.Total) / (1024 * 1024 * 1024)
		availableGB := float64(vmStat.Available) / (1024 * 1024 * 1024)
		usedGB := float64(vmStat.Used) / (1024 * 1024 * 1024)
		swapTotalGB := float64(swapStat.Total) / (1024 * 1024 * 1024)
		swapUsedGB := float64(swapStat.Used) / (1024 * 1024 * 1024)

		memoryInfo := MemoryInfo{
			TotalMemory:      fmt.Sprintf("%.2f GB", totalGB),
			AvailableMemory:  fmt.Sprintf("%.2f GB", availableGB),
			UsedMemory:       fmt.Sprintf("%.2f GB", usedGB),
			MemoryPercentage: fmt.Sprintf("%.1f%%", vmStat.UsedPercent),
			SwapTotal:        fmt.Sprintf("%.2f GB", swapTotalGB),
			SwapUsed:         fmt.Sprintf("%.2f GB", swapUsedGB),
			SwapPercentage:   fmt.Sprintf("%.1f%%", swapStat.UsedPercent),
		}

		stats := MemoryStats{
			MemoryUsagePercentage: vmStat.UsedPercent,
			SwapUsagePercentage:   swapStat.UsedPercent,
			TotalMemoryGB:         totalGB,
			AvailableMemoryGB:     availableGB,
		}

		// Check for concerns
		highMemoryUsage := vmStat.UsedPercent > 80
		highSwapUsage := swapStat.UsedPercent > 80

		var memConcern, swapConcern *string
		if highMemoryUsage {
			concern := "High memory usage detected"
			memConcern = &concern
		}
		if highSwapUsage {
			concern := "High swap usage detected"
			swapConcern = &concern
		}

		additionalInfo := AdditionalInfo{
			DataFormat:          "dictionary",
			CollectionTimestamp: float64(time.Now().Unix()),
			PerformanceConcern:  memConcern,
			SwapConcern:         swapConcern,
		}

		fmt.Printf("   ✓ Collected: %.2f GB total, %.2f GB available, %.1f%% used\n",
			totalGB, availableGB, vmStat.UsedPercent)

		return MemoryInfoResults{
			Result:         memoryInfo,
			Stats:          stats,
			AdditionalInfo: additionalInfo,
		}, nil
	}

	return functiontool.New(
		functiontool.Config{
			Name:        "get_memory_info",
			Description: "Gather real memory information including RAM and swap usage from the system",
		},
		getMemoryInfo,
	)
}
