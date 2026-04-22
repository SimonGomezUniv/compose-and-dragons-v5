package main

import (
	"context"

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
	query := "Generate a female elf sorcerer"
	//query := "Generate a male dwarf warrior"

	//sheetFilePath := "./sheets/character-sheet-thorin-stonefist.md"
	sheetFilePath := "./sheets/female-elf-sorcerer.md"

	// === CONFIGURATION ===
	engineURL := "http://localhost:12434/engines/llama.cpp/v1"
	//generatorModelID := "huggingface.co/unsloth/nvidia-nemotron-3-nano-4b-gguf:Q4_K_M"

	generatorModelID := "huggingface.co/menlo/jan-nano-128k-gguf:Q4_K_M"

	// NOTE: Using a smaller model for the interactive chat to reduce latency
	npcModelID := "ai/qwen2.5:0.5B-F16"

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

	// === INTERACTIVE NPC CHAT ===
	err := startInteractiveChat(ctx, engineURL, npcModelID, sheetFilePath)
	if err != nil {
		display.Errorf("❌ Error starting interactive chat: %v", err)
		return
	}
}
