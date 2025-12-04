// Package main implements a system monitor parallel agent in Go.
// This example demonstrates how to create a hybrid workflow using both Parallel and Sequential agents.
//
// The system monitoring workflow:
// 1. Parallel Information Gathering: Concurrently collect CPU, Memory, and Disk information
// 2. Sequential Report Synthesis: Combine all information into a comprehensive report
//
// This hybrid approach shows how to combine workflow agent types for optimal performance
// and logical flow - parallel for independent tasks, sequential for dependent processing.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/parallelagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"

	"github.com/muchlist/agent-dev-kit/11-parallel-agent/system_monitor_agent/agents"
)

const (
	MODEL_NAME = "gemini-2.0-flash"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	// Create the Gemini model
	model, err := gemini.NewModel(ctx, MODEL_NAME, &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	fmt.Println("🖥️  Creating System Monitor Parallel Agent...")

	// Create sub-agents for parallel system information gathering
	cpuInfoAgent, err := agents.NewCPUInfoAgent(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create CPU info agent: %v", err)
	}
	fmt.Println("  ✓ CPU Info Agent created")

	memoryInfoAgent, err := agents.NewMemoryInfoAgent(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create memory info agent: %v", err)
	}
	fmt.Println("  ✓ Memory Info Agent created")

	diskInfoAgent, err := agents.NewDiskInfoAgent(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create disk info agent: %v", err)
	}
	fmt.Println("  ✓ Disk Info Agent created")

	// Create report synthesizer agent
	reportSynthesizer, err := agents.NewSystemReportSynthesizer(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create report synthesizer agent: %v", err)
	}
	fmt.Println("  ✓ System Report Synthesizer created")

	// Create Parallel Agent for concurrent system information gathering
	fmt.Println("⚡ Creating Parallel Information Gatherer...")
	parallelInfoGatherer, err := parallelagent.New(parallelagent.Config{
		AgentConfig: agent.Config{
			Name:        "system_info_gatherer",
			Description: "Gathers system information concurrently from CPU, memory, and disk",
			SubAgents:   []agent.Agent{cpuInfoAgent, memoryInfoAgent, diskInfoAgent},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create parallel info gatherer: %v", err)
	}
	fmt.Println("  ✓ Parallel Information Gatherer created")

	// Create Sequential Agent for the overall workflow
	fmt.Println("🔗 Creating Sequential Pipeline...")
	sequentialAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "system_monitor_agent",
			Description: "Monitors system health using parallel data gathering and sequential synthesis",
			SubAgents:   []agent.Agent{parallelInfoGatherer, reportSynthesizer},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create system monitor sequential agent: %v", err)
	}
	fmt.Println("  ✓ System Monitor Sequential Agent created")

	fmt.Println("\n🚀 Launching System Monitor Parallel Agent...")
	fmt.Println("========================================================")
	fmt.Println("Example prompts to try:")
	fmt.Println("• 'Check my system health'")
	fmt.Println("• 'Provide a comprehensive system report with recommendations'")
	fmt.Println("• 'Is my system running out of memory or disk space?'")
	fmt.Println("• 'Generate a detailed system status report'")
	fmt.Println("========================================================")

	// Configure and launch the agent
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(sequentialAgent),
	}

	l := full.NewLauncher()
	if err := l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}