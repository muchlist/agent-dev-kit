// Package tools implements tools for the LinkedIn post generator loop workflow.
package tools

import (
	"fmt"
	"log"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// CharacterCounterArgs represents the input arguments for the character counter tool
type CharacterCounterArgs struct {
	Text string `json:"text"`
}

// CharacterCounterResult represents the result from the character counter tool
type CharacterCounterResult struct {
	Result        string `json:"result"`
	CharCount     int    `json:"char_count"`
	CharsNeeded   int    `json:"chars_needed,omitempty"`
	CharsToRemove int    `json:"chars_to_remove,omitempty"`
	Message       string `json:"message"`
}

// NewCharacterCounter creates a tool to count characters and provide length-based feedback.
// This tool helps validate LinkedIn post length requirements (1000-1500 characters).
func NewCharacterCounter() (tool.Tool, error) {
	charCounter := func(ctx tool.Context, args CharacterCounterArgs) (CharacterCounterResult, error) {
		charCount := len(args.Text)
		const (
			MIN_LENGTH = 1000
			MAX_LENGTH = 1500
		)

		log.Printf("\n----------- TOOL DEBUG -----------")
		log.Printf("Checking text length: %d characters", charCount)
		log.Printf("----------------------------------\n")

		// Update review status in state
		if charCount < MIN_LENGTH {
			charsNeeded := MIN_LENGTH - charCount
			ctx.State().Set("review_status", "fail")
			return CharacterCounterResult{
				Result:      "fail",
				CharCount:   charCount,
				CharsNeeded: charsNeeded,
				Message:     fmt.Sprintf("Post is too short. Add %d more characters to reach minimum length of %d.", charsNeeded, MIN_LENGTH),
			}, nil
		} else if charCount > MAX_LENGTH {
			charsToRemove := charCount - MAX_LENGTH
			ctx.State().Set("review_status", "fail")
			return CharacterCounterResult{
				Result:       "fail",
				CharCount:    charCount,
				CharsToRemove: charsToRemove,
				Message:      fmt.Sprintf("Post is too long. Remove %d characters to meet maximum length of %d.", charsToRemove, MAX_LENGTH),
			}, nil
		} else {
			ctx.State().Set("review_status", "pass")
			return CharacterCounterResult{
				Result:    "pass",
				CharCount: charCount,
				Message:   fmt.Sprintf("Post length is good (%d characters).", charCount),
			}, nil
		}
	}

	return functiontool.New(
		functiontool.Config{
			Name:        "count_characters",
			Description: "Counts characters in text and provides length-based feedback for LinkedIn posts",
		},
		charCounter,
	)
}