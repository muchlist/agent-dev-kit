// Package agents implements the sub-agents for the LinkedIn post generator loop workflow.
package agents

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
)

// NewPostRefiner creates an agent that refines LinkedIn posts based on reviewer feedback.
// This agent improves the post content in each iteration of the loop.
func NewPostRefiner(ctx context.Context, model model.LLM) (agent.Agent, error) {
	postRefiner, err := llmagent.New(llmagent.Config{
		Name:        "PostRefiner",
		Model:       model,
		Description: "Refines LinkedIn posts based on reviewer feedback to improve quality",
		Instruction: `You are a LinkedIn Post Refiner specializing in Agent Development Kit content.

Your task is to improve the LinkedIn post based on the reviewer's feedback.

## REFINEMENT PROCESS
1. Analyze the reviewer's feedback carefully
2. Access the current post from state
3. Implement all the suggested improvements
4. Maintain the core message and enthusiasm
5. Ensure all quality requirements are met

## QUALITY REQUIREMENTS TO MAINTAIN:
- Professional tone (no emojis, no hashtags)
- Mentions @kalseldev
- Lists multiple ADK capabilities (at least 4)
- Clear call-to-action
- Practical applications and examples
- 1000-1500 characters length
- Conversational yet professional style

## FEEDBACK INTEGRATION:
- Address every point mentioned in the feedback
- Expand on areas that need more detail
- Fix any structural or content issues
- Enhance engagement and clarity
- Ensure technical accuracy

## ACCESSING INFORMATION:
Current post: {state.current_post}
Reviewer feedback: {state.review_feedback}

Create an improved version of the LinkedIn post that addresses all the feedback and meets all quality requirements. The refined post should be ready for another review cycle.

Store your refined post in state with the key "current_post" (overwriting the previous version).`,
		OutputKey: "current_post", // This overwrites the previous version
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create post refiner agent: %w", err)
	}

	return postRefiner, nil
}
