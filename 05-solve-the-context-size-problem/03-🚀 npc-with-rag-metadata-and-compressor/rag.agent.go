package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/snipwise/nova/nova-sdk/agents"
	"github.com/snipwise/nova/nova-sdk/agents/rag"
	"github.com/snipwise/nova/nova-sdk/agents/rag/chunks"
	"github.com/snipwise/nova/nova-sdk/agents/structured"
	"github.com/snipwise/nova/nova-sdk/messages"
	"github.com/snipwise/nova/nova-sdk/messages/roles"
	"github.com/snipwise/nova/nova-sdk/models"
	"github.com/snipwise/nova/nova-sdk/toolbox/files"
	"github.com/snipwise/nova/nova-sdk/ui/display"
)

// getRagAgent creates or loads a RAG agent with JSON store persistence
func getRagAgent(ctx context.Context, engineURL, embeddingModel, sheetFilePath string, metadataExtractorAgent *structured.Agent[KeywordMetadata]) (*rag.Agent, error) {
	// === CONFIGURATION ===
	// Extract filename without extension from sheetFilePath
	baseName := filepath.Base(sheetFilePath)
	fileName := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	storePathFile := filepath.Join("./store", fileName+".json")

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

	// === LOAD OR CREATE STORE ===
	if ragAgent.StoreFileExists(storePathFile) {
		// Load existing store
		err := ragAgent.LoadStore(storePathFile)
		if err != nil {
			display.Errorf("❌ Error loading store %s: %v", storePathFile, err)
			return nil, err
		}
		display.Infof("✅ RAG store loaded from %s", storePathFile)
	} else {
		display.Infof("📝 Store not found. Creating new store and indexing character sheet...")

		// Read character sheet content
		characterSheetContent, err := files.ReadTextFile(sheetFilePath)
		if err != nil {
			display.Errorf("❌ Error reading character sheet: %v", err)
			return nil, err
		}

		// === CHUNK AND INDEX CONTENT ===
		// Split markdown by sections
		contentPieces := chunks.SplitMarkdownBySections(characterSheetContent)
		display.Infof("📄 Split character sheet into %d sections", len(contentPieces))

		// Index each section
		for idx, piece := range contentPieces[1:] { // Skip title section
			// === EXTRACT METADATA FOR SECTION ===
			// Extract keywords and metadata using structured agent
			extractionPrompt := fmt.Sprintf(`Analyze the following content and extract relevant metadata.
			Content:
			%s

			Extract:
			- Keywords: only 4 keywords, important terms and concepts from the markdown section title then from the content
			- Main topic: the primary subject (use the markdown section title)
			- Category: type of content
			`,
				piece,
			)

			metadata, _, err := metadataExtractorAgent.GenerateStructuredData([]messages.Message{
				{Role: roles.User, Content: extractionPrompt},
			})
			if err != nil {
				display.Errorf("❌ Error extracting keywords from section %d: %v", idx, err)
				// Continue with embedding even if keyword extraction fails
			} else {
				display.Infof("🏷️  Keywords: %v", metadata.Keywords)
				display.Infof("📌 Topic: %s | Category: %s",
					metadata.MainTopic, metadata.Category)

				// Enrich the chunk with metadata
				enrichedPiece := fmt.Sprintf("[METADATA]\nKeywords: %v\nTopic: %s\nCategory: %s\n\nContent:\n%s",
					metadata.Keywords, metadata.MainTopic, metadata.Category, piece,
				)

				piece = enrichedPiece
			}
			
			// === SAVE EMBEDDING FOR SECTION ===
			err = ragAgent.SaveEmbedding(piece)
			if err != nil {
				display.Errorf("❌ Error embedding section %d: %v", idx, err)
			} else {
				display.Infof("✅ Indexed section %d/%d", idx+1, len(contentPieces))
				fmt.Println(piece)
			}

		}

		// === PERSIST STORE TO DISK ===
		err = ragAgent.PersistStore(storePathFile)
		if err != nil {
			display.Errorf("❌ Error persisting store: %v", err)
			return nil, err
		}
		display.Infof("💾 RAG store saved to %s", storePathFile)
	}

	return ragAgent, nil
}
