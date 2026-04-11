package main

import (
	"context"
	"fmt"
	"os"

	"github.com/snipwise/nova/nova-sdk/agents"
	"github.com/snipwise/nova/nova-sdk/agents/structured"
	"github.com/snipwise/nova/nova-sdk/messages"
	"github.com/snipwise/nova/nova-sdk/messages/roles"
	"github.com/snipwise/nova/nova-sdk/models"
	"github.com/snipwise/nova/nova-sdk/toolbox/files"
	"github.com/snipwise/nova/nova-sdk/ui/display"
)

// === NPCCharacter represents a generated D&D character with structured output ===
type NPCCharacter struct {
	FirstName  string `json:"firstName" jsonschema:"required,description=Character's first name following race conventions"`
	FamilyName string `json:"familyName" jsonschema:"required,description=Clan/house/family name following race conventions"`
	Race       string `json:"race" jsonschema:"required,description=Character race (Dwarf/Elf/Human)"`
	Class      string `json:"class" jsonschema:"required,description=D&D character class (Warrior/Mage/Ranger/Cleric/Rogue/Paladin/etc)"`
	Gender     string `json:"gender" jsonschema:"required,description=Character gender (male/female)"`
}

func main() {

	ctx := context.Background()

	engineURL := "http://localhost:12434/engines/llama.cpp/v1"
	modelID := "huggingface.co/qwen/qwen2.5-0.5b-instruct-gguf:Q4_K_M"
	//modelID := "huggingface.co/unsloth/qwen3.5-0.8b-gguf:Q4_K_M"

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

	// fmt.Println(strings.Repeat("=", 50))
	// fmt.Println(systemInstructions)
	// fmt.Println(strings.Repeat("=", 50))

	// === CREATE D&D NPC GENERATOR AGENT ===
	agent, err := structured.NewAgent[NPCCharacter](
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

	// === EXAMPLES OF USAGE ===
	testCases := []string{
		"Generate a dwarf character",
		"Create a female elf ranger",
		"Generate a male human paladin",
	}

	display.Separator()
	display.Colorf(display.ColorBrightBlue, "🎲 Starting D&D NPC Character Generation Tests...\n")
	display.Colorf(display.ColorBrightBlue, "🧠 Using Model: %s\n", modelID)
	display.Separator()

	for i, query := range testCases {

		display.Colorf(display.ColorBrightCyan, "📝 Request %d: %s\n", i+1, query)
		display.Colorln(display.ColorBrightPurple, "🔄 Generating NPC...")

		// === Generate structured output ===
		npc, _, err := agent.GenerateStructuredData([]messages.Message{
			{Role: roles.User, Content: query},
		})
		if err != nil {
			fmt.Printf("❌ Error generating NPC: %v\n\n", err)
			continue
		}

		// Display generated NPC
		display.Title("🧙 Generated NPC Summary:")
		display.Table("Name", npc.FirstName+" "+npc.FamilyName)
		display.Table("Race", npc.Race)
		display.Table("Class", npc.Class)
		display.Table("Gender", npc.Gender)
		display.Separator()

	}
	
	display.Styledln("✨ Generation complete!",display.ColorBrightWhite, display.BgBlue)
}
