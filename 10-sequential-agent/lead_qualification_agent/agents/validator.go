// Package agents implements the sub-agents for the lead qualification sequential pipeline.
package agents

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
)

// NewLeadValidator creates an agent that validates lead information for completeness.
// This agent checks if a lead has sufficient information to proceed with qualification.
func NewLeadValidator(ctx context.Context, model model.LLM) (agent.Agent, error) {
	validator, err := llmagent.New(llmagent.Config{
		Name:        "LeadValidatorAgent",
		Model:       model,
		Description: "Validates lead information for completeness",
		Instruction: `You are a Lead Validation AI.

Examine the lead information provided by the user and determine if it's complete enough for qualification.
A complete lead should include:
- Contact information (name, email or phone)
- Some indication of interest or need
- Company or context information if applicable

Output ONLY 'valid' or 'invalid' with a single reason if invalid.

Example valid output: 'valid'
Example invalid output: 'invalid: missing contact information'

Store your validation result in state with the key "validation_status".`,
		OutputKey: "validation_status",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create lead validator agent: %w", err)
	}

	return validator, nil
}