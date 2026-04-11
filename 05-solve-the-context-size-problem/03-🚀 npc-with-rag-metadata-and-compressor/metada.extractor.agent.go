package main

import (
	"context"

	"github.com/snipwise/nova/nova-sdk/agents"
	"github.com/snipwise/nova/nova-sdk/agents/structured"
	"github.com/snipwise/nova/nova-sdk/models"
	"github.com/snipwise/nova/nova-sdk/ui/display"
)

// KeywordMetadata represents extracted keywords from content
type KeywordMetadata struct {
	Keywords  []string `json:"keywords"`
	MainTopic string   `json:"main_topic"`
	Category  string   `json:"category"`
}

func getMetadataExtractorAgent(ctx context.Context, engineURL, metadataModel string) (*structured.Agent[KeywordMetadata], error) {

	// Create structured agent for keyword extraction
	structuredAgent, err := structured.NewAgent[KeywordMetadata](
		ctx,
		agents.Config{
			EngineURL: engineURL,
		},
		models.Config{
			Name: metadataModel,
		},
	)
	if err != nil {
		display.Errorf("❌ Error creating structured agent: %v", err)
		panic(err)
	}
	display.Infof("✅ Structured agent created for keyword extraction")
	
	return structuredAgent, nil
}
