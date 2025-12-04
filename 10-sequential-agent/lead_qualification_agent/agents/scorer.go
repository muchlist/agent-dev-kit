// Package agents implements the sub-agents for the lead qualification sequential pipeline.
package agents

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
)

// NewLeadScorer creates an agent that scores qualified leads on a scale of 1-10.
// This agent analyzes various criteria to determine lead qualification level.
func NewLeadScorer(ctx context.Context, model model.LLM) (agent.Agent, error) {
	scorer, err := llmagent.New(llmagent.Config{
		Name:        "LeadScorerAgent",
		Model:       model,
		Description: "Scores qualified leads on a scale of 1-10 based on qualification criteria",
		Instruction: `You are a Lead Scoring AI.

Analyze the lead information and assign a qualification score from 1-10 based on:
- Expressed need (urgency/clarity of problem)
- Decision-making authority
- Budget indicators
- Timeline indicators

Output ONLY a numeric score and ONE sentence justification.

Example output: '8: Decision maker with clear budget and immediate need'
Example output: '3: Vague interest with no timeline or budget mentioned'

You can access the validation status from previous step using state if needed.
Store your scoring result in state with the key "lead_score".`,
		OutputKey: "lead_score",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create lead scorer agent: %w", err)
	}

	return scorer, nil
}