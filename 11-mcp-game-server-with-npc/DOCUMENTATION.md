# Comprehensive Documentation: Dungeon Master Game System

## System Overview

This is a **multi-agent D&D-style dungeon crawler game** that combines several advanced technologies:
- **MCP (Model Context Protocol)** for tool orchestration
- **AI Agents** for game interaction and NPC roleplay
- **RAG (Retrieval Augmented Generation)** for NPC character knowledge
- **Docker Compose** for service orchestration

The system consists of **8 interconnected Docker services** that work together to create an interactive text-based RPG experience.

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                        DUNGEON MASTER (Go)                          │
│  - Main game client application                                     │
│  - Two AI agents: dmToolsAgent & dmBuddyAgent                       │
│  - User interaction loop                                            │
└───────────┬──────────────────────────────────────────┬──────────────┘
            │                                          │
            │ HTTP/MCP calls                          │ HTTP calls
            ↓                                          ↓
┌────────────────────────┐              ┌──────────────────────────────┐
│   MCP GATEWAY          │              │  NPC Remote Agents (4 total) │
│  - Security layer      │              │  ┌────────────────────────┐  │
│  - Tool filtering      │              │  │  dwarf-warrior:8080    │  │
│  - Proxy to MCP server │              │  │  elf-mage:8080         │  │
└──────────┬─────────────┘              │  │  human-rogue:8080      │  │
           │                            │  │  sphinx-boss:8080      │  │
           │ Proxied MCP calls          │  └────────────────────────┘  │
           ↓                            │  - RAG-enabled chat agents   │
┌────────────────────────┐              │  - Character sheet knowledge │
│  DUNGEON MCP SERVER    │              │  - Context compression       │
│  - Game state manager  │              └──────────────────────────────┘
│  - 14 MCP tools        │
│  - Dungeon logic       │
│  - Combat system       │
└────────────────────────┘

Additional Services:
┌────────────────────────┐
│  MCP INSPECTOR         │  - Debugging/inspection UI
└────────────────────────┘
```

---

## Service Breakdown

### 1. **dungeon-master** (Main Application)

**File**: `dungeon-master/main.go`

**Purpose**: This is the **game client** that orchestrates everything. It's the main entry point where the user interacts with the game.

**Key Components**:

#### Two AI Agents:

1. **dmToolsAgent** (lines 128-148)
   - **Model**: `jan-nano` or `qwen2.5` (small, fast model)
   - **Purpose**: Executes MCP tools to interact with the game
   - **Capabilities**:
     - Detects which game actions to take based on user input
     - Makes parallel tool calls
     - Temperature: 0.0 (deterministic)
   - **System Instructions**: `/app/docs/dm.tools.system.instructions.md`
   - **Available Tools**: All 14 MCP tools from the dungeon server

2. **dmBuddyAgent** (lines 98-117)
   - **Model**: `lucy` (narrative model)
   - **Purpose**: Creates engaging narrative reports of game events
   - **Capabilities**:
     - Converts JSON game results into engaging stories
     - Reasoning capability for enhanced storytelling
     - Temperature: 0.8 (creative)
   - **System Instructions**: `/app/docs/dm.buddy.system.instructions.md`

#### Remote NPC Agents (lines 150-221):
- **4 remote agent connections** to NPC services
- Each NPC runs as a separate HTTP service
- Agents: `dwarfAgent`, `elfAgent`, `humanAgent`, `sphinxAgent`

#### Main Game Loop (lines 253-448):

The game loop follows this flow:

1. **User Input** (line 257-258): Accepts natural language commands
2. **Tool Detection** (line 276): `dmToolsAgent.DetectParallelToolCalls()` determines which MCP tools to call
3. **Tool Execution** (lines 228-248): `executeFunction()` calls the MCP server via gateway
4. **Response Processing** (lines 289-447): Handles different result types:
   - **Map display** (line 346): Special handling for colored ASCII maps
   - **NPC conversations** (lines 350-440): Initiates interactive chat with remote NPCs
   - **General results** (lines 442-446): Creates narrative with dmBuddyAgent

#### NPC Interaction Flow (lines 350-440):

When talking to NPCs:
1. Extracts NPC type and ID from tool result
2. Selects corresponding remote agent
3. **Enters sub-loop** for conversation with that specific NPC
4. Streams responses from remote agent
5. User can type `/bye` to exit conversation

**Environment Variables**:
- `MCP_GATEWAY_URL`: URL to MCP gateway (default: `http://mcp-gateway:9011/mcp`)
- `DWARF_AGENT_URL`: `http://dwarf-warrior:8080`
- `ELF_AGENT_URL`: `http://elf-mage:8080`
- `HUMAN_AGENT_URL`: `http://human-rogue:8080`
- `SPHINX_AGENT_URL`: `http://sphinx-end-of-level-boss:8080`
- `ENGINE_BASE_URL`: AI inference engine
- `TOOLS_MODEL`: Model for tools agent
- `BUDDY_MODEL`: Model for buddy agent

---

### 2. **mcp-gateway** (Security Layer)

**Service**: Docker MCP Gateway v2

**Purpose**: Acts as a **security proxy** between the dungeon-master and the MCP server.

**Configuration** (compose.yml lines 41-62):
- **Port**: 9011
- **Transport**: Streaming HTTP
- **Tool Filtering**: Only allows specific tools (line 54):
  ```
  get_map, answer_riddle, attack, collect_items, drink_potion,
  move, save_game, start_combat, talk_to_npcs, get_current_room,
  get_game_status, get_inventory, get_help
  ```

**Why it exists**:
- Provides centralized access control
- Can whitelist/blacklist specific tools
- Adds logging and monitoring capability
- Decouples client from server implementation

**Dependency**: Waits for `dungeon-mcp-server` to be healthy before starting

---

### 3. **dungeon-mcp-server** (Game Engine)

**Files**:
- Main: `dungeon-mcp-server/main.go`
- Handlers: `dungeon-mcp-server/handlers.go`
- Game State: `dungeon-mcp-server/game.go`

**Purpose**: This is the **core game engine** that manages all game state and logic.

#### Game State (game.go):

The `GameState` struct contains:
- **Dungeon**: The dungeon structure (rooms, connections)
- **Player**: Player character with stats and inventory
- **Collectibles**: Map of gold coins and potions
- **Monsters**: Map of enemies (skeletons, goblins, vampires)
- **NPCs**: Map of friendly characters (elves, dwarves, humans, sphinx)
- **InCombat**: Combat state flag
- **CurrentEnemy**: Pointer to active enemy
- **RiddleSolved**: Sphinx riddle completion status

#### 14 MCP Tools:

**Status Tools**:
- `get_game_status`: Returns combat state, riddle state, player info
- `save_game`: Saves current state to timestamped files

**Query Tools**:
- `get_current_room`: Returns detailed room info (exits, monsters, NPCs, items)
- `get_inventory`: Returns player stats and inventory
- `get_map`: Returns ASCII art map with ANSI colors
- `get_help`: Context-aware help system

**Action Tools**:
- `move`: Move player (north/south/east/west), also flees combat
- `collect_items`: Picks up gold and potions
- `drink_potion`: Restores 5 HP, costs 1 potion
- `talk_to_npcs`: Returns NPC dialogues (with optional NPC type filter)
- `answer_riddle`: Submit answer to Sphinx's riddle

**Combat Tools**:
- `start_combat`: Initiates combat with monster
- `attack`: Execute combat round

#### Key Handler Logic:

**handleTalkToNPCs**:
- Accepts optional `npc_type` parameter to filter specific NPCs
- Returns JSON with NPC type, ID, and dialogue
- For Sphinx: dialogue changes based on `RiddleSolved` state

**handleMove**:
- Automatically flees combat if in combat
- Validates direction and connection existence
- Updates player position
- Returns new room information

**handleStartCombat**:
- Prevents double combat initiation
- Finds first monster in room
- Sets `InCombat = true` and `CurrentEnemy`
- Returns monster stats

**handleAttack**:
- Calls `PerformCombatRound()` to execute combat logic
- Returns combat result (victory/defeat)

#### Combat System:

**Combat Mechanics**:
- Dice-based: 3d6 + Strength modifier
- Damage: `(winner_total - loser_total + 1) / 2`
- Victory: Monster removed from game
- Defeat: Game over (exits process)

**Health Endpoint**:
- `/health` endpoint for Docker healthcheck
- Returns 200 OK when server is ready

---

### 4. **NPC Services** (4 instances)

**Files**:
- Main: `non-player-characters/main.go`
- Server: `non-player-characters/interactive.chat.go`
- RAG: `non-player-characters/rag.agent.go`

**Purpose**: Each NPC is a **RAG-enabled conversational AI agent** that roleplays a D&D character.

#### Four NPC Instances:

1. **dwarf-warrior** (compose.yml lines 152-173)
   - Character: Thorin Sturgeshield (male dwarf warrior)
   - Port: 9091
   - Sheet: `male-dwarf-warrior.md`

2. **elf-mage** (compose.yml lines 176-197)
   - Character: Female elf mage
   - Port: 9092
   - Sheet: `female-elf-mage.md`

3. **human-rogue** (compose.yml lines 201-221)
   - Character: Male human rogue
   - Port: 9093
   - Sheet: `male-human-rogue.md`

4. **sphinx-end-of-level-boss** (compose.yml lines 226-247)
   - Character: Female sphinx (end boss)
   - Port: 9094
   - Sheet: `female-sphinx-boss.md`
   - Special instructions: Riddle logic

#### NPC Agent Architecture:

**Three Sub-Agents**:

1. **Metadata Extractor Agent**
   - Extracts keywords from character sheet
   - Used for RAG document chunking

2. **RAG Agent**
   - Loads pre-built vector store from JSON
   - Provides character sheet context via similarity search
   - Store path: `./store/{character-name}.json`

3. **Compressor Agent**
   - Compresses conversation history when context grows too large
   - Maintains conversation coherence with limited context

#### NPC Server Configuration:

**Server Agent Setup**:
- **Port**: 8080 (internal to container)
- **Model**: `qwen2.5:1.5B-F16` (small, fast conversational model)
- **Temperature**: 0.9 (high creativity for roleplay)
- **Keep History**: True (maintains conversation context)
- **RAG Integration**:
  - Similarity threshold: 0.4
  - Max similarities: 7 results
- **Compression**:
  - Context size limit: 80,000 characters
  - Auto-compresses when exceeded

**System Instructions**:
- Dynamically formatted with character data:
  - First name, family name, gender, race, class
  - Secret word (used in Sphinx riddle puzzle)
- Instructions emphasize staying in character

**Character Sheet Structure**:
```markdown
# CHARACTER SHEET
- Name and Title
- Age
- Background Story
- Personality and Character Traits
- Occupation
- Abilities and Skills
- Physical Appearance
- Clothing
- Food Preferences
- Favorite Quote
```

#### Special: Sphinx NPC:

The Sphinx has **unique behavior**:
- If player provides secret words: "Shadowfell Starforge Runewardens"
  - Grants passage
  - Displays celebration emoji: 🎉🎉🎉
- Otherwise: Denies passage and poses riddles
- Speaks in archaic, mystical tone with rhetorical questions

---

### 5. **mcp-inspector** (Debugging Tool)

**Service**: MCP Inspector (official debugging UI)

**Purpose**: Web-based tool for **inspecting and testing MCP tools** directly.

**Configuration** (compose.yml lines 89-100):
- **Ports**:
  - 6274: Web UI
  - 6277: Inspector API
- **Connection**: Points to `http://mcp-gateway:9011/mcp`
- **Access**: View logs to get connection URL

**Use Case**: Debug MCP tools without running the full game client

---

## Data Flow Examples

### Example 1: User Moves North

```
1. User types: "go north"
   Location: dungeon-master/main.go line 258

2. dmToolsAgent analyzes input
   Location: dungeon-master/main.go line 276
   Output: Determines to call "move" tool with direction="north"

3. executeFunction() calls MCP gateway
   Location: dungeon-master/main.go line 231
   HTTP POST: http://mcp-gateway:9011/mcp
   Tool: "move"
   Args: {"direction": "north"}

4. MCP Gateway forwards to dungeon-mcp-server
   Location: dungeon-mcp-server receives at /mcp endpoint

5. handleMove() processes request
   Location: dungeon-mcp-server/handlers.go line 317
   - Validates direction
   - Checks for connection
   - Updates player position in GameState
   - Returns new room info as JSON

6. Gateway returns result to dungeon-master

7. dmBuddyAgent creates narrative
   Location: dungeon-master/main.go line 304
   Converts: {"success": true, "new_room": {...}}
   Into: "You moved north into the Ancient Hall. The air grows colder..."

8. Result streamed to user with Markdown formatting
   Location: dungeon-master/main.go line 323
```

### Example 2: Talking to an Elf NPC

```
1. User types: "talk to the elf"

2. dmToolsAgent calls "talk_to_npcs" with npc_type="elf"
   Location: dungeon-master/main.go line 276

3. handleTalkToNPCs() returns dialogue
   Location: dungeon-mcp-server/handlers.go line 452
   Returns: {
     "success": true,
     "conversations": [{
       "npc_type": "elf",
       "npc_id": "elf_001",
       "dialogue": "Greetings, traveler! ..."
     }]
   }

4. dungeon-master detects NPC conversation
   Location: dungeon-master/main.go line 350
   - Extracts NPC type: "elf"
   - Extracts NPC ID: "elf_001"

5. Selects elfAgent from npcRemoteAgents map
   Location: dungeon-master/main.go line 362

6. Enters conversation loop
   Location: dungeon-master/main.go line 368
   - Displays special prompt: "😃[elf_001] Ask me something?"

7. User asks: "What do you know about the dungeon?"

8. Sends HTTP request to elf-mage:8080
   Location: elf-mage container (remote.Agent)

9. elf-mage NPC server processes request
   Location: non-player-characters/interactive.chat.go line 85
   - RAG agent retrieves relevant character sheet info
   - Server agent generates roleplay response
   - Uses character personality and knowledge

10. Response streamed back to dungeon-master
    Location: dungeon-master/main.go line 385

11. User can continue conversation or type "/bye" to exit
    Location: dungeon-master/main.go line 379
```

### Example 3: Combat Sequence

```
1. User: "attack the skeleton"

2. dmToolsAgent calls "start_combat"
   - handleStartCombat() sets InCombat=true
   - Sets CurrentEnemy to skeleton object

3. User: "fight"

4. dmToolsAgent calls "attack"
   - handleAttack() calls PerformCombatRound()
   - Rolls dice: Player 3d6+strength vs Monster 3d6+strength
   - Calculates damage
   - Updates health

5. If monster dies:
   - Removed from gs.Monsters map
   - Removed from dungeon room
   - InCombat set to false

6. dmBuddyAgent narrates combat:
   "Your blade strikes true! The skeleton collapses into a pile of bones..."
```

---

## Key Design Patterns

### 1. **Agent Separation of Concerns**

- **dmToolsAgent**: Technical execution (deterministic, temp=0.0)
- **dmBuddyAgent**: Creative narration (creative, temp=0.8)
- **NPC Agents**: Character roleplay (creative, temp=0.9)

This separation ensures:
- Tool calls are reliable and consistent
- Narratives are engaging and varied
- NPCs have distinct personalities

### 2. **MCP Gateway Pattern**

Benefits:
- **Security**: Whitelist approved tools
- **Decoupling**: Client doesn't know server details
- **Scalability**: Can add more MCP servers
- **Monitoring**: Central point for logging

### 3. **Remote Agent Pattern**

Each NPC is a **remote HTTP service** rather than in-process:
- **Isolation**: NPC crashes don't affect main game
- **Scalability**: NPCs can run on different machines
- **Flexibility**: Different models/configurations per NPC
- **Resource Management**: Heavy RAG operations isolated

### 4. **RAG for Character Consistency**

NPCs use RAG to:
- Access detailed character sheets (background, personality, skills)
- Provide consistent responses across conversations
- Reference specific character knowledge
- Maintain roleplay authenticity

### 5. **Context Compression**

Long conversations are compressed to:
- Prevent context window overflow
- Maintain conversation history
- Preserve important details
- Keep inference fast

---

## Configuration Management

### Docker Compose Features:

**Common Variables**:
```yaml
NOVA_LOG_LEVEL: INFO
DOCS_PATH: ./docs
SHEETS_PATH: ./sheets
STORE_PATH: ./store
WARM_UP: true
```

**Common Models**:
- RAG embedding model
- NPC conversation model
- Context compressor model
- Metadata extraction model

**Model Definitions**:
```yaml
dungeon-tools-model: jan-nano (16K context)
dungeon-master-buddy: lucy (16K context)
npc-model: qwen2.5:1.5B (16K context)
rag-model: embeddinggemma
compressor-model: qwen2.5:0.5B
metadata-model: jan-nano
```

---

## How to Extend the System

### Add a New NPC:

1. Create character sheet in `non-player-characters/sheets/new-character.md`
2. Generate vector store (RAG embeddings)
3. Add service definition in `compose.yml`:
   ```yaml
   new-npc:
     build: ./non-player-characters
     ports: [9095:8080]
     environment:
       SHEET_FILE_NAME: "new-character.md"
   ```
4. Add remote agent connection in `dungeon-master/main.go`:
   ```go
   newAgent, err := remote.NewAgent(ctx, agents.Config{...}, newAgentURL)
   npcRemoteAgents["new_type"] = newAgent
   ```
5. Update `talk_to_npcs` handler to recognize new NPC type

### Add a New MCP Tool:

1. Define tool in `dungeon-mcp-server/main.go`:
   ```go
   func createNewTool() mcp.Tool {
       return mcp.NewTool("new_tool", ...)
   }
   ```
2. Register tool: `s.AddTool(createNewTool(), handleNewTool)`
3. Implement handler in `handlers.go`:
   ```go
   func handleNewTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
       // Implementation
   }
   ```
4. Add tool to gateway whitelist in `compose.yml` line 54

### Add a New Agent to dungeon-master:

1. Load system instructions
2. Create agent with `NewAgent()` or appropriate agent type
3. Use in game loop for specific processing tasks

---

## Performance Considerations

### Model Selection:
- **Tools Agent**: Small model (jan-nano) for speed and reliability
- **Buddy Agent**: Medium model (lucy) for narrative quality
- **NPC Agents**: Small model (qwen2.5:1.5B) for quick responses
- **Embedding**: Specialized embedding model (gemma)

### Context Management:
- RAG reduces context size by retrieving only relevant character info
- Compression agent prevents context overflow in long conversations
- Tools agent doesn't keep history (stateless tool execution)
- Buddy agent doesn't keep history (one-shot narrative generation)
- NPC agents keep history (ongoing conversations)

### Parallelization:
- Multiple NPCs run simultaneously as separate services
- Combat calculations are synchronous (game state consistency)
- Tool calls can be parallel (`ParallelToolCalls: true`)

---

## Security & State Management

### State Persistence:
- Game state saved via `save_game` tool
- Creates timestamped files: `dungeon_generated_{timestamp}.json`
- Metadata saved separately: `dungeon_metadata_{timestamp}.json`
- RAG stores pre-built and loaded at startup

### Error Handling:
- MCP tool errors return JSON with error messages
- Agent errors logged and displayed to user
- Combat failure (player death) exits program
- Service dependencies managed by Docker Compose healthchecks

### Access Control:
- MCP Gateway filters allowed tools
- No authentication (single-player local game)
- Each service isolated in Docker network

---

## Summary

This system demonstrates **advanced multi-agent orchestration** for game development:

1. **Separation of Concerns**: Different agents handle different cognitive tasks
2. **MCP Protocol**: Standardized tool calling for game actions
3. **RAG Technology**: Knowledge-grounded NPC conversations
4. **Microservices**: Each NPC is an independent service
5. **Context Management**: Compression and RAG for long conversations
6. **Docker Orchestration**: 8 services working in harmony

The architecture is highly **modular and extensible**, making it easy to add new NPCs, game mechanics, or AI capabilities.
