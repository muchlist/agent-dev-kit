package tools

import (
	"fmt"
	"time"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// ===== Time Tool Structures =====

type getCurrentTimeArgs struct{}

type getCurrentTimeResults struct {
	CurrentTime string `json:"current_time"`
}

// ===== Tool Implementation =====

// getCurrentTime returns the current time in YYYY-MM-DD HH:MM:SS format
func getCurrentTime(ctx tool.Context, input getCurrentTimeArgs) (getCurrentTimeResults, error) {
	fmt.Println("--- Tool: get_current_time called ---")
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	return getCurrentTimeResults{
		CurrentTime: currentTime,
	}, nil
}

// ===== Tool Creation =====

// NewGetCurrentTimeTool creates a tool for getting the current time
func NewGetCurrentTimeTool() (tool.Tool, error) {
	getCurrentTimeTool, err := functiontool.New(
		functiontool.Config{
			Name:        "get_current_time",
			Description: "Get the current time in the format YYYY-MM-DD HH:MM:SS",
		},
		getCurrentTime)
	if err != nil {
		return nil, fmt.Errorf("failed to create get_current_time tool: %w", err)
	}

	return getCurrentTimeTool, nil
}
