// Package tools implements tools for the LinkedIn post generator loop workflow.
package tools

import (
	"log"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// ExitLoopArgs represents the input arguments for the exit loop tool
type ExitLoopArgs struct {
}

// ExitLoopResult represents the result from the exit loop tool
type ExitLoopResult struct {
	Success bool `json:"success"`
}

// NewExitLoop creates a tool to exit the loop when quality requirements are met.
// This tool signals the LoopAgent to stop iterating by setting escalate=true.
func NewExitLoop() (tool.Tool, error) {
	exitLoop := func(ctx tool.Context, args ExitLoopArgs) (ExitLoopResult, error) {
		log.Printf("\n----------- EXIT LOOP TRIGGERED -----------")
		log.Printf("Post review completed successfully")
		log.Printf("Loop will exit now")
		log.Printf("------------------------------------------\n")

		// Signal to the LoopAgent that we should stop iterating
		ctx.Actions().Escalate = true
		return ExitLoopResult{Success: true}, nil
	}

	return functiontool.New(
		functiontool.Config{
			Name:        "exit_loop",
			Description: "Call this function ONLY when the post meets all quality requirements, signaling the iterative process should end",
		},
		exitLoop,
	)
}