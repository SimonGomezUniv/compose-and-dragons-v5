package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"codeberg.org/rpg/dungeon/dungeon"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NPCDialogue represents a dialogue from an NPC
type NPCDialogue struct {
	NPCType  string `json:"npc_type"`
	NPCID    string `json:"npc_id"`
	Dialogue string `json:"dialogue"`
}

// GameState holds the current state of the game
type GameState struct {
	Dungeon        *dungeon.Dungeon
	Player         *Player
	Collectibles   map[string]interface{} // Map of object_id to collectible object
	Monsters       map[string]interface{} // Map of object_id to monster object
	NPCs           map[string]interface{} // Map of object_id to NPC object
	Scanner        *bufio.Scanner          // Optional: only used in CLI mode
	InCombat       bool                   // Is the player currently in combat
	CurrentEnemy   interface{}            // Pointer to the current enemy
	RiddleSolved   bool                   // Has the player solved the Sphinx's riddle
}

// NewGameState creates a new game state
func NewGameState(dng *dungeon.Dungeon) *GameState {
	return &GameState{
		Dungeon:      dng,
		Collectibles: make(map[string]interface{}),
		Monsters:     make(map[string]interface{}),
		NPCs:         make(map[string]interface{}),
		Scanner:      bufio.NewScanner(os.Stdin),
		InCombat:     false,
		CurrentEnemy: nil,
	}
}

// DisplayRoom shows the current room information
func (gs *GameState) DisplayRoom() {
	room := gs.GetCurrentRoom()
	if room == nil {
		fmt.Println("Error: Room not found!")
		return
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("🏰 Room #%d: %s\n", room.ID, room.Name)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("📖 %s\n", room.Description)
	fmt.Println(strings.Repeat("-", 60))

	// Show available exits
	fmt.Println("\n🚪 Available exits:")
	if len(room.Connections) == 0 {
		fmt.Println("  None - You're trapped!")
	} else {
		for direction, targetRoom := range room.Connections {
			fmt.Printf("  - %s (leads to Room #%d: %s)\n", direction, targetRoom.ID, targetRoom.Name)
		}
	}

	// Show monsters with their stats
	monsters := gs.GetMonstersInRoom(room.ID)
	if len(monsters) > 0 {
		fmt.Println("\n⚔️  Monsters in this room:")

		// Get unique monster IDs to avoid duplicates (from Shape3x3)
		seenMonsters := make(map[string]bool)
		for _, char := range room.Chars {
			kind := string(char.Kind)
			objectID := char.ObjectID

			if (kind == "skeleton" || kind == "goblin" || kind == "vampire") && !seenMonsters[objectID] {
				seenMonsters[objectID] = true

				// Get monster stats from the Monsters map
				if monsterObj, exists := gs.Monsters[objectID]; exists {
					switch m := monsterObj.(type) {
					case *Skeleton:
						fmt.Printf("  - 💀 Skeleton [Health: %d ❤️, Strength: %d 💪]\n", m.Health, m.Strength)
					case *Goblin:
						fmt.Printf("  - 👺 Goblin [Health: %d ❤️, Strength: %d 💪]\n", m.Health, m.Strength)
					case *Vampire:
						fmt.Printf("  - 🧛 Vampire [Health: %d ❤️, Strength: %d 💪]\n", m.Health, m.Strength)
					}
				}
			}
		}

		// Show player stats when there are monsters
		fmt.Println("\n🛡️  Your Stats:")
		fmt.Printf("  - Health: %d ❤️, Strength: %d 💪\n", gs.Player.Health, gs.Player.Strength)

		// Show combat prompt if not already in combat
		if !gs.InCombat {
			fmt.Println("\n⚠️  DANGER! Type 'f' or 'fight' to engage in combat")
			fmt.Println("💬 Or move to another room to avoid the fight")
		}
	}

	// Show NPCs in the room
	npcs := gs.GetNPCsInRoom(room.ID)
	if len(npcs) > 0 {
		fmt.Println("\n👥 NPCs in this room:")

		// Get unique NPC IDs to avoid duplicates (from Shape3x3)
		seenNPCs := make(map[string]bool)
		for _, char := range room.Chars {
			kind := string(char.Kind)
			objectID := char.ObjectID

			if (kind == "elf" || kind == "dwarf" || kind == "human" || kind == "sphinx") && !seenNPCs[objectID] {
				seenNPCs[objectID] = true

				// Get NPC info from the NPCs map
				if npcObj, exists := gs.NPCs[objectID]; exists {
					switch n := npcObj.(type) {
					case *Elf:
						fmt.Printf("  - 🧝 %s (Elf)\n", n.Id)
					case *Dwarf:
						fmt.Printf("  - 🧔 %s (Dwarf)\n", n.Id)
					case *Human:
						fmt.Printf("  - 👤 %s (Human)\n", n.Id)
					case *Sphinx:
						fmt.Printf("  - 🦁 %s (Sphinx) - Guards the exit with a riddle\n", n.Id)
					}
				}
			}
		}
		fmt.Println("💬 Type 't' or 'talk' to speak with them")
	}

	// Show collectibles
	collectibles := gs.GetCollectiblesInRoom(room.ID)
	if len(collectibles) > 0 {
		fmt.Println("\n💎 Items available:")
		for kind := range collectibles {
			fmt.Printf("  - %s\n", kind)
		}
	}

	fmt.Println(strings.Repeat("=", 60))
}

// GetCurrentRoom returns the room where the player is currently located
func (gs *GameState) GetCurrentRoom() *dungeon.Room {
	rooms := gs.Dungeon.GetRooms()
	for _, room := range rooms {
		if room.ID == gs.Player.RoomId {
			return room
		}
	}
	return nil
}

// GetMonstersInRoom returns a map of monster types in the given room
func (gs *GameState) GetMonstersInRoom(roomID int) map[string]int {
	monsters := make(map[string]int)
	rooms := gs.Dungeon.GetRooms()

	for _, room := range rooms {
		if room.ID == roomID {
			for _, char := range room.Chars {
				kind := string(char.Kind)
				if kind == "skeleton" || kind == "goblin" || kind == "vampire" {
					monsters[kind]++
				}
			}
		}
	}
	return monsters
}

// GetNPCsInRoom returns a map of NPC types in the given room
func (gs *GameState) GetNPCsInRoom(roomID int) map[string]int {
	npcs := make(map[string]int)
	rooms := gs.Dungeon.GetRooms()

	for _, room := range rooms {
		if room.ID == roomID {
			for _, char := range room.Chars {
				kind := string(char.Kind)
				if kind == "elf" || kind == "dwarf" || kind == "human" || kind == "sphinx" {
					npcs[kind]++
				}
			}
		}
	}
	return npcs
}

// GetCollectiblesInRoom returns a map of collectible types in the given room
// Returns ALL object IDs found in the room (including duplicates)
func (gs *GameState) GetCollectiblesInRoom(roomID int) map[string][]string {
	collectibles := make(map[string][]string)
	rooms := gs.Dungeon.GetRooms()

	for _, room := range rooms {
		if room.ID == roomID {
			for _, char := range room.Chars {
				kind := string(char.Kind)
				if kind == "gold" || kind == "potion" {
					// Add ALL occurrences, including duplicates
					collectibles[kind] = append(collectibles[kind], char.ObjectID)
				}
			}
		}
	}
	return collectibles
}

// ParseDirection converts user input to a direction
func ParseDirection(input string) (dungeon.Direction, bool) {
	input = strings.ToLower(strings.TrimSpace(input))

	switch input {
	case "n", "north":
		return dungeon.North, true
	case "s", "south":
		return dungeon.South, true
	case "e", "east":
		return dungeon.East, true
	case "w", "west":
		return dungeon.West, true
	default:
		return dungeon.North, false
	}
}

// MovePlayer moves the player in the specified direction
func (gs *GameState) MovePlayer(direction dungeon.Direction) bool {
	room := gs.GetCurrentRoom()
	if room == nil {
		fmt.Println("\n❌ Cannot find current room!")
		return false
	}

	// Check if there's a connection in the desired direction
	targetRoom, hasConnection := room.Connections[direction]
	if !hasConnection {
		fmt.Printf("\n❌ Cannot move %s: No connection from Room #%d\n", direction, room.ID)
		return false
	}

	// Determine player position based on room contents
	// If there's an NPC/monster, move to (9, 2), otherwise center at (6, 2)
	playerX := 6
	if gs.hasNPCOrMonster(targetRoom.ID) {
		playerX = 9
	}

	// Move the player to the target room
	success := gs.Player.MoveToRoom(dungeon.MoveToRoomArgs{
		Id:         gs.Player.Id,
		RoomNumber: targetRoom.ID,
		X:          playerX,
		Y:          2,
	})

	if !success {
		fmt.Printf("\n❌ Failed to move to %s\n", targetRoom.Name)
		return false
	}

	// Force the player's color to ColorBrightBlue to ensure it's always preserved
	gs.forcePlayerColor()

	fmt.Printf("\n✅ You moved %s to %s (Room #%d)\n", direction, targetRoom.Name, targetRoom.ID)
	return true
}

// CollectItems allows the player to collect items in the current room
func (gs *GameState) CollectItems() {
	room := gs.GetCurrentRoom()
	if room == nil {
		return
	}

	collectibles := gs.GetCollectiblesInRoom(room.ID)
	if len(collectibles) == 0 {
		fmt.Println("\n📭 No items to collect in this room.")
		return
	}

	// Collect gold - call Remove for EACH instance
	if goldIDs, hasGold := collectibles["gold"]; hasGold {
		for _, goldID := range goldIDs {
			// Find the gold object in the collectibles map
			if goldObj, exists := gs.Collectibles[goldID]; exists {
				if gold, ok := goldObj.(*GoldCoins); ok {
					gs.Player.CollectGold(gold.Amount, goldID)
					// Remove this specific instance from the dungeon
					success := gs.Dungeon.RemoveCharObjectFromRoom(room.ID, goldID)
					fmt.Printf("💰 Collected %d gold coins (ID: %s) - Removed: %v\n", gold.Amount, goldID, success)
				}
			}
		}
	}

	// Collect potions - call Remove for EACH instance
	if potionIDs, hasPotions := collectibles["potion"]; hasPotions {
		for _, potionID := range potionIDs {
			if potionObj, exists := gs.Collectibles[potionID]; exists {
				if potion, ok := potionObj.(*MagicPotion); ok {
					gs.Player.CollectPotion(potionID)
					// Remove this specific instance from the dungeon
					success := gs.Dungeon.RemoveCharObjectFromRoom(room.ID, potionID)
					fmt.Printf("🧪 Collected a magic potion (ID: %s, Health: %d) - Removed: %v\n", potionID, potion.Health, success)
				}
			}
		}
	}
}

// ShowHelp displays available commands
func ShowHelp() {
	fmt.Println("\n📜 Available Commands:")
	fmt.Println("  n, north - Move north")
	fmt.Println("  s, south - Move south")
	fmt.Println("  e, east  - Move east")
	fmt.Println("  w, west  - Move west")
	fmt.Println("  c, collect - Collect items in current room")
	fmt.Println("  d, drink - Drink a potion to restore 5 health")
	fmt.Println("  t, talk - Talk to NPCs in current room")
	fmt.Println("  f, fight - Start/continue combat with monsters")
	fmt.Println("  i, inventory - Show your inventory")
	fmt.Println("  m, map - Show the dungeon map")
	fmt.Println("  l, look - Look around the current room")
	fmt.Println("  h, help - Show this help")
	fmt.Println("  q, quit - Quit the game")
	fmt.Println("\n💡 During combat:")
	fmt.Println("  f - Attack the monster")
	fmt.Println("  d - Drink a potion (monster attacks)")
	fmt.Println("  n/s/e/w - Flee from combat")
}

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// TalkToNPCs allows the player to talk to NPCs in the current room
func (gs *GameState) TalkToNPCs() {
	room := gs.GetCurrentRoom()
	if room == nil {
		return
	}

	// Get unique NPCs in the room
	npcList := []interface{}{}
	seenNPCs := make(map[string]bool)

	for _, char := range room.Chars {
		kind := string(char.Kind)
		objectID := char.ObjectID

		if (kind == "elf" || kind == "dwarf" || kind == "human" || kind == "sphinx") && !seenNPCs[objectID] {
			seenNPCs[objectID] = true
			if npcObj, exists := gs.NPCs[objectID]; exists {
				npcList = append(npcList, npcObj)
			}
		}
	}

	if len(npcList) == 0 {
		fmt.Println("\n❌ There are no NPCs in this room to talk to.")
		return
	}

	// Talk to each NPC
	for _, npcObj := range npcList {
		switch n := npcObj.(type) {
		case *Elf:
			fmt.Printf("\n🧝 %s (Elf) says:\n", n.Id)
			fmt.Println("   \"Greetings, traveler! The ancient trees whisper of your quest.\"")
			fmt.Println("   \"May the light of the stars guide your path through these dark halls.\"")
		case *Dwarf:
			fmt.Printf("\n🧔 %s (Dwarf) says:\n", n.Id)
			fmt.Println("   \"Aye, welcome to these forsaken halls! I've been mining here for years.\"")
			fmt.Println("   \"Watch out for the monsters - they're tougher than they look!\"")
			fmt.Println("   \"If you need strength, remember: a good meal and rest do wonders.\"")
		case *Human:
			fmt.Printf("\n👤 %s (Human) says:\n", n.Id)
			fmt.Println("   \"Hello there! I'm a fellow adventurer trapped in this dungeon.\"")
			fmt.Println("   \"I've heard rumors of great treasures deeper within...\"")
			fmt.Println("   \"Be careful, and may fortune favor you!\"")
		case *Sphinx:
			fmt.Printf("\n🦁 %s (Sphinx) says:\n", n.Id)
			if gs.RiddleSolved {
				fmt.Println("   \"You have proven your wisdom, brave one.\"")
				fmt.Println("   \"The path to freedom is now open to you.\"")
			} else {
				fmt.Println("   \"Halt, mortal! I am the guardian of the exit.\"")
				fmt.Println("   \"Answer my riddle correctly, and you may leave this dungeon.\"")
				fmt.Println("   \"Fail, and you shall wander these halls forever...\"")
				fmt.Println("\n   🔮 THE RIDDLE:")
				fmt.Println("   \"I speak without a mouth and hear without ears.")
				fmt.Println("   I have no body, but I come alive with wind.")
				fmt.Println("   What am I?\"")
				fmt.Print("\n   Your answer: ")

				if gs.Scanner.Scan() {
					answer := strings.TrimSpace(strings.ToLower(gs.Scanner.Text()))

					if answer == "echo" || answer == "an echo" {
						gs.RiddleSolved = true
						fmt.Println("\n   ✨ \"Correct! You are wise indeed.\"")
						fmt.Println("   \"The exit is now accessible to you. Go forth, brave adventurer!\"")
						fmt.Println("\n🎉 Congratulations! You have solved the Sphinx's riddle!")
						fmt.Println("You can now exit the dungeon to complete your quest!")
					} else {
						fmt.Println("\n   ❌ \"Wrong! That is not the answer I seek.\"")
						fmt.Println("   \"You may try again when you return...\"")
					}
				}
			}
		}
	}
}

// TalkToNPCsNonInteractive returns dialogues from NPCs without interactive input
// This is the MCP-friendly version that returns data instead of printing
func (gs *GameState) TalkToNPCsNonInteractive() []NPCDialogue {
	room := gs.GetCurrentRoom()
	if room == nil {
		return []NPCDialogue{}
	}

	// Get unique NPCs in the room
	var dialogues []NPCDialogue
	seenNPCs := make(map[string]bool)

	for _, char := range room.Chars {
		kind := string(char.Kind)
		objectID := char.ObjectID

		if (kind == "elf" || kind == "dwarf" || kind == "human" || kind == "sphinx") && !seenNPCs[objectID] {
			seenNPCs[objectID] = true
			if npcObj, exists := gs.NPCs[objectID]; exists {
				var dialogue string
				var npcID string

				switch n := npcObj.(type) {
				case *Elf:
					npcID = n.Id
					dialogue = "Greetings, traveler! The ancient trees whisper of your quest. May the light of the stars guide your path through these dark halls."
				case *Dwarf:
					npcID = n.Id
					dialogue = "Aye, welcome to these forsaken halls! I've been mining here for years. Watch out for the monsters - they're tougher than they look! If you need strength, remember: a good meal and rest do wonders."
				case *Human:
					npcID = n.Id
					dialogue = "Hello there! I'm a fellow adventurer trapped in this dungeon. I've heard rumors of great treasures deeper within... Be careful, and may fortune favor you!"
				case *Sphinx:
					npcID = n.Id
					if gs.RiddleSolved {
						dialogue = "You have proven your wisdom, brave one. The path to freedom is now open to you."
					} else {
						dialogue = "Halt, mortal! I am the guardian of the exit. Answer my riddle correctly, and you may leave this dungeon. Fail, and you shall wander these halls forever... THE RIDDLE: I speak without a mouth and hear without ears. I have no body, but I come alive with wind. What am I?"
					}
				}

				dialogues = append(dialogues, NPCDialogue{
					NPCType:  kind,
					NPCID:    npcID,
					Dialogue: dialogue,
				})
			}
		}
	}

	return dialogues
}

// AnswerRiddle attempts to answer the Sphinx's riddle
// Returns true if correct, false otherwise, along with a message
func (gs *GameState) AnswerRiddle(answer string) (bool, string) {
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "echo" || answer == "an echo" {
		gs.RiddleSolved = true
		return true, "Correct! You are wise indeed. The exit is now accessible to you. Go forth, brave adventurer!"
	}

	return false, "Wrong! That is not the answer I seek. You may try again when you return..."
}

// SaveGameState saves the current game state to timestamped files
func (gs *GameState) SaveGameState(dungeonPath, metadataPath string) error {
	// Save the dungeon state
	err := gs.Dungeon.SaveToFile(dungeonPath)
	if err != nil {
		return fmt.Errorf("error saving dungeon: %v", err)
	}

	// Create metadata with current state
	metadata := &DungeonMetadata{
		Entities:     []EntityMetadata{},
		RiddleSolved: gs.RiddleSolved,
	}

	// Add player metadata with current state including inventory
	metadata.AddEntityMetadata(EntityMetadata{
		ID:       gs.Player.Id,
		Type:     "player",
		Name:     gs.Player.Name,
		Race:     gs.Player.Race,
		Class:    gs.Player.Class,
		Health:   gs.Player.Health,
		Strength: gs.Player.Strength,
		Inventory: &PlayerInventoryMetadata{
			GoldCoins: gs.Player.Inventory.GoldCoins,
			GoldIDs:   gs.Player.Inventory.GoldIDs,
			Potions:   gs.Player.Inventory.Potions,
			PotionIDs: gs.Player.Inventory.PotionIDs,
		},
	})

	// Add all monsters with current state
	for id, monsterObj := range gs.Monsters {
		var entityType string
		var health, strength int

		switch m := monsterObj.(type) {
		case *Skeleton:
			entityType = "skeleton"
			health = m.Health
			strength = m.Strength
		case *Goblin:
			entityType = "goblin"
			health = m.Health
			strength = m.Strength
		case *Vampire:
			entityType = "vampire"
			health = m.Health
			strength = m.Strength
		case *Sphinx:
			entityType = "sphinx"
			health = m.Health
			strength = m.Strength
		}

		if entityType != "" {
			metadata.AddEntityMetadata(EntityMetadata{
				ID:       id,
				Type:     entityType,
				Health:   health,
				Strength: strength,
			})
		}
	}

	// Add all NPCs with current state
	for id, npcObj := range gs.NPCs {
		var entityType string
		var health, strength int

		switch n := npcObj.(type) {
		case *Elf:
			entityType = "elf"
			health = n.Health
			strength = n.Strength
		case *Dwarf:
			entityType = "dwarf"
			health = n.Health
			strength = n.Strength
		case *Human:
			entityType = "human"
			health = n.Health
			strength = n.Strength
		}

		if entityType != "" {
			metadata.AddEntityMetadata(EntityMetadata{
				ID:       id,
				Type:     entityType,
				Health:   health,
				Strength: strength,
			})
		}
	}

	// Save metadata
	err = SaveMetadata(metadata, metadataPath)
	if err != nil {
		return fmt.Errorf("error saving metadata: %v", err)
	}

	return nil
}

// rollDice rolls n dice with d faces and returns the sum
func rollDice(n, d int) int {
	total := 0
	for i := 0; i < n; i++ {
		total += rand.Intn(d) + 1
	}
	return total
}

// InitiateCombat starts combat with the first monster in the current room
func (gs *GameState) InitiateCombat() {
	room := gs.GetCurrentRoom()
	if room == nil {
		return
	}

	// Find the first monster in the room
	seenMonsters := make(map[string]bool)
	for _, char := range room.Chars {
		kind := string(char.Kind)
		objectID := char.ObjectID

		if (kind == "skeleton" || kind == "goblin" || kind == "vampire") && !seenMonsters[objectID] {
			seenMonsters[objectID] = true

			if monsterObj, exists := gs.Monsters[objectID]; exists {
				gs.InCombat = true
				gs.CurrentEnemy = monsterObj

				// Display combat start message
				var monsterName string
				var monsterHealth, monsterStrength int

				switch m := monsterObj.(type) {
				case *Skeleton:
					monsterName = "💀 Skeleton"
					monsterHealth = m.Health
					monsterStrength = m.Strength
				case *Goblin:
					monsterName = "👺 Goblin"
					monsterHealth = m.Health
					monsterStrength = m.Strength
				case *Vampire:
					monsterName = "🧛 Vampire"
					monsterHealth = m.Health
					monsterStrength = m.Strength
				}

				fmt.Println("\n" + strings.Repeat("=", 60))
				fmt.Println("⚔️  COMBAT STARTED!")
				fmt.Println(strings.Repeat("=", 60))
				fmt.Printf("You: [%d ❤️, %d 💪] vs %s [%d ❤️, %d 💪]\n",
					gs.Player.Health, gs.Player.Strength, monsterName, monsterHealth, monsterStrength)
				fmt.Println(strings.Repeat("-", 60))
				fmt.Println("💬 Commands: 'f' to attack, 'd' to drink potion, or move to flee")
				fmt.Println(strings.Repeat("=", 60))
				return
			}
		}
	}

	fmt.Println("\n❌ No monsters to fight in this room!")
}

// PerformCombatRound executes one round of combat
func (gs *GameState) PerformCombatRound() {
	if !gs.InCombat || gs.CurrentEnemy == nil {
		fmt.Println("\n❌ You are not in combat!")
		return
	}

	var monsterName string
	var monsterHealth, monsterStrength *int
	var monsterId string

	switch m := gs.CurrentEnemy.(type) {
	case *Skeleton:
		monsterName = "💀 Skeleton"
		monsterHealth = &m.Health
		monsterStrength = &m.Strength
		monsterId = m.Id
	case *Goblin:
		monsterName = "👺 Goblin"
		monsterHealth = &m.Health
		monsterStrength = &m.Strength
		monsterId = m.Id
	case *Vampire:
		monsterName = "🧛 Vampire"
		monsterHealth = &m.Health
		monsterStrength = &m.Strength
		monsterId = m.Id
	default:
		fmt.Println("❌ Invalid enemy!")
		gs.InCombat = false
		gs.CurrentEnemy = nil
		return
	}

	fmt.Println("\n" + strings.Repeat("-", 60))

	// Player's turn
	playerRoll := rollDice(3, 6)
	playerTotal := playerRoll + gs.Player.Strength
	fmt.Printf("🎲 Your roll: 3d6(%d) + %d💪 = %d\n", playerRoll, gs.Player.Strength, playerTotal)

	// Monster's turn
	monsterRoll := rollDice(3, 6)
	monsterTotal := monsterRoll + *monsterStrength
	fmt.Printf("🎲 %s roll: 3d6(%d) + %d💪 = %d\n", monsterName, monsterRoll, *monsterStrength, monsterTotal)

	fmt.Println(strings.Repeat("-", 60))

	// Determine outcome
	if playerTotal > monsterTotal {
		damage := (playerTotal - monsterTotal + 1) / 2
		*monsterHealth -= damage
		fmt.Printf("💥 You hit! %s takes %d damage. (%d ❤️ remaining)\n", monsterName, damage, *monsterHealth)

		// Check if monster is dead
		if *monsterHealth <= 0 {
			fmt.Println(strings.Repeat("=", 60))
			fmt.Printf("🎉 VICTORY! You defeated the %s!\n", monsterName)
			fmt.Println(strings.Repeat("=", 60))

			// Remove monster from the game
			delete(gs.Monsters, monsterId)
			gs.Dungeon.RemoveCharObjectFromRoom(gs.Player.RoomId, monsterId)

			gs.InCombat = false
			gs.CurrentEnemy = nil
		}
	} else if monsterTotal > playerTotal {
		damage := (monsterTotal - playerTotal + 1) / 2
		gs.Player.Health -= damage
		fmt.Printf("💔 %s hits you! You take %d damage. (%d ❤️ remaining)\n", monsterName, damage, gs.Player.Health)

		// Check if player is dead
		if gs.Player.Health <= 0 {
			fmt.Println(strings.Repeat("=", 60))
			fmt.Println("💀 GAME OVER! You have been defeated...")
			fmt.Println(strings.Repeat("=", 60))
			fmt.Println("The dungeon claims another soul.")
			os.Exit(0)
		}
	} else {
		fmt.Println("⚔️  The attacks clash! No one takes damage this round.")
	}

	if gs.InCombat {
		fmt.Println(strings.Repeat("-", 60))
		fmt.Printf("Status: You [%d ❤️] vs %s [%d ❤️]\n", gs.Player.Health, monsterName, *monsterHealth)
		fmt.Println("💬 'f' to attack, 'd' to drink, or move to flee")
	}
}

// FleeCombat allows the player to flee from combat
func (gs *GameState) FleeCombat() {
	if gs.InCombat {
		fmt.Println("\n🏃 You flee from combat!")
		gs.InCombat = false
		gs.CurrentEnemy = nil
	}
}

// MonsterCounterAttack allows the monster to attack when player drinks potion
func (gs *GameState) MonsterCounterAttack() {
	if !gs.InCombat || gs.CurrentEnemy == nil {
		return
	}

	var monsterName string
	var monsterStrength *int

	switch m := gs.CurrentEnemy.(type) {
	case *Skeleton:
		monsterName = "💀 Skeleton"
		monsterStrength = &m.Strength
	case *Goblin:
		monsterName = "👺 Goblin"
		monsterStrength = &m.Strength
	case *Vampire:
		monsterName = "🧛 Vampire"
		monsterStrength = &m.Strength
	default:
		return
	}

	fmt.Println("\n⚠️  The monster attacks while you drink!")
	monsterRoll := rollDice(3, 6)
	monsterTotal := monsterRoll + *monsterStrength
	fmt.Printf("🎲 %s roll: 3d6(%d) + %d💪 = %d\n", monsterName, monsterRoll, *monsterStrength, monsterTotal)

	// Monster damage calculation (half of what it would be in combat)
	damage := (monsterTotal + 1) / 4
	if damage < 1 {
		damage = 1
	}

	gs.Player.Health -= damage
	fmt.Printf("💔 %s hits you! You take %d damage. (%d ❤️ remaining)\n", monsterName, damage, gs.Player.Health)

	// Check if player is dead
	if gs.Player.Health <= 0 {
		fmt.Println(strings.Repeat("=", 60))
		fmt.Println("💀 GAME OVER! You have been defeated...")
		fmt.Println(strings.Repeat("=", 60))
		fmt.Println("The dungeon claims another soul.")
		os.Exit(0)
	}
}

// hasNPCOrMonster checks if a room contains any NPCs or monsters
func (gs *GameState) hasNPCOrMonster(roomID int) bool {
	for _, room := range gs.Dungeon.Rooms {
		if room.ID == roomID {
			for _, char := range room.Chars {
				kind := string(char.Kind)
				if kind == "elf" || kind == "dwarf" || kind == "human" || kind == "sphinx" ||
					kind == "skeleton" || kind == "goblin" || kind == "vampire" {
					return true
				}
			}
			return false
		}
	}
	return false
}

// forcePlayerColor ensures the player @ always keeps their ColorBrightBlue color
func (gs *GameState) forcePlayerColor() {
	// Access the Dungeon.Rooms directly (it's a public field with []*Room)
	for _, room := range gs.Dungeon.Rooms {
		if room.ID == gs.Player.RoomId {
			// Find and update the player character in this room
			for i := range room.Chars {
				if room.Chars[i].ObjectID == gs.Player.Id && room.Chars[i].Kind == dungeon.Player {
					// Force the player's color to ColorBrightBlue
					room.Chars[i].ForeColor = dungeon.ColorBrightBlue
					// Clear any background color
					room.Chars[i].BackColor = ""
					return
				}
			}
		}
	}
}
