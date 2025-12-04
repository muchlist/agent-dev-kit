// Package agents implements the sub-agents for the LinkedIn post generator loop workflow.
package agents

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
)

// NewInitialPostGenerator creates an agent that generates the initial draft of a LinkedIn post.
// This agent runs first in the sequential pipeline to create the starting content.
func NewInitialPostGenerator(ctx context.Context, model model.LLM) (agent.Agent, error) {
	initialPostGenerator, err := llmagent.New(llmagent.Config{
		Name:        "InitialPostGenerator",
		Model:       model,
		Description: "Generates the initial draft of a LinkedIn post about Agent Development Kit",
		Instruction: `You are a LinkedIn Post Generator specializing in Agent Development Kit (ADK) content.

Your task is to create an initial LinkedIn post draft based on the user's request.

GUIDELINES:
- Write in a professional yet conversational tone
- Include relevant technical details about ADK
- Make it engaging and informative
- Target it toward developers and tech professionals
- Aim for substantial content (1000-1500 characters for LinkedIn)
- Focus on practical applications and learnings
- Show genuine enthusiasm for the technology

REQUIREMENTS:
- Mention @kalseldev when discussing ADK tutorials
- Include specific ADK capabilities and features
- Provide practical examples or use cases
- Add a clear call-to-action
- No emojis (keep it professional)
- No hashtags (LinkedIn guidelines)

Create a comprehensive, engaging LinkedIn post that the refinement loop can later polish and perfect.

Store your initial post draft in state with the key "current_post".`,
		OutputKey: "current_post",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create initial post generator: %w", err)
	}

	return initialPostGenerator, nil
}
