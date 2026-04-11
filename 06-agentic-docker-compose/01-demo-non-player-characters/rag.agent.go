package main

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/snipwise/nova/nova-sdk/agents"
	"github.com/snipwise/nova/nova-sdk/agents/rag"
	"github.com/snipwise/nova/nova-sdk/agents/structured"
	"github.com/snipwise/nova/nova-sdk/models"
	"github.com/snipwise/nova/nova-sdk/toolbox/env"
	"github.com/snipwise/nova/nova-sdk/ui/display"
)

// getRagAgent creates or loads a RAG agent with JSON store persistence
func getRagAgent(ctx context.Context, engineURL, embeddingModel, sheetFilePath string, metadataExtractorAgent *structured.Agent[KeywordMetadata]) (*rag.Agent, error) {
	// === CONFIGURATION ===
	// Extract filename without extension from sheetFilePath
	baseName := filepath.Base(sheetFilePath)
	fileName := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	storePath := env.GetEnvOrDefault("STORE_PATH", "./store")

	storePathFile := filepath.Join(storePath, fileName+".json")

	display.Infof("📦 RAG Store path: %s", storePathFile)

	// === CREATE RAG AGENT ===
	ragAgent, err := rag.NewAgent(
		ctx,
		agents.Config{
			EngineURL: engineURL,
		},
		models.Config{
			Name: embeddingModel,
		},
	)
	if err != nil {
		display.Errorf("❌ Error creating RAG agent: %v", err)
		return nil, err
	}

	// === LOAD STORE ===
	if ragAgent.StoreFileExists(storePathFile) {
		// Load existing store
		err := ragAgent.LoadStore(storePathFile)
		if err != nil {
			display.Errorf("❌ Error loading store %s: %v", storePathFile, err)
			return nil, err
		}
		display.Infof("✅ RAG store loaded from %s", storePathFile)
	} else {
		err := errors.New("store file not found")
		display.Errorf("❌ Error: %v", err)
		return nil, err
	}

	return ragAgent, nil
}
