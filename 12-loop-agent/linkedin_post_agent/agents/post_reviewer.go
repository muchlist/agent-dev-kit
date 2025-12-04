// Package agents implements the sub-agents for the LinkedIn post generator loop workflow.
package agents

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"

	"github.com/muchlist/agent-dev-kit/12-loop-agent/linkedin_post_agent/tools"
)

// NewPostReviewer creates an agent that reviews LinkedIn posts for quality and can exit the loop.
// This agent evaluates posts against quality criteria and calls exit_loop when requirements are met.
func NewPostReviewer(ctx context.Context, model model.LLM) (agent.Agent, error) {
	// Create the tools for the post reviewer
	charCounterTool, err := tools.NewCharacterCounter()
	if err != nil {
		return nil, fmt.Errorf("failed to create character counter tool: %w", err)
	}

	exitLoopTool, err := tools.NewExitLoop()
	if err != nil {
		return nil, fmt.Errorf("failed to create exit loop tool: %w", err)
	}

	postReviewer, err := llmagent.New(llmagent.Config{
		Name:        "PostReviewer",
		Model:       model,
		Description: "Reviews post quality and provides feedback or exits loop when requirements are met",
		Instruction: `You are a LinkedIn Post Quality Reviewer.

Your task is to evaluate the quality of a LinkedIn post about Agent Development Kit (ADK).

## EVALUATION PROCESS
1. Use the count_characters tool to check the post's length.
   Pass the current post text from state to the tool.

2. If the length check fails (tool result is "fail"), provide specific feedback on what needs to be fixed.
   Use the tool's message as a guideline, but add your own professional critique.

3. If length check passes, evaluate the post against these criteria:

   REQUIRED ELEMENTS:
   1. Mentions @kalseldev
   2. Lists multiple ADK capabilities (at least 4)
   3. Has a clear call-to-action
   4. Includes practical applications
   5. Shows genuine enthusiasm

   STYLE REQUIREMENTS:
   1. NO emojis
   2. NO hashtags
   3. Professional tone
   4. Conversational style
   5. Clear and concise writing

## OUTPUT INSTRUCTIONS
IF the post fails ANY of the checks above:
  - Return concise, specific feedback on what to improve

ELSE IF the post meets ALL requirements:
  - Call the exit_loop function
  - Return "Post meets all requirements. Exiting the refinement loop."

Access the current post from state: {state.current_post}

Do not embellish your response. Either provide feedback on what to improve OR call exit_loop and return the completion message.`,
		Tools:     []tool.Tool{charCounterTool, exitLoopTool},
		OutputKey: "review_feedback",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create post reviewer agent: %w", err)
	}

	return postReviewer, nil
}
