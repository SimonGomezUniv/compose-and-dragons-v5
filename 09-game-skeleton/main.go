package main

import (
	"fmt"
	"log"
	"strings"

	"codeberg.org/rpg/dungeon/dungeon"
)

func main() {
	// Find the most recent save files
	dungeonFile := FindMostRecentSave("./data", "dungeon_generated.json")
	metadataFile := FindMostRecentSave("./data", "dungeon_metadata.json")

	// Check if we're loading a saved game
	if dungeonFile != "./data/dungeon_generated.json" {
		fmt.Printf("📂 Loading saved game: %s\n", dungeonFile)
	}

	// Load the dungeon from the JSON file
	dng, err := dungeon.LoadFromFile(dungeonFile)
	if err != nil {
		log.Fatalf("Error loading dungeon from file: %v", err)
	}

	fmt.Println("✨ Dungeon loaded successfully!")
	fmt.Println()

	// Load metadata
	metadata, err := LoadMetadata(metadataFile)
	if err != nil {
		log.Printf("⚠️  Warning: Could not load metadata file: %v", err)
		metadata = &DungeonMetadata{Entities: []EntityMetadata{}}
	} else {
		fmt.Printf("✨ Metadata loaded successfully! (%d entities)\n", len(metadata.Entities))
		fmt.Println()
	}

	// Initialize game state
	gameState := NewGameState(dng)

	// Load all collectibles and monsters from the dungeon
	loadEntitiesFromDungeon(gameState, metadata)

	// Find the player in the dungeon
	player := findPlayer(gameState, metadata)
	if player == nil {
		log.Fatal("Player not found in the dungeon!")
	}
	gameState.Player = player

	// Load riddle state from metadata
	gameState.RiddleSolved = metadata.RiddleSolved

	// Display welcome message
	fmt.Println("🎮 Welcome to the Dungeon Adventure!")
	fmt.Println("====================================")
	fmt.Println("Your quest: Explore the dungeon, collect treasures, and defeat monsters!")
	fmt.Println("Type 'h' or 'help' for available commands.")
	fmt.Println()

	// Show initial room
	gameState.DisplayRoom()

	// Game loop
	running := true
	for running {
		fmt.Print("\n> ")
		if !gameState.Scanner.Scan() {
			break
		}

		input := strings.ToLower(strings.TrimSpace(gameState.Scanner.Text()))

		switch input {
		case "q", "quit", "exit":
			// Save game state before quitting
			fmt.Println("\n💾 Saving game...")
			dungeonSavePath := "./data/" + GenerateTimestampedFilename("dungeon_generated.json")
			metadataSavePath := "./data/" + GenerateTimestampedFilename("dungeon_metadata.json")

			err := gameState.SaveGameState(dungeonSavePath, metadataSavePath)
			if err != nil {
				fmt.Printf("⚠️  Warning: Failed to save game: %v\n", err)
			} else {
				fmt.Printf("✅ Game saved successfully!\n")
				fmt.Printf("   Dungeon: %s\n", dungeonSavePath)
				fmt.Printf("   Metadata: %s\n", metadataSavePath)
			}

			fmt.Println("\n👋 Thanks for playing! Goodbye!")
			running = false

		case "h", "help":
			ShowHelp()

		case "l", "look":
			gameState.DisplayRoom()

		case "m", "map":
			fmt.Println("\n🗺️  Dungeon Map:")
			fmt.Println(dng.GetDetailedGridCompact())

		case "i", "inventory":
			gameState.Player.ShowInventory()

		case "c", "collect":
			gameState.CollectItems()

		case "d", "drink":
			success := gameState.Player.DrinkPotion()
			// If in combat and potion was drunk, monster gets a free attack
			if gameState.InCombat && success {
				gameState.MonsterCounterAttack()
			}

		case "t", "talk":
			gameState.TalkToNPCs()

		case "f", "fight":
			if gameState.InCombat {
				gameState.PerformCombatRound()
			} else {
				gameState.InitiateCombat()
			}

		case "n", "north", "s", "south", "e", "east", "w", "west":
			// If in combat, fleeing
			if gameState.InCombat {
				gameState.FleeCombat()
			}

			direction, valid := ParseDirection(input)
			if valid {
				if gameState.MovePlayer(direction) {
					gameState.DisplayRoom()
				}
			} else {
				fmt.Println("❌ Invalid direction!")
			}

		default:
			if input != "" {
				fmt.Printf("❓ Unknown command: '%s'. Type 'h' for help.\n", input)
			}
		}
	}
}

// loadEntitiesFromDungeon loads all collectibles and monsters from the dungeon data
func loadEntitiesFromDungeon(gs *GameState, metadata *DungeonMetadata) {
	rooms := gs.Dungeon.GetRooms()

	for _, room := range rooms {
		for _, char := range room.Chars {
			kind := string(char.Kind)
			objectID := char.ObjectID

			// Get metadata for this entity
			entityMeta := metadata.GetEntityMetadata(objectID)

			switch kind {
			case "gold":
				gold := &GoldCoins{
					Id:     objectID,
					Type:   dungeon.Gold,
					dng:    gs.Dungeon,
					patterns: dungeon.ObjectPatterns{
						Simple: '*',
					},
					Amount: 100, // Default amount
					RoomId: room.ID,
				}
				gs.Collectibles[objectID] = gold

			case "potion":
				potion := &MagicPotion{
					Id:     objectID,
					Type:   dungeon.Potion,
					dng:    gs.Dungeon,
					patterns: dungeon.ObjectPatterns{
						Simple: 'Y',
					},
					Health: 50, // Default health
					RoomId: room.ID,
				}
				gs.Collectibles[objectID] = potion

			case "skeleton":
				health := 10
				strength := 5
				if entityMeta != nil {
					health = entityMeta.Health
					strength = entityMeta.Strength
				}
				skeleton := &Skeleton{
					Id:       objectID,
					Type:     dungeon.Skeleton,
					dng:      gs.Dungeon,
					Health:   health,
					Strength: strength,
					RoomId:   room.ID,
				}
				gs.Monsters[objectID] = skeleton

			case "goblin":
				health := 10
				strength := 5
				if entityMeta != nil {
					health = entityMeta.Health
					strength = entityMeta.Strength
				}
				goblin := &Goblin{
					Id:       objectID,
					Type:     dungeon.Goblin,
					dng:      gs.Dungeon,
					Health:   health,
					Strength: strength,
					RoomId:   room.ID,
				}
				gs.Monsters[objectID] = goblin

			case "vampire":
				health := 10
				strength := 5
				if entityMeta != nil {
					health = entityMeta.Health
					strength = entityMeta.Strength
				}
				vampire := &Vampire{
					Id:       objectID,
					Type:     dungeon.Vampire,
					dng:      gs.Dungeon,
					Health:   health,
					Strength: strength,
					RoomId:   room.ID,
				}
				gs.Monsters[objectID] = vampire

			case "sphinx":
				health := 10
				strength := 5
				if entityMeta != nil {
					health = entityMeta.Health
					strength = entityMeta.Strength
				}
				sphinx := &Sphinx{
					Id:       objectID,
					Type:     dungeon.Sphinx,
					dng:      gs.Dungeon,
					Health:   health,
					Strength: strength,
					RoomId:   room.ID,
				}
				gs.NPCs[objectID] = sphinx

			case "elf":
				health := 15
				strength := 7
				if entityMeta != nil {
					health = entityMeta.Health
					strength = entityMeta.Strength
				}
				elf := &Elf{
					Id:       objectID,
					Type:     dungeon.Elf,
					dng:      gs.Dungeon,
					Health:   health,
					Strength: strength,
					RoomId:   room.ID,
				}
				gs.NPCs[objectID] = elf

			case "dwarf":
				health := 20
				strength := 10
				if entityMeta != nil {
					health = entityMeta.Health
					strength = entityMeta.Strength
				}
				dwarf := &Dwarf{
					Id:       objectID,
					Type:     dungeon.Dwarf,
					dng:      gs.Dungeon,
					Health:   health,
					Strength: strength,
					RoomId:   room.ID,
				}
				gs.NPCs[objectID] = dwarf

			case "human":
				health := 15
				strength := 8
				if entityMeta != nil {
					health = entityMeta.Health
					strength = entityMeta.Strength
				}
				human := &Human{
					Id:       objectID,
					Type:     dungeon.Human,
					dng:      gs.Dungeon,
					Health:   health,
					Strength: strength,
					RoomId:   room.ID,
				}
				gs.NPCs[objectID] = human
			}
		}
	}
}

// findPlayer finds the player character in the dungeon
func findPlayer(gs *GameState, metadata *DungeonMetadata) *Player {
	rooms := gs.Dungeon.GetRooms()

	for _, room := range rooms {
		for _, char := range room.Chars {
			if string(char.Kind) == "player" {
				// Default values
				name := "Unknown Hero"
				race := "Human"
				class := "Warrior"
				health := 50
				strength := 5
				inventory := Inventory{
					GoldCoins: 0,
					GoldIDs:   []string{},
					Potions:   0,
					PotionIDs: []string{},
				}

				// Load from metadata if available
				entityMeta := metadata.GetEntityMetadata(char.ObjectID)
				if entityMeta != nil {
					if entityMeta.Name != "" {
						name = entityMeta.Name
					}
					if entityMeta.Race != "" {
						race = entityMeta.Race
					}
					if entityMeta.Class != "" {
						class = entityMeta.Class
					}
					health = entityMeta.Health
					strength = entityMeta.Strength

					// Restore inventory if available
					if entityMeta.Inventory != nil {
						inventory.GoldCoins = entityMeta.Inventory.GoldCoins
						inventory.GoldIDs = entityMeta.Inventory.GoldIDs
						inventory.Potions = entityMeta.Inventory.Potions
						inventory.PotionIDs = entityMeta.Inventory.PotionIDs
					}
				}

				player := &Player{
					Id:        char.ObjectID,
					Name:      name,
					Race:      race,
					Class:     class,
					Type:      dungeon.Player,
					dng:       gs.Dungeon,
					RoomId:    room.ID,
					Health:    health,
					Strength:  strength,
					Inventory: inventory,
				}
				return player
			}
		}
	}
	return nil
}
