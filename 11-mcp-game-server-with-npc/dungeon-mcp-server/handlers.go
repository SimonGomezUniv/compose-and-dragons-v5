package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// Helper functions

func extractStringParam(args map[string]interface{}, key string) (string, error) {
	val, exists := args[key]
	if !exists {
		return "", fmt.Errorf("missing required parameter: %s", key)
	}
	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("parameter %s must be a string", key)
	}
	return str, nil
}

func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return `{"error": "failed to marshal JSON"}`
	}
	return string(data)
}

// Status handlers

func handleGetGameStatus(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"in_combat":     gameState.InCombat,
		"riddle_solved": gameState.RiddleSolved,
		"player": map[string]interface{}{
			"name":     gameState.Player.Name,
			"race":     gameState.Player.Race,
			"class":    gameState.Player.Class,
			"health":   gameState.Player.Health,
			"strength": gameState.Player.Strength,
		},
	}
	return mcp.NewToolResultText(toJSON(result)), nil
}

func handleSaveGame(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	prefix := "dungeon_generated"
	if val, exists := args["filename_prefix"]; exists && val != nil {
		if str, ok := val.(string); ok && str != "" {
			prefix = str
		}
	}

	dungeonPath := fmt.Sprintf("./data/%s_%s.json", prefix, GenerateTimestamp())
	metadataPath := fmt.Sprintf("./data/%s_metadata_%s.json", prefix, GenerateTimestamp())

	err := gameState.SaveGameState(dungeonPath, metadataPath)

	result := map[string]interface{}{
		"success":       err == nil,
		"dungeon_path":  dungeonPath,
		"metadata_path": metadataPath,
	}
	if err != nil {
		result["error"] = err.Error()
	}

	return mcp.NewToolResultText(toJSON(result)), nil
}

// Query handlers

func handleGetCurrentRoom(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	room := gameState.GetCurrentRoom()
	if room == nil {
		return mcp.NewToolResultText(`{"error": "room not found"}`), nil
	}

	// Build exits list
	exits := []map[string]interface{}{}
	for direction, targetRoom := range room.Connections {
		exits = append(exits, map[string]interface{}{
			"direction":      string(direction),
			"target_room_id": targetRoom.ID,
			"target_name":    targetRoom.Name,
		})
	}

	// Get monsters
	monsters := []map[string]interface{}{}
	seenMonsters := make(map[string]bool)
	for _, char := range room.Chars {
		kind := string(char.Kind)
		objectID := char.ObjectID
		if (kind == "skeleton" || kind == "goblin" || kind == "vampire") && !seenMonsters[objectID] {
			seenMonsters[objectID] = true
			if monsterObj, exists := gameState.Monsters[objectID]; exists {
				var health, strength int
				switch m := monsterObj.(type) {
				case *Skeleton:
					health, strength = m.Health, m.Strength
				case *Goblin:
					health, strength = m.Health, m.Strength
				case *Vampire:
					health, strength = m.Health, m.Strength
				}
				monsters = append(monsters, map[string]interface{}{
					"type":     kind,
					"health":   health,
					"strength": strength,
				})
			}
		}
	}

	// Get NPCs
	npcs := []map[string]interface{}{}
	seenNPCs := make(map[string]bool)
	for _, char := range room.Chars {
		kind := string(char.Kind)
		objectID := char.ObjectID
		if (kind == "elf" || kind == "dwarf" || kind == "human" || kind == "sphinx") && !seenNPCs[objectID] {
			seenNPCs[objectID] = true
			npcs = append(npcs, map[string]interface{}{
				"type": kind,
				"id":   objectID,
			})
		}
	}

	// Get items
	items := []map[string]interface{}{}
	collectibles := gameState.GetCollectiblesInRoom(room.ID)
	for itemType, ids := range collectibles {
		items = append(items, map[string]interface{}{
			"type":  itemType,
			"count": len(ids),
		})
	}

	result := map[string]interface{}{
		"room": map[string]interface{}{
			"id":          room.ID,
			"name":        room.Name,
			"description": room.Description,
			"exits":       exits,
			"monsters":    monsters,
			"npcs":        npcs,
			"items":       items,
		},
	}

	return mcp.NewToolResultText(toJSON(result)), nil
}

func handleGetInventory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"player": map[string]interface{}{
			"name":     gameState.Player.Name,
			"race":     gameState.Player.Race,
			"class":    gameState.Player.Class,
			"health":   gameState.Player.Health,
			"strength": gameState.Player.Strength,
			"inventory": map[string]interface{}{
				"gold_coins": gameState.Player.Inventory.GoldCoins,
				"potions":    gameState.Player.Inventory.Potions,
			},
		},
	}
	return mcp.NewToolResultText(toJSON(result)), nil
}

func handleGetMap(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mapStr := gameState.Dungeon.GetDetailedGridCompact()
	result := map[string]interface{}{
		"map": mapStr,
	}
	return mcp.NewToolResultText(toJSON(result)), nil
}

func handleGetHelp(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Determine current game context
	inCombat := gameState.InCombat

	// Build comprehensive help
	help := map[string]interface{}{
		"game_name": "Dungeon RPG",
		"version":   "1.0.0",
		"current_mode": func() string {
			if inCombat {
				return "combat"
			}
			return "exploration"
		}(),
		"commands": map[string]interface{}{
			"status": []map[string]string{
				{
					"tool":        "get_game_status",
					"description": "Get current game status (combat state, riddle status, player info)",
					"parameters":  "none",
					"example":     "get_game_status",
				},
				{
					"tool":        "save_game",
					"description": "Save current game state to file",
					"parameters":  "filename_prefix (optional)",
					"example":     "save_game",
				},
			},
			"query": []map[string]string{
				{
					"tool":        "get_current_room",
					"description": "Get detailed information about current room (exits, monsters, NPCs, items)",
					"parameters":  "none",
					"example":     "get_current_room",
				},
				{
					"tool":        "get_inventory",
					"description": "View player stats and inventory (health, strength, gold, potions)",
					"parameters":  "none",
					"example":     "get_inventory",
				},
				{
					"tool":        "get_map",
					"description": "Display ASCII art map of entire dungeon with player position",
					"parameters":  "none",
					"example":     "get_map",
				},
				{
					"tool":        "get_help",
					"description": "Show this help information with all available commands",
					"parameters":  "none",
					"example":     "get_help",
				},
			},
			"actions": []map[string]string{
				{
					"tool":        "move",
					"description": "Move player in a direction (flees from combat if in combat)",
					"parameters":  "direction: north|south|east|west",
					"example":     "move direction=north",
				},
				{
					"tool":        "collect_items",
					"description": "Collect all items in current room (gold and potions)",
					"parameters":  "none",
					"example":     "collect_items",
				},
				{
					"tool":        "drink_potion",
					"description": "Drink a health potion to restore 20 HP (monster attacks if in combat)",
					"parameters":  "none",
					"example":     "drink_potion",
				},
				{
					"tool":        "talk_to_npcs",
					"description": "Talk to all NPCs in current room (elf, dwarf, human, sphinx)",
					"parameters":  "none",
					"example":     "talk_to_npcs",
				},
				{
					"tool":        "answer_riddle",
					"description": "Answer the Sphinx's riddle (hint: what speaks without a mouth?)",
					"parameters":  "answer: string",
					"example":     "answer_riddle answer=echo",
				},
			},
			"combat": []map[string]string{
				{
					"tool":        "start_combat",
					"description": "Initiate combat with a monster in current room",
					"parameters":  "none",
					"example":     "start_combat",
				},
				{
					"tool":        "attack",
					"description": "Attack current enemy (must be in combat)",
					"parameters":  "none",
					"example":     "attack",
				},
			},
		},
		"current_context": map[string]interface{}{
			"in_combat":      inCombat,
			"riddle_solved":  gameState.RiddleSolved,
			"player_health":  gameState.Player.Health,
			"player_potions": gameState.Player.Inventory.Potions,
			"available_actions": func() []string {
				if inCombat {
					return []string{"attack", "drink_potion", "move (to flee)"}
				}
				return []string{"move", "collect_items", "talk_to_npcs", "start_combat", "get_current_room", "get_map"}
			}(),
		},
		"tips": []string{
			"Use get_current_room to see what's in your current location",
			"Talk to NPCs for hints and information",
			"Collect potions before entering combat",
			"The Sphinx's riddle answer is: echo",
			"Moving while in combat will flee the battle",
			"Check get_game_status to see if you're in combat",
		},
	}

	return mcp.NewToolResultText(toJSON(help)), nil
}

// Action handlers

func handleMove(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	directionStr, err := extractStringParam(args, "direction")
	if err != nil {
		return mcp.NewToolResultText(toJSON(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})), nil
	}

	direction, valid := ParseDirection(directionStr)
	if !valid {
		return mcp.NewToolResultText(toJSON(map[string]interface{}{
			"success": false,
			"message": "Invalid direction. Use: north, south, east, or west",
		})), nil
	}

	fledCombat := gameState.InCombat
	if gameState.InCombat {
		gameState.FleeCombat()
	}

	success := gameState.MovePlayer(direction)

	result := map[string]interface{}{
		"success":     success,
		"fled_combat": fledCombat && success,
	}

	if success {
		newRoom := gameState.GetCurrentRoom()
		if newRoom != nil {
			result["new_room"] = map[string]interface{}{
				"id":          newRoom.ID,
				"name":        newRoom.Name,
				"description": newRoom.Description,
			}
			result["message"] = fmt.Sprintf("Moved %s to %s", direction, newRoom.Name)
		}
	} else {
		result["message"] = "Cannot move in that direction"
	}

	return mcp.NewToolResultText(toJSON(result)), nil
}

func handleCollectItems(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	room := gameState.GetCurrentRoom()
	if room == nil {
		return mcp.NewToolResultText(toJSON(map[string]interface{}{
			"items_collected": []interface{}{},
		})), nil
	}

	collectibles := gameState.GetCollectiblesInRoom(room.ID)
	itemsCollected := []map[string]interface{}{}

	// Collect gold
	if goldIDs, hasGold := collectibles["gold"]; hasGold {
		for _, goldID := range goldIDs {
			if goldObj, exists := gameState.Collectibles[goldID]; exists {
				if gold, ok := goldObj.(*GoldCoins); ok {
					gameState.Player.CollectGold(gold.Amount, goldID)
					gameState.Dungeon.RemoveCharObjectFromRoom(room.ID, goldID)
					itemsCollected = append(itemsCollected, map[string]interface{}{
						"type":   "gold",
						"id":     goldID,
						"amount": gold.Amount,
					})
				}
			}
		}
	}

	// Collect potions
	if potionIDs, hasPotions := collectibles["potion"]; hasPotions {
		for _, potionID := range potionIDs {
			if potionObj, exists := gameState.Collectibles[potionID]; exists {
				if _, ok := potionObj.(*MagicPotion); ok {
					gameState.Player.CollectPotion(potionID)
					gameState.Dungeon.RemoveCharObjectFromRoom(room.ID, potionID)
					itemsCollected = append(itemsCollected, map[string]interface{}{
						"type": "potion",
						"id":   potionID,
					})
				}
			}
		}
	}

	result := map[string]interface{}{
		"items_collected": itemsCollected,
		"inventory": map[string]interface{}{
			"gold_coins": gameState.Player.Inventory.GoldCoins,
			"potions":    gameState.Player.Inventory.Potions,
		},
	}

	return mcp.NewToolResultText(toJSON(result)), nil
}

func handleDrinkPotion(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if gameState.Player.Inventory.Potions <= 0 {
		return mcp.NewToolResultText(toJSON(map[string]interface{}{
			"success":      false,
			"message":      "No potions available",
			"potions_left": 0,
		})), nil
	}

	oldHealth := gameState.Player.Health
	success := gameState.Player.DrinkPotion()
	healthGained := gameState.Player.Health - oldHealth

	result := map[string]interface{}{
		"success":        success,
		"health_gained":  healthGained,
		"current_health": gameState.Player.Health,
		"potions_left":   gameState.Player.Inventory.Potions,
	}

	// If in combat, monster attacks
	if gameState.InCombat {
		result["monster_attacked"] = true
		gameState.MonsterCounterAttack()
		// Note: damage is printed by MonsterCounterAttack, we could refactor to return it
	} else {
		result["monster_attacked"] = false
	}

	return mcp.NewToolResultText(toJSON(result)), nil
}

func handleTalkToNPCs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	// Check if a specific NPC type was requested
	var targetNPC string
	if val, exists := args["npc_type"]; exists && val != nil {
		if str, ok := val.(string); ok {
			targetNPC = strings.ToLower(strings.TrimSpace(str))
		}
	}

	room := gameState.GetCurrentRoom()
	if room == nil {
		return mcp.NewToolResultText(toJSON(map[string]interface{}{
			"success": false,
			"message": "Room not found",
		})), nil
	}

	// Get unique NPCs in the room
	var dialogues []NPCDialogue
	seenNPCs := make(map[string]bool)
	foundNPCTypes := []string{}

	for _, char := range room.Chars {
		kind := string(char.Kind)
		objectID := char.ObjectID

		if (kind == "elf" || kind == "dwarf" || kind == "human" || kind == "sphinx") && !seenNPCs[objectID] {
			seenNPCs[objectID] = true
			foundNPCTypes = append(foundNPCTypes, kind)

			// If a specific NPC was requested, skip others
			if targetNPC != "" && kind != targetNPC {
				continue
			}

			if npcObj, exists := gameState.NPCs[objectID]; exists {
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
					if gameState.RiddleSolved {
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

	// If a specific NPC was requested but not found, return error message
	if targetNPC != "" && len(dialogues) == 0 {
		return mcp.NewToolResultText(toJSON(map[string]interface{}{
			"success":        false,
			"message":        fmt.Sprintf("NPC type '%s' not found in this room", targetNPC),
			"available_npcs": foundNPCTypes,
			"conversations":  []NPCDialogue{},
		})), nil
	}

	// If no NPCs at all
	if len(foundNPCTypes) == 0 {
		return mcp.NewToolResultText(toJSON(map[string]interface{}{
			"success":        false,
			"message":        "No NPCs in this room",
			"available_npcs": []string{},
			"conversations":  []NPCDialogue{},
		})), nil
	}

	result := map[string]interface{}{
		"success":        true,
		"conversations":  dialogues,
		"available_npcs": foundNPCTypes,
	}

	return mcp.NewToolResultText(toJSON(result)), nil
}

func handleAnswerRiddle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	answer, err := extractStringParam(args, "answer")
	if err != nil {
		return mcp.NewToolResultText(toJSON(map[string]interface{}{
			"correct": false,
			"message": "Answer required",
		})), nil
	}

	correct, message := gameState.AnswerRiddle(answer)

	result := map[string]interface{}{
		"correct":       correct,
		"message":       message,
		"riddle_solved": gameState.RiddleSolved,
	}

	return mcp.NewToolResultText(toJSON(result)), nil
}

// Combat handlers

func handleStartCombat(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	room := gameState.GetCurrentRoom()
	if room == nil {
		return mcp.NewToolResultText(toJSON(map[string]interface{}{
			"success": false,
			"message": "Room not found",
		})), nil
	}

	// Check if already in combat
	if gameState.InCombat {
		return mcp.NewToolResultText(toJSON(map[string]interface{}{
			"success": false,
			"message": "Already in combat",
		})), nil
	}

	// Find first monster
	seenMonsters := make(map[string]bool)
	for _, char := range room.Chars {
		kind := string(char.Kind)
		objectID := char.ObjectID

		if (kind == "skeleton" || kind == "goblin" || kind == "vampire") && !seenMonsters[objectID] {
			seenMonsters[objectID] = true

			if monsterObj, exists := gameState.Monsters[objectID]; exists {
				gameState.InCombat = true
				gameState.CurrentEnemy = monsterObj

				var monsterType string
				var monsterHealth, monsterStrength int

				switch m := monsterObj.(type) {
				case *Skeleton:
					monsterType, monsterHealth, monsterStrength = "skeleton", m.Health, m.Strength
				case *Goblin:
					monsterType, monsterHealth, monsterStrength = "goblin", m.Health, m.Strength
				case *Vampire:
					monsterType, monsterHealth, monsterStrength = "vampire", m.Health, m.Strength
				}

				result := map[string]interface{}{
					"success": true,
					"enemy": map[string]interface{}{
						"type":     monsterType,
						"health":   monsterHealth,
						"strength": monsterStrength,
					},
					"player": map[string]interface{}{
						"health":   gameState.Player.Health,
						"strength": gameState.Player.Strength,
					},
				}

				return mcp.NewToolResultText(toJSON(result)), nil
			}
		}
	}

	return mcp.NewToolResultText(toJSON(map[string]interface{}{
		"success": false,
		"message": "No monsters in this room",
	})), nil
}

func handleAttack(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !gameState.InCombat || gameState.CurrentEnemy == nil {
		return mcp.NewToolResultText(toJSON(map[string]interface{}{
			"combat_ended": true,
			"message":      "Not in combat",
		})), nil
	}

	// Execute combat round (simplified - actual implementation would need to capture dice rolls)
	gameState.PerformCombatRound()

	// Check if combat ended
	combatEnded := !gameState.InCombat
	victor := ""

	if combatEnded {
		if gameState.Player.Health > 0 {
			victor = "player"
		} else {
			victor = "monster"
		}
	}

	result := map[string]interface{}{
		"combat_ended": combatEnded,
		"victor":       victor,
		// Note: detailed combat stats would require refactoring PerformCombatRound to return data
	}

	return mcp.NewToolResultText(toJSON(result)), nil
}

// GenerateTimestamp generates a timestamp string for filenames
func GenerateTimestamp() string {
	return fmt.Sprintf("%s", GenerateTimestampedFilename("temp")[5:20])
}
