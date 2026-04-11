package main

import "codeberg.org/rpg/dungeon/dungeon"

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
