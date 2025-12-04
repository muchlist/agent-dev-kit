// Package agents implements the sub-agents for the lead qualification sequential pipeline.
package agents

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
)

// NewActionRecommender creates an agent that recommends next actions based on lead qualification.
// This agent uses the validation and scoring results to suggest appropriate follow-up actions.
func NewActionRecommender(ctx context.Context, model model.LLM) (agent.Agent, error) {
	recommender, err := llmagent.New(llmagent.Config{
		Name:        "ActionRecommenderAgent",
		Model:       model,
		Description: "Recommends next actions based on lead qualification results",
		Instruction: `You are an Action Recommendation AI.

Based on the lead information and scoring:

- For invalid leads: Suggest what additional information is needed
- For leads scored 1-3: Suggest nurturing actions (educational content, etc.)
- For leads scored 4-7: Suggest qualifying actions (discovery call, needs assessment)
- For leads scored 8-10: Suggest sales actions (demo, proposal, etc.)

Format your response as a complete recommendation to the sales team.

You can access previous results from state:
- validation_status: Lead validation result
- lead_score: Lead scoring result

Store your recommendation in state with the key "action_recommendation".`,
		OutputKey: "action_recommendation",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create action recommender agent: %w", err)
	}

	return recommender, nil
}