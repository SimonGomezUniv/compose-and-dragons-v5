package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/snipwise/nova/nova-sdk/agents"
	"github.com/snipwise/nova/nova-sdk/agents/chat"
	"github.com/snipwise/nova/nova-sdk/agents/structured"
	"github.com/snipwise/nova/nova-sdk/messages"
	"github.com/snipwise/nova/nova-sdk/messages/roles"
	"github.com/snipwise/nova/nova-sdk/models"
	"github.com/snipwise/nova/nova-sdk/toolbox/files"
	"github.com/snipwise/nova/nova-sdk/ui/display"
)

// generateNewCharacter creates a new NPC character and its sheet
func generateNewCharacter(ctx context.Context, engineURL, modelID, query, sheetFilePath string) error {
	// === NAMING RULES DATABASE ===
	// This resource guides the AI in generating authentic names
	namingRules, err := files.ReadTextFile("./dnd.naming.rules.md")
	if err != nil {
		display.Errorf("❌ Error reading naming rules file: %v\n", err)
		os.Exit(1)
	}

	// === SYSTEM INSTRUCTIONS TEMPLATE ===
	systemInstructionsTemplate, err := files.ReadTextFile("./dnd.system.instructions.md")
	if err != nil {
		display.Errorf("❌ Error reading system instructions template file: %v\n", err)
		os.Exit(1)
	}

	// === Inject naming rules into system instructions ===
	systemInstructions := fmt.Sprintf(systemInstructionsTemplate, namingRules)

	// === CREATE D&D NPC GENERATOR AGENT ===
	npcGeneratorAgent, err := structured.NewAgent[NPCCharacter](
		ctx,
		agents.Config{
			Name:               "dnd-npc-generator",
			EngineURL:          engineURL,
			SystemInstructions: systemInstructions,
		},
		models.Config{
			Name: modelID,
			// Some creativity for name generation
			Temperature: models.Float64(0.7),
			TopP:        models.Float64(0.9),
			TopK:        models.Int64(40),
		},
	)
	if err != nil {
		display.Errorf("❌ Error creating agent: %v", err)
		return err
	}

	// === GENERATE NPC CHARACTER ===
	display.Infof("🎲 D&D NPC Character Generator")
	display.Separator()
	display.Infof("📝 Request %s:", query)
	display.Println("🔄 Generating NPC...")

	// Generate structured output
	npc, _, err := npcGeneratorAgent.GenerateStructuredData([]messages.Message{
		{Role: roles.User, Content: query},
	})
	if err != nil {
		display.Errorf("❌ Error generating NPC: %v", err)
		return err
	}

	// Display NPC summary
	display.Infof("🧙 Generated NPC Summary:")
	display.Table("Name", npc.FirstName+" "+npc.FamilyName)
	display.Table("Race", npc.Race)
	display.Table("Class", npc.Class)
	display.Table("Gender", npc.Gender)
	display.Table("SecretWord", npc.SecretWord)
	display.Separator()

	// === CREATE STORY GENERATOR AGENT ===
	display.Infof("📖 Creating character sheet for %s %s...", npc.FirstName, npc.FamilyName)
	display.Separator()

	storySystemInstructions, err := files.ReadTextFile("./dnd.story.system.instructions.md")
	if err != nil {
		display.Errorf("❌ Error reading story system instructions file: %v\n", err)
		os.Exit(1)
	}

	// === CREATE NPC STORY GENERATOR AGENT ===
	npcStoryGeneratorAgent, err := chat.NewAgent(
		ctx,
		agents.Config{
			Name:               "npc-story-generator",
			EngineURL:          engineURL,
			SystemInstructions: storySystemInstructions,
		},
		models.Config{
			Name:        modelID,
			Temperature: models.Float64(0.8),  // More creativity for storytelling
			MaxTokens:   models.Int(4096),     // Increased to ensure all sections are generated
			TopP:        models.Float64(0.95), // Diverse vocabulary
		},
	)
	if err != nil {
		display.Errorf("❌ Error creating story generator agent: %v", err)
		return err
	}

	// === Prepare prompt for character sheet generation ===
	characterPrompt := fmt.Sprintf(`Generate a complete character sheet for:
		Name: %s %s
		Race: %s
		Class: %s
		Gender: %s
		`, npc.FirstName, npc.FamilyName, npc.Race, npc.Class, npc.Gender,
	)
	
	// === GENERATE CHARACTER SHEET WITH STREAMING ===
	display.Infof("🤖 Generating character sheet...")

	// Stream the character sheet generation
	var fullResponse strings.Builder
	result, err := npcStoryGeneratorAgent.GenerateStreamCompletion(
		[]messages.Message{
			{Role: roles.User, Content: characterPrompt},
		},
		func(chunk string, finishReason string) error {
			// Print chunk to console
			fmt.Print(chunk)
			// Accumulate full response for saving
			fullResponse.WriteString(chunk)
			return nil
		},
	)
	if err != nil {
		display.Errorf("❌ Error generating character sheet: %v", err)
		return err
	}

	display.NewLine()

	errSave := saveNPCSheetToFile(sheetFilePath, result.Response, npc)
	if errSave != nil {
		display.Errorf("❌ Error saving character sheet and NPC data: %v", errSave)
		return errSave
	}

	display.Infof("✅ Character sheet and NPC data saved successfully.")
	display.KeyValue("Character Sheet Path", sheetFilePath)
	display.KeyValue("NPC JSON Path", strings.TrimSuffix(sheetFilePath, ".md")+".json")
	display.KeyValue("Context Size", fmt.Sprintf("%d characters", npcStoryGeneratorAgent.GetContextSize()))
	display.KeyValue("Finish reason", result.FinishReason)

	display.Println("✨ Generation complete!")

	return nil
}
