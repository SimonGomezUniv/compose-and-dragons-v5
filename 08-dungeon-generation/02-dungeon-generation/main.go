package main

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"codeberg.org/rpg/dungeon/dungeon"
	"github.com/snipwise/nova/nova-sdk/agents"
	"github.com/snipwise/nova/nova-sdk/agents/structured"
	"github.com/snipwise/nova/nova-sdk/messages"
	"github.com/snipwise/nova/nova-sdk/messages/roles"
	"github.com/snipwise/nova/nova-sdk/models"
	"github.com/snipwise/nova/nova-sdk/toolbox/conversion"
	"github.com/snipwise/nova/nova-sdk/toolbox/env"
	"github.com/snipwise/nova/nova-sdk/toolbox/files"
	"github.com/snipwise/nova/nova-sdk/toolbox/logger"
	"github.com/snipwise/nova/nova-sdk/ui/display"
	"github.com/snipwise/nova/nova-sdk/ui/spinner"
)

/*
This program is used to generate a room in a dungeon, with its name, description,
and the items it contains (treasures, potions, ...).
A room can be empty, contain multiple items, or just one...
*/

type Item struct {
	Type dungeon.Item `json:"type"`
}

type GeneratedRoom struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	ShortDescription string `json:"short_description"`
	Items            []Item `json:"items"`
}

// interpolateEnvVars replaces environment variable references in the text
// Supports both ${VAR_NAME} and $VAR_NAME formats
func interpolateEnvVars(text string) string {
	// Replace ${VAR_NAME} format
	re1 := regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`)
	text = re1.ReplaceAllStringFunc(text, func(match string) string {
		varName := re1.FindStringSubmatch(match)[1]
		if value := os.Getenv(varName); value != "" {
			return value
		}
		return match // Keep original if env var not found
	})

	// Replace $VAR_NAME format (not followed by { and word boundary)
	re2 := regexp.MustCompile(`\$([A-Z_][A-Z0-9_]*)(?:\b|$)`)
	text = re2.ReplaceAllStringFunc(text, func(match string) string {
		varName := re2.FindStringSubmatch(match)[1]
		if value := os.Getenv(varName); value != "" {
			return value
		}
		return match // Keep original if env var not found
	})

	return text
}

func main() {

	log := logger.NewConsoleLogger(logger.LevelDebug)

	generatingDungeonSpinner := spinner.
		NewWithColor("").
		SetSuffix("generating dungeon...").
		SetFrames(spinner.FramesPulsingStar).
		SetSuffixColor(spinner.ColorBrightBlue).
		SetFrameColor(spinner.ColorBrightBlue)

	generatingRoomsSpinner := spinner.
		NewWithColor("").
		SetSuffix("generating rooms...").
		SetFrames(spinner.FramesPulsingStar).
		SetSuffixColor(spinner.ColorBrightGreen).
		SetFrameColor(spinner.ColorBrightGreen)

	generatingDungeonSpinner.Start()

	dng := dungeon.NewDungeon(dungeon.DungeonColors{
		StartColor: dungeon.ColorBrightYellow,
		EndColor:   dungeon.ColorBrightGreen,
		EmptyColor: dungeon.ColorGray,
	})

	// Get number of rooms from environment variable or default to 8

	nbDungeonRooms := env.GetEnvIntOrDefault("NUMBER_OF_ROOMS", 8)

	dng.Generate(nbDungeonRooms)

	dungeonMap := dng.GetBWDetailedGridCompact()

	dungeonRoomsConnections := dng.GetDungeonDescriptionText()

	rooms := dng.GetRooms()

	generatingDungeonSpinner.Stop()

	modelConfig := models.Config{
		Name:        env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m"),
		Temperature: models.Float64(1.0),
		TopP:        models.Float64(0.9),
		TopK:        models.Int(40),
	}

	userInstructions, err := files.ReadTextFile("./user.instructions.md")
	if err != nil {
		log.Error("Error reading user instructions file: %v", err)
	} else {
		log.Info("User instructions loaded from file.")
		// Interpolate environment variables in user instructions
		userInstructions = interpolateEnvVars(userInstructions)
		fmt.Println("User Instructions:")
		fmt.Println(userInstructions)
	}

	systemInstructions, err := files.ReadTextFile("./system.instructions.md")
	if err != nil {
		log.Error("Error reading system instructions file: %v", err)
	} else {
		log.Info("System instructions loaded from file.")
	}

	ctx := context.Background()
	dungeonGeneratorAgent, err := structured.NewAgent[[]GeneratedRoom](
		ctx,
		agents.Config{
			EngineURL:          env.GetEnvOrDefault("ENGINE_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1"),
			SystemInstructions: systemInstructions,
			KeepConversationHistory: true,
		},
		modelConfig,
	)
	if err != nil {
		panic(err)
	}

	dungeonGeneratorAgent.AddMessage(
		roles.System,
		`DUNGEON MAP:\n`+dungeonMap+`\n`,
	)
	dungeonGeneratorAgent.AddMessage(
		roles.System,
		`DUNGEON ROOMS CONNECTIONS LIST:\n`+dungeonRoomsConnections+`\n`,
	)

	generatingRoomsSpinner.Start()

	generatedRooms, finishReason, err := dungeonGeneratorAgent.GenerateStructuredData(
		[]messages.Message{
			{
				Role:    roles.User,
				Content: userInstructions,
			},
		},
	)
	_ = finishReason

	if err != nil {
		log.Error("Error generating rooms %v", err)
		return
	}

	for idx, generatedRoom := range *generatedRooms {
		generatingRoomsSpinner.SetSuffix("generating room #" + conversion.IntToString(idx) + "...")

		room := rooms[idx]
		room.SetName(generatedRoom.Name)
		room.SetDescription(generatedRoom.Description)

		display.Title("Room Generated")
		display.KeyValue("Name", generatedRoom.Name)
		display.KeyValue("Description", generatedRoom.Description)
		display.KeyValue("Short Description", generatedRoom.ShortDescription)
		display.Separator()
		if len(generatedRoom.Items) == 0 {
			display.Info("No items in this room.")
		} else {
			for itemIdx, item := range generatedRoom.Items {
				display.Println("  - " + string(item.Type))

				switch item.Type {
				case dungeon.Goblin:

					goblin := NewGoblin(
						dng, dungeon.ItemArgs{
							RoomNumber: room.ID,
							Id:         "goblin_" + conversion.IntToString(itemIdx),
							ForeColor:  dungeon.ColorBrightGreen,
							Shape:      dungeon.Shape3x3,
						},
					)
					_ = goblin

				case dungeon.Skeleton:

					skeleton := NewSkeleton(
						dng, dungeon.ItemArgs{
							RoomNumber: room.ID,
							Id:         "skeleton_" + conversion.IntToString(itemIdx),
							ForeColor:  dungeon.ColorBrightWhite,
							Shape:      dungeon.Shape3x3,
						},
					)
					_ = skeleton

				case dungeon.Vampire:

					vampire := NewVampire(
						dng, dungeon.ItemArgs{
							RoomNumber: room.ID,
							Id:         "vampire_" + conversion.IntToString(itemIdx),
							ForeColor:  dungeon.ColorBrightMagenta,
							Shape:      dungeon.Shape3x3,
						},
					)
					_ = vampire

				case dungeon.Sphinx:

					sphinx := NewSphinx(
						dng, dungeon.ItemArgs{
							RoomNumber: room.ID,
							Id:         "sphinx_" + conversion.IntToString(itemIdx),
							ForeColor:  dungeon.ColorBrightCyan,
							Shape:      dungeon.Shape3x3,
							
						},
					)
					_ = sphinx

				case dungeon.Treasure:

					gold := NewGoldCoins(dng, dungeon.ItemArgs{
						Id:         "gold_" + conversion.IntToString(itemIdx),
						RoomNumber: room.ID,
						ForeColor:  dungeon.ColorBrightYellow,
						X:          2,
						Y:          2,
					})
					_ = gold

				case dungeon.Potion:

					potion := NewMagicPotion(dng, dungeon.ItemArgs{
						Id:         "potion_" + conversion.IntToString(itemIdx),
						RoomNumber: room.ID,
						ForeColor:  dungeon.ColorBrightRed,
						X:          12,
						Y:          2,
					})
					_ = potion

				}
			}
		}
		display.Separator()
	}

	generatingRoomsSpinner.Stop()

	playerId := env.GetEnvOrDefault("PLAYER_ID", "hero_1")
	playerName := env.GetEnvOrDefault("PLAYER_NAME", "Aragorn")
	playerKind := env.GetEnvOrDefault("PLAYER_KIND", "Human")
	playerClass := env.GetEnvOrDefault("PLAYER_CLASS", "Warrior")

	player := NewPlayer(dng, dungeon.ItemArgs{
		Id:         playerId,
		RoomNumber: 1,
		ForeColor:  dungeon.ColorBrightBlue,
	}, playerName, playerKind, playerClass)

	// Add boss sphinx to the last room (exit room)
	_ = NewSphinx(dng, dungeon.ItemArgs{
		Id:         "boss_sphinx",
		RoomNumber: nbDungeonRooms,
		ForeColor:  dungeon.ColorBrightCyan,
		Shape:      dungeon.Shape3x3,
	})

	// Find rooms without monsters to add NPCs
	roomsWithoutMonsters := []int{}
	for _, room := range rooms {
		hasMonster := false
		for _, char := range room.Chars {
			kind := string(char.Kind)
			if kind == "goblin" || kind == "skeleton" || kind == "vampire" || kind == "sphinx" {
				hasMonster = true
				break
			}
		}
		if !hasMonster && room.ID != 1 && room.ID != nbDungeonRooms {
			roomsWithoutMonsters = append(roomsWithoutMonsters, room.ID)
		}
	}

	// Add NPCs to safe rooms

	dwarfId := env.GetEnvOrDefault("DWARF_ID", "Thorin")
	elfId := env.GetEnvOrDefault("ELF_ID", "Elrond")
	humanId := env.GetEnvOrDefault("HUMAN_ID", "Beren")


	if len(roomsWithoutMonsters) > 0 {
		// Add Dwarf to the first safe room
		dwarf := NewDwarf(dng, dungeon.ItemArgs{
			Id:         dwarfId,
			RoomNumber: roomsWithoutMonsters[0],
			ForeColor:  dungeon.ColorBrightYellow,
			Shape:      dungeon.Shape3x3,
		})
		_ = dwarf
		log.Info("Added Dwarf '%s' to room #%d", dwarfId, roomsWithoutMonsters[0])
	}

	if len(roomsWithoutMonsters) > 1 {
		// Add Elf to the second safe room
		elf := NewElf(dng, dungeon.ItemArgs{
			Id:         elfId,
			RoomNumber: roomsWithoutMonsters[1],
			ForeColor:  dungeon.ColorBrightGreen,
			Shape:      dungeon.Shape3x3,
		})
		_ = elf
		log.Info("Added Elf '%s' to room #%d", elfId, roomsWithoutMonsters[1])
	}

	if len(roomsWithoutMonsters) > 2 {
		// Add Human to the third safe room
		human := NewHuman(dng, dungeon.ItemArgs{
			Id:         humanId,
			RoomNumber: roomsWithoutMonsters[2],
			ForeColor:  dungeon.ColorBrightBlue,
			Shape:      dungeon.Shape3x3,
		})
		_ = human
		log.Info("Added Human '%s' to room #%d", humanId, roomsWithoutMonsters[2])
	}

	// Display the dungeon
	fmt.Println(dng.GetDetailedGridCompact())

	// Save dungeon to file
	errSave := dng.SaveToFile("./data/dungeon_generated.json")
	if errSave != nil {
		log.Error("Error saving dungeon to file: %v", err)
	}
	log.Info("Dungeon saved to data/dungeon_generated.json")

	// Create and save metadata
	metadata := &DungeonMetadata{
		Entities: []EntityMetadata{},
	}

	// Add player metadata
	metadata.AddEntityMetadata(EntityMetadata{
		ID:       player.Id,
		Type:     "player",
		Name:     player.Name,
		Race:     player.Race,
		Class:    player.Class,
		Health:   player.Health,
		Strength: player.Strength,
	})

	// Collect all monsters from the dungeon
	rooms = dng.GetRooms()
	for _, room := range rooms {
		for _, char := range room.Chars {
			entityType := string(char.Kind)
			entityID := char.ObjectID

			// Only save metadata for entities with health/strength
			switch entityType {
			case "goblin", "skeleton", "vampire", "sphinx":
				// Default values for monsters (they all have the same stats)
				metadata.AddEntityMetadata(EntityMetadata{
					ID:       entityID,
					Type:     entityType,
					Health:   10,
					Strength: 5,
				})
			case "elf":
				metadata.AddEntityMetadata(EntityMetadata{
					ID:       entityID,
					Type:     entityType,
					Health:   15,
					Strength: 7,
				})
			case "dwarf":
				metadata.AddEntityMetadata(EntityMetadata{
					ID:       entityID,
					Type:     entityType,
					Health:   20,
					Strength: 10,
				})
			case "human":
				metadata.AddEntityMetadata(EntityMetadata{
					ID:       entityID,
					Type:     entityType,
					Health:   15,
					Strength: 8,
				})
			}
		}
	}

	errMetadata := SaveMetadata(metadata, "./data/dungeon_metadata.json")
	if errMetadata != nil {
		log.Error("Error saving metadata to file: %v", errMetadata)
	}
	log.Info("Metadata saved to data/dungeon_metadata.json")

}
