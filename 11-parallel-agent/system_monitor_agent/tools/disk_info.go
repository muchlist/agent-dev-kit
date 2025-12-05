// Package tools implements real system information gathering tools using gopsutil.
package tools

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// DiskInfoArgs represents the input arguments for disk info gathering
type DiskInfoArgs struct{}

// DiskInfoResults represents the result from disk info gathering
type DiskInfoResults struct {
	Result         DiskInfo       `json:"result"`
	Stats          DiskStats      `json:"stats"`
	AdditionalInfo AdditionalInfo `json:"additional_info"`
}

// DiskInfo contains detailed disk information
type DiskInfo struct {
	TotalSpace      string   `json:"total_space"`
	UsedSpace       string   `json:"used_space"`
	FreeSpace       string   `json:"free_space"`
	UsagePercentage string   `json:"usage_percentage"`
	MountPoint      string   `json:"mount_point"`
	FileSystem      string   `json:"file_system"`
	Partitions      []string `json:"partitions,omitempty"`
}

// DiskStats contains disk statistics
type DiskStats struct {
	UsagePercentage float64 `json:"usage_percentage"`
	TotalSpaceGB    float64 `json:"total_space_gb"`
	FreeSpaceGB     float64 `json:"free_space_gb"`
	UsedSpaceGB     float64 `json:"used_space_gb"`
}

// NewGetDiskInfo creates a tool to gather real disk information using gopsutil.
// This tool collects actual disk usage from the system.
func NewGetDiskInfo() (tool.Tool, error) {
	getDiskInfo := func(ctx tool.Context, input DiskInfoArgs) (DiskInfoResults, error) {
		fmt.Println("\n🔧 Tool: get_disk_info called - gathering real disk metrics")

		// Determine root path based on OS
		mountPoint := "/"
		if runtime.GOOS == "windows" {
			mountPoint = "C:"
		}

		// Get disk usage for the primary mount point
		usage, err := disk.Usage(mountPoint)
		if err != nil {
			return DiskInfoResults{}, fmt.Errorf("failed to get disk usage: %w", err)
		}

		// Convert bytes to GB
		totalGB := float64(usage.Total) / (1024 * 1024 * 1024)
		usedGB := float64(usage.Used) / (1024 * 1024 * 1024)
		freeGB := float64(usage.Free) / (1024 * 1024 * 1024)

		// Get partition information
		partitions, err := disk.Partitions(false)
		var partitionInfo []string
		if err == nil {
			for _, partition := range partitions {
				partitionInfo = append(partitionInfo, fmt.Sprintf("%s (%s)", partition.Device, partition.Mountpoint))
			}
		}

		diskInfo := DiskInfo{
			TotalSpace:      fmt.Sprintf("%.2f GB", totalGB),
			UsedSpace:       fmt.Sprintf("%.2f GB", usedGB),
			FreeSpace:       fmt.Sprintf("%.2f GB", freeGB),
			UsagePercentage: fmt.Sprintf("%.1f%%", usage.UsedPercent),
			MountPoint:      mountPoint,
			FileSystem:      usage.Fstype,
			Partitions:      partitionInfo,
		}

		stats := DiskStats{
			UsagePercentage: usage.UsedPercent,
			TotalSpaceGB:    totalGB,
			FreeSpaceGB:     freeGB,
			UsedSpaceGB:     usedGB,
		}

		// Check for disk space concerns
		highDiskUsage := usage.UsedPercent > 80
		var diskConcern *string
		if highDiskUsage {
			concern := "High disk usage detected"
			diskConcern = &concern
		}

		additionalInfo := AdditionalInfo{
			DataFormat:          "dictionary",
			CollectionTimestamp: float64(time.Now().Unix()),
			DiskSpaceConcern:    diskConcern,
		}

		fmt.Printf("   ✓ Collected: %.2f GB total, %.2f GB free, %.1f%% used\n",
			totalGB, freeGB, usage.UsedPercent)

		return DiskInfoResults{
			Result:         diskInfo,
			Stats:          stats,
			AdditionalInfo: additionalInfo,
		}, nil
	}

	return functiontool.New(
		functiontool.Config{
			Name:        "get_disk_info",
			Description: "Gather real disk information including space usage and partitions from the system",
		},
		getDiskInfo,
	)
}
