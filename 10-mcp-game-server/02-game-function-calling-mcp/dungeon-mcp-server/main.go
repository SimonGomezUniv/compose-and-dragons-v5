package main

import (
	"log"
	"net/http"
	"os"

	"codeberg.org/rpg/dungeon/dungeon"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Global game state (mono-user)
var gameState *GameState

func main() {
	log.Println("🎮 Initializing Dungeon RPG MCP Server...")

	// Initialize game with default dungeon
	gameState = initializeDefaultGame()
	if gameState == nil {
		log.Fatal("❌ Failed to initialize game state")
	}

	log.Println("✅ Game state initialized successfully")

	// Create MCP server
	s := server.NewMCPServer("dungeon-rpg", "1.0.0")

	// Register all tools
	registerTools(s)

	log.Println("✅ All MCP tools registered")

	// Setup HTTP server
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", healthCheckHandler)

	// MCP endpoint
	httpServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
	)
	mux.Handle("/mcp", httpServer)

	// Start server

	port := os.Getenv("MCP_HTTP_PORT")
	if port == "" {
		port = "9093"
	}

	log.Printf("🚀 MCP Dungeon Server starting on port %s", port)
	log.Printf("   Health endpoint: http://localhost:%s/health", port)
	log.Printf("   MCP endpoint: http://localhost:%s/mcp", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

// initializeDefaultGame loads the default dungeon and initializes game state
func initializeDefaultGame() *GameState {
	// Load default dungeon files
	dungeonFile := "./data/dungeon_generated.json"
	metadataFile := "./data/dungeon_metadata.json"

	// Load dungeon
	dng, err := dungeon.LoadFromFile(dungeonFile)
	if err != nil {
		log.Printf("⚠️  Error loading dungeon: %v", err)
		return nil
	}

	// Load metadata
	metadata, err := LoadMetadata(metadataFile)
	if err != nil {
		log.Printf("⚠️  Warning: Could not load metadata file: %v", err)
		metadata = &DungeonMetadata{Entities: []EntityMetadata{}}
	}

	// Create game state
	gs := NewGameState(dng)

	// Load entities
	loadEntitiesFromDungeon(gs, metadata)

	// Find and set player
	player := findPlayer(gs, metadata)
	if player == nil {
		log.Println("⚠️  Player not found in dungeon")
		return nil
	}
	gs.Player = player

	// Load riddle state
	gs.RiddleSolved = metadata.RiddleSolved

	log.Printf("✅ Dungeon loaded: %d rooms, player: %s (%s %s)",
		len(dng.GetRooms()), player.Name, player.Race, player.Class)

	return gs
}

// registerTools registers all MCP tools with the server
func registerTools(s *server.MCPServer) {
	// Status tools
	s.AddTool(createGetGameStatusTool(), handleGetGameStatus)
	s.AddTool(createSaveGameTool(), handleSaveGame)

	// Query tools
	s.AddTool(createGetCurrentRoomTool(), handleGetCurrentRoom)
	s.AddTool(createGetInventoryTool(), handleGetInventory)
	s.AddTool(createGetMapTool(), handleGetMap)
	s.AddTool(createGetHelpTool(), handleGetHelp)

	// Action tools
	s.AddTool(createMoveTool(), handleMove)
	s.AddTool(createCollectItemsTool(), handleCollectItems)
	s.AddTool(createDrinkPotionTool(), handleDrinkPotion)
	s.AddTool(createTalkToNPCsTool(), handleTalkToNPCs)
	s.AddTool(createAnswerRiddleTool(), handleAnswerRiddle)

	// Combat tools
	s.AddTool(createStartCombatTool(), handleStartCombat)
	s.AddTool(createAttackTool(), handleAttack)
}

// Tool definitions
func createGetGameStatusTool() mcp.Tool {
	return mcp.NewTool("get_game_status",
		mcp.WithDescription("Get current game status including combat state, riddle state, and player info"),
	)
}

func createSaveGameTool() mcp.Tool {
	return mcp.NewTool("save_game",
		mcp.WithDescription("Save the current game state to files"),
		mcp.WithString("filename_prefix",
			mcp.Description("Optional prefix for save filenames (default: dungeon_generated)"),
		),
	)
}

func createGetCurrentRoomTool() mcp.Tool {
	return mcp.NewTool("get_current_room",
		mcp.WithDescription("Get detailed information about the current room"),
	)
}

func createGetInventoryTool() mcp.Tool {
	return mcp.NewTool("get_inventory",
		mcp.WithDescription("Get player inventory and stats"),
	)
}

func createGetMapTool() mcp.Tool {
	return mcp.NewTool("get_map",
		mcp.WithDescription("Get ASCII art map of the entire dungeon"),
	)
}

func createGetHelpTool() mcp.Tool {
	return mcp.NewTool("get_help",
		mcp.WithDescription("Get comprehensive help about available commands, their usage, and current game context"),
	)
}

func createMoveTool() mcp.Tool {
	return mcp.NewTool("move",
		mcp.WithDescription("Move the player in a direction (north, south, east, west)"),
		mcp.WithString("direction",
			mcp.Required(),
			mcp.Description("Direction to move: north, south, east, or west"),
		),
	)
}

func createCollectItemsTool() mcp.Tool {
	return mcp.NewTool("collect_items",
		mcp.WithDescription("Collect all items (gold, potions) in the current room"),
	)
}

func createDrinkPotionTool() mcp.Tool {
	return mcp.NewTool("drink_potion",
		mcp.WithDescription("Drink a potion to restore health (costs 1 potion, restores 5 health)"),
	)
}

func createTalkToNPCsTool() mcp.Tool {
	return mcp.NewTool("talk_to_npcs",
		mcp.WithDescription("Talk to NPCs in the current room. Optional: specify npc_type (elf, dwarf, human, sphinx) to talk to a specific NPC type"),
		mcp.WithString("npc_type",
			mcp.Description("Optional: specific NPC type to talk to (elf, dwarf, human, sphinx). If not provided, talks to all NPCs in the room"),
		),
	)
}

func createAnswerRiddleTool() mcp.Tool {
	return mcp.NewTool("answer_riddle",
		mcp.WithDescription("Answer the Sphinx's riddle"),
		mcp.WithString("answer",
			mcp.Required(),
			mcp.Description("Your answer to the riddle"),
		),
	)
}

func createStartCombatTool() mcp.Tool {
	return mcp.NewTool("start_combat",
		mcp.WithDescription("Initiate combat with a monster in the current room"),
	)
}

func createAttackTool() mcp.Tool {
	return mcp.NewTool("attack",
		mcp.WithDescription("Attack the current enemy (must be in combat)"),
	)
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
