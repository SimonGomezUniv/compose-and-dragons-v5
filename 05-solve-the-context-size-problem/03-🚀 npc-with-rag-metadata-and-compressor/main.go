package main

import (
	"context"
	"strings"

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

	// NOTE: You can change the query to generate different NPCs
	query := "Generate a male dwarf warrior character."
	sheetFilePath := "./sheets/male-dwarf-warrior.md"

	engineURL := "http://localhost:12434/engines/llama.cpp/v1"
	generatorModelID := "huggingface.co/unsloth/nvidia-nemotron-3-nano-4b-gguf:Q4_K_M"
	npcModelID := "ai/qwen2.5:1.5B-F16"
	ragEmbeddingModel := "ai/embeddinggemma:latest"
	compressorModelId := "ai/qwen2.5:0.5B-F16"
	metadataModel := "huggingface.co/menlo/jan-nano-gguf:q4_k_m"



	// === CHECK FOR EXISTING CHARACTER SHEETS ===
	// test if the file sheetFilePath already exists
	if files.FileExists(sheetFilePath) {
		display.Infof("📂 Found existing character sheet: %s", sheetFilePath)
	} else {
		display.Infof("🔎 No existing character sheet found at: %s", sheetFilePath)
		// Generate new character
		display.Infof("🎲 Generating new character...")
		err := generateNewCharacter(ctx, engineURL, generatorModelID, query, sheetFilePath)
		if err != nil {
			display.Errorf("❌ Error generating new character: %v", err)
			return
		}
		display.Infof("✅ New character generated successfully.")
	}

	// === CREATE METADATA EXTRACTOR AGENT ===
	metadataExtractorAgent , err := getMetadataExtractorAgent(ctx, engineURL, metadataModel)
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
	errChat := startInteractiveChat(ctx, engineURL, npcModelID, sheetJsonFilePath, ragAgent, compressorAgent)
	if errChat != nil {
		display.Errorf("❌ Error during interactive chat: %v", errChat)
		return
	}
}
