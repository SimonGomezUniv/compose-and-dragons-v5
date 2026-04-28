package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
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

// NPCCharacter represents a generated D&D character with structured output
type NPCCharacter struct {
	FirstName  string `json:"firstName" jsonschema:"required,description=Character's first name following race conventions"`
	FamilyName string `json:"familyName" jsonschema:"required,description=Clan/house/family name following race conventions"`
	Race       string `json:"race" jsonschema:"required,description=Character race (Dwarf/Elf/Human)"`
	Class      string `json:"class" jsonschema:"required,description=D&D character class (Warrior/Mage/Ranger/Cleric/Rogue/Paladin/etc)"`
	Gender     string `json:"gender" jsonschema:"required,description=Character gender (male/female)"`
}

func main() {
	ctx := context.Background()

	// NOTE: You can change the query to generate different NPCs
	// query := "Generate a male dwarf warrior"
	query := "Generate a female elf warrior"

	engineURL := "http://localhost:11434/v1"
	//modelID := "huggingface.co/unsloth/nvidia-nemotron-3-nano-4b-gguf:Q4_K_M"
	//modelID := "ai/qwen2.5:3B-F16"
	//modelID := "ai/qwen2.5:1.5B-F16"
	modelID := "qwen2:0.5b"

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
		display.Errorf("❌ Error creating agent: %v\n", err)
		os.Exit(1)
	}

	display.Separator()
	display.Colorf(display.ColorBrightBlue, "🎲 Starting D&D NPC Character Generation Tests...\n")
	display.Colorf(display.ColorBrightBlue, "🧠 Using Model: %s\n", modelID)
	display.Separator()

	display.Colorf(display.ColorBrightCyan, "📝 Request: %s\n", query)
	display.Colorln(display.ColorBrightPurple, "🔄 Generating NPC...")

	// === Generate structured output ===
	npc, _, err := npcGeneratorAgent.GenerateStructuredData([]messages.Message{
		{Role: roles.User, Content: query},
	})
	if err != nil {
		fmt.Printf("❌ Error generating NPC: %v\n\n", err)
		panic(err)
	}

	// Display generated NPC
	display.Title("🧙 Generated NPC Summary:")
	display.Table("Name", npc.FirstName+" "+npc.FamilyName)
	display.Table("Race", npc.Race)
	display.Table("Class", npc.Class)
	display.Table("Gender", npc.Gender)
	display.Separator()

	// === CREATE STORY GENERATOR AGENT ===
	display.Println("📖 Creating character sheet with streaming...")
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
			Temperature: models.Float64(0.8), // More creativity for storytelling
			//MaxTokens:   models.Int(4096),     // Increased to ensure all sections are generated
			TopP: models.Float64(0.95), // Diverse vocabulary
		},
	)
	if err != nil {
		display.Errorf("❌ Error creating story generator agent: %v\n", err)
		os.Exit(1)
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
	display.Println("🤖 Generating character sheet...")

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
		display.Errorf("\n❌ Error generating character sheet: %v\n", err)
		os.Exit(1)
	}

	display.NewLine()
	display.Separator()

	// === SAVE CHARACTER SHEET TO FILE ===
	// === Generate filename from character name ===
	// Convert "Eldorin Shadowleaf" -> "character-sheet-eldorin-shadowleaf.md"
	filename := generateFilename(npc.FirstName, npc.FamilyName)

	err = files.WriteTextFile("./sheets/"+filename, result.Response)

	if err != nil {
		display.Errorf("❌ Error saving character sheet: %v\n", err)
		os.Exit(1)
	}

	display.KeyValue("✅ Character sheet saved to", filename)
	display.KeyValue("🏁 Finish reason", result.FinishReason)
	display.Separator()
	display.Styledln("✨ Generation complete!", display.ColorBrightWhite, display.BgBlue)

}

// generateFilename creates a filename from character name
// Example: "Eldorin Shadowleaf" -> "character-sheet-eldorin-shadowleaf.md"
func generateFilename(firstName, familyName string) string {
	// Combine names
	fullName := fmt.Sprintf("%s %s", firstName, familyName)

	// Convert to lowercase
	name := strings.ToLower(fullName)

	// Replace spaces with hyphens
	name = strings.ReplaceAll(name, " ", "-")

	// Remove any non-alphanumeric characters except hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	name = reg.ReplaceAllString(name, "")

	// Create final filename
	return fmt.Sprintf("character-sheet-%s.md", name)
}
