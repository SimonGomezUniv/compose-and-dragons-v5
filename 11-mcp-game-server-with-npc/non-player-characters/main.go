package main

import (
	"context"
	"errors"
	"strings"

	"github.com/snipwise/nova/nova-sdk/toolbox/env"
	"github.com/snipwise/nova/nova-sdk/toolbox/files"
	"github.com/snipwise/nova/nova-sdk/ui/display"
)

// NPCCharacter represents a generated D&D character with structured output
type NPCCharacter struct {
	FirstName  string `json:"firstName" jsonschema:"required,description=Character's first name following race conventions"`
	FamilyName string `json:"familyName" jsonschema:"required,description=Clan/house/family name following race conventions"`
	Race       string `json:"race" jsonschema:"required,description=Character race (Dwarf/Elf/Human)"`
	Class      string `json:"class" jsonschema:"required,description=D&D character class (Warrior/Mage/Ranger/Cleric/Rogue/Paladin/etc)"`
	Gender     string `json:"gender" jsonschema:"required,description=Character gender (male/female)"`
	SecretWord string `json:"secretWord" jsonschema:"required,description=A unique secret word associated with the character"`
}

func main() {
	ctx := context.Background()

	sheetsPath := env.GetEnvOrDefault("SHEETS_PATH", "./sheets")
	sheetFilePath := sheetsPath + "/" + env.GetEnvOrDefault("SHEET_FILE_NAME", "female-klingon-warrior.md")

	//engineURL := "http://localhost:11434/v1"
	engineURL := env.GetEnvOrDefault("ENGINE_URL", "http://localhost:11434/v1")
	
	//npcModelID := "qwen2:0.5b"
	npcModelID := env.GetEnvOrDefault("NPC_MODEL_ID", "qwen2:0.5b")

	//ragEmbeddingModel := "qwen2:0.5b"
	ragEmbeddingModel := env.GetEnvOrDefault("RAG_EMBEDDING_MODEL_ID", "qwen2:0.5b")

	//compressorModelId := "qwen2:0.5b"
	compressorModelId := env.GetEnvOrDefault("COMPRESSOR_MODEL_ID", "qwen2:0.5b")	

	//metadataModel := "qwen2:0.5b"
	metadataModel := env.GetEnvOrDefault("METADATA_MODEL_ID", "qwen2:0.5b")

	// === CHECK FOR EXISTING CHARACTER SHEETS ===
	// test if the file sheetFilePath already exists
	if files.FileExists(sheetFilePath) {
		display.Infof("📂 Found existing character sheet: %s", sheetFilePath)
	} else {
		display.Infof("🔎 No existing character sheet found at: %s", sheetFilePath)
		err := errors.New("❌ Character sheet file not found, please create one before starting the NPC server")
		display.Errorf("%v", err)
		return
	}

	// === CREATE METADATA EXTRACTOR AGENT ===
	metadataExtractorAgent, err := getMetadataExtractorAgent(ctx, engineURL, metadataModel)
	if err != nil {
		display.Errorf("❌ Error creating metadata extractor agent: %v", err)
		return
	}

	// === CREATE/LOAD RAG AGENT ===
	ragAgent, err := getRagAgent(ctx, engineURL, ragEmbeddingModel, sheetFilePath, metadataExtractorAgent)
	if err != nil {
		display.Errorf("❌ Error creating/loading RAG agent: %v", err)
		return
	}

	// === CREATE COMPRESSOR AGENT ===
	compressorAgent, err := getCompressorAgent(ctx, engineURL, compressorModelId)
	if err != nil {
		display.Errorf("❌ Error creating compressor agent: %v", err)
		return
	}

	sheetJsonFilePath := strings.TrimSuffix(sheetFilePath, ".md") + ".json"
	// === INTERACTIVE NPC CHAT ===
	startNPCServer(ctx, engineURL, npcModelID, sheetJsonFilePath, ragAgent, compressorAgent)
}
