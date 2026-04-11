package main

import (
	"context"
	"fmt"

	"github.com/snipwise/nova/nova-sdk/agents"
	"github.com/snipwise/nova/nova-sdk/agents/compressor"
	"github.com/snipwise/nova/nova-sdk/agents/rag"
	"github.com/snipwise/nova/nova-sdk/agents/server"
	"github.com/snipwise/nova/nova-sdk/models"
	"github.com/snipwise/nova/nova-sdk/toolbox/env"
	"github.com/snipwise/nova/nova-sdk/toolbox/files"
	"github.com/snipwise/nova/nova-sdk/ui/display"
)

// startNPCServer starts an interactive roleplay session with the NPC
func startNPCServer(ctx context.Context, engineURL, npcAgentModelID, sheetJsonFilePath string, ragAgent *rag.Agent, compressorAgent *compressor.Agent) error {

	npc, err := loadNPCSheetFromJsonFile(sheetJsonFilePath)
	if err != nil {
		display.Errorf("failed to load NPC character sheet: %v", err)
		return err
	}

	// === INTERACTIVE NPC CHAT ===
	display.Infof("💬 Starting interactive chat with NPC...")
	display.Separator()

	// Create system instructions for the roleplay agent
	// dnd.chat.system.instructions.md

	docsPath := env.GetEnvOrDefault("DOCS_PATH", "./docs")

	roleplaySystemInstructionsTemplate, err := files.ReadTextFile(docsPath + "/dnd.chat.system.instructions.md")
	if err != nil {
		display.Errorf("failed to read roleplay system instructions file: %v", err)
		return err
	}

	// Create system instructions for the roleplay agent
	roleplaySystemInstructions := fmt.Sprintf(roleplaySystemInstructionsTemplate,
		npc.FirstName,
		npc.FamilyName,
		npc.Gender,
		npc.Race,
		npc.Class,
		npc.SecretWord,
		npc.FirstName, // Always speak in first person as %s

	)
	// NOTE: this time we DO NOT inject the character sheet content directly
	// into the system instructions, as we will use RAG to provide context

	// === CREATE THE ROLEPLAY CHAT AGENT ===
	npcServerAgent, err := server.NewAgent(
		ctx,
		agents.Config{
			Name:                    "npc-roleplay-agent",
			EngineURL:               engineURL,
			SystemInstructions:      roleplaySystemInstructions,
			KeepConversationHistory: true, // IMPORTANT: Keep history for better context
		},
		models.Config{
			Name:        npcAgentModelID,
			Temperature: models.Float64(0.9), // High creativity for roleplay
			TopP:        models.Float64(0.95),
		},
		server.WithPort(8080),
		server.WithRagAgentAndSimilarityConfig(ragAgent, 0.4, 7),
		server.WithCompressorAgentAndContextSize(compressorAgent, 80000),
	)
	if err != nil {
		display.Errorf("failed to create NPC agent: %v", err)
		return err
	}

	// npcServerAgent.SetRagAgent(ragAgent)
	// npcServerAgent.SetSimilarityLimit(0.4)
	// npcServerAgent.SetMaxSimilarities(7)

	// npcServerAgent.SetCompressorAgent(compressorAgent)
	// npcServerAgent.SetContextSizeLimit(8000)

	initialContextSize := npcServerAgent.GetContextSize()
	display.Colorf(display.ColorBrightRed, "Initial Context Size: %d characters\n", initialContextSize)
	display.Separator()

	display.Colorf(display.ColorCyan, "🚀 Server starting on http://localhost%s\n", npcServerAgent.GetPort())

	// Start the server
	if err := npcServerAgent.StartServer(); err != nil {
		display.Errorf("failed to start NPC server: %v", err)
		return err
	}
	return nil
}
