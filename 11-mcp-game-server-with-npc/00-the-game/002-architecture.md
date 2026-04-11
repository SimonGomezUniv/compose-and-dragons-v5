<style>
.dodgerblue {
  color: dodgerblue;
}
.indianred {
  color: indianred;
}
</style>
# 🏰 🐙🐲 Compose & Dragons


## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                        DUNGEON MASTER (Go)                          │
│  - Main game client application                                     │
│  - Two AI agents: dmToolsAgent & dmBuddyAgent                       │
│  - User interaction loop                                            │
└───────────┬──────────────────────────────────────────┬──────────────┘
            │                                          │
            │ HTTP/MCP calls                           │ HTTP calls
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






