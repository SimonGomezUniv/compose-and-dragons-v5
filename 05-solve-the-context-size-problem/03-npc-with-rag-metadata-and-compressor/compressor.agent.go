package main

import (
	"context"

	"github.com/snipwise/nova/nova-sdk/agents"
	"github.com/snipwise/nova/nova-sdk/agents/compressor"
	"github.com/snipwise/nova/nova-sdk/models"
)

// getCompressorAgent creates and returns a Nova compressor agent
// for context compression to save tokens while preserving key information
func getCompressorAgent(ctx context.Context, engineURL, compressorModel string) (*compressor.Agent, error) {
	// Create compressor agent with specialized system instructions
	compressorAgent, err := compressor.NewAgent(
		ctx,
		agents.Config{
			Name:      "compressor",
			EngineURL: engineURL,
			SystemInstructions: `You are an expert at summarizing and compressing conversations.
			Your role is to create concise summaries that preserve:
			- Key information and important facts
			- Decisions made
			- User preferences
			- Emotional context if relevant
			- Ongoing or pending actions

			Output format:
			## Conversation Summary
			[Concise summary of exchanges]

			## Key Points
			- [Point 1]
			- [Point 2]

			## To Remember
			[Important information for continuity]
			`,
		},
		models.Config{
			Name:        compressorModel,
			Temperature: models.Float64(0.0), // Deterministic for consistent summaries
		},
		compressor.WithCompressionPrompt(compressor.Prompts.UltraShort),
	)
	if err != nil {
		return nil, err
	}

	return compressorAgent, nil
}
