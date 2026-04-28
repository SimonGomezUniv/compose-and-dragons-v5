package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/snipwise/nova/nova-sdk/agents"
	"github.com/snipwise/nova/nova-sdk/agents/chat"
	"github.com/snipwise/nova/nova-sdk/agents/tools"
	"github.com/snipwise/nova/nova-sdk/mcptools"
	"github.com/snipwise/nova/nova-sdk/messages"
	"github.com/snipwise/nova/nova-sdk/messages/roles"
	"github.com/snipwise/nova/nova-sdk/models"
	"github.com/snipwise/nova/nova-sdk/toolbox/env"
	"github.com/snipwise/nova/nova-sdk/ui/display"
	"github.com/snipwise/nova/nova-sdk/ui/prompt"
)

/*
To execute this sample, make sure you have an MCP server running locally.

```
cd ../mcp-servers
docker compose up --build
```
*/

// MapResponse represents the JSON response from get_map tool
type MapResponse struct {
	Map string `json:"map"`
}

// displayColoredMap parses the JSON response and displays the map with ANSI colors
func displayColoredMap(jsonResult string) error {
	var mapResp MapResponse

	// Unmarshal the JSON string into our struct
	err := json.Unmarshal([]byte(jsonResult), &mapResp)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// The map string already contains ANSI escape sequences
	// We just need to print it directly to the terminal
	display.Separator()
	fmt.Println(mapResp.Map)
	display.Separator()

	return nil
}

func main() {
	ctx := context.Background()
	// Initialize MCP client
	// Url of the MCP Gateway
	mcpGatewayURL := env.GetEnvOrDefault("MCP_GATEWAY_URL", "http://localhost:9011")
	mcpClient, err := mcptools.NewStreamableHttpMCPClient(ctx, mcpGatewayURL)

	if err != nil {
		panic(err)
	}

	// Dungeaon MCP tools are now available via mcpClient.GetTools()
	dungeonTools := mcpClient.GetTools()

	// Print available tools
	for _, tool := range dungeonTools {
		println("Tool:", tool.Name, "-", tool.Description)
	}

	fmt.Println(strings.Repeat("=", 50))

	engineUrl := env.GetEnvOrDefault("ENGINE_BASE_URL", "http://localhost:11434/v1")
	toolsModelId := env.GetEnvOrDefault("TOOLS_MODEL", "qwen2:0.5b")

	buddyModelId := env.GetEnvOrDefault("BUDDY_MODEL", "qwen2:0.5b")

	dmBuddyAgent, err := chat.NewAgent(
		ctx,
		agents.Config{
			EngineURL: engineUrl,
			SystemInstructions: `
				You are Bud, a helpful AI assistant.
				You are an expert to create fancy reports about a text-based adventure game 
				based on the interactions of another AI agent called Zephyr with the game world.
				You will receive information in a JSON format about the actions taken by Zephyr and the results of those actions.
				Your goal is to create a comprehensive and engaging report of the adventure so far.
				The format of the report need to be very user friendly and engaging. 
				Use markdown formatting where appropriate to enhance readability.
				Add a touch of humor and personality to make the report more enjoyable to read.
				Add emojis where appropriate to enhance the tone of the report.
			`,
			KeepConversationHistory: false,
		},
		models.Config{
			Name:        buddyModelId,
			Temperature: models.Float64(0.8),
			//TopK: 	 models.Int(40),
			TopP: models.Float64(0.9),
		},
	)

	if err != nil {
		panic(err)
	}

	dmToolsAgent, err := tools.NewAgent(
		ctx,
		agents.Config{
			EngineURL: engineUrl,
			SystemInstructions: `
				You are Zephyr, a helpful AI assistant.
				You are playing a text-based adventure game where you can call functions to interact with the game world.
				Use the functions to explore the dungeon, move around, and get information about your surroundings.
				I repeat here the list of the available functions you can call:

				- answer_riddle - Answer the Sphinx's riddle
				- attack - Attack the current enemy (must be in combat)
				- collect_items - Collect all items (gold, potions) in the current room
				- drink_potion - Drink a potion to restore health (costs 1 potion, restores 5 health)
				- get_current_room - Get detailed information about the current room
				- get_game_status - Get current game status including combat state, riddle state, and player info
				- get_help - Get comprehensive help about available commands, their usage, and current game context
				- get_inventory - Get player inventory and stats
				- get_map - Get ASCII art map of the entire dungeon
				- move - Move the player in a direction (north, south, east, west)
				- save_game - Save the current game state to files
				- start_combat - Initiate combat with a monster in the current room
				- talk_to_npcs - Talk to NPCs in the current room
			`,
			//KeepConversationHistory: false,
		},
		models.Config{
			Name:              toolsModelId,
			Temperature:       models.Float64(0.0),
			ParallelToolCalls: models.Bool(true),
		},

		tools.WithMCPTools(dungeonTools),
	)

	if err != nil {
		panic(err)
	}

	// messages := []messages.Message{
	// 	{
	// 		Content: `
	// 		Get the map of the dungeon
	// 		Move to the south
	// 		Move to the east
	// 		Get information about the current room
	// 		Get the map of the dungeon
	// 		`,
	// 		Role: roles.User,
	// 	},
	// }

	executeFunction := func(functionName, arguments string) (string, error) {
		display.Colorf(display.ColorGreen, "🟢 Executing function: %s with arguments: %s\n", functionName, arguments)

		result, err := mcpClient.ExecToolWithString(functionName, arguments)
		if err != nil {
			display.Errorf("Failed when executing %s, %s: %v", functionName, arguments, err)
			return `{"error":"I think the Dungeon Master is drunk"}`, err
		}
		if functionName == "get_map" {
			// Display the colored map in the terminal
			if err := displayColoredMap(result); err != nil {
				return "", fmt.Errorf("failed to display map: %w", err)
			}
			return `{"map": "Displayed in terminal"}`, nil
		}
		return result, nil

		// TODO: add confirmation for running a tool

	}

	// TODO: add a chat agent to transform the result into a more friendly format
	// call it a first time to "warm up"

	for {

		markdownParser := display.NewMarkdownChunkParser()

		input := prompt.NewWithColor("🤖 Ask me something?")
		question, err := input.Run()

		if err != nil {
			display.Errorf("failed to get input: %v", err)
			return
		}
		if strings.HasPrefix(question, "/bye") {
			display.Infof("👋 Goodbye!")
			break
		}

		promptMessages := []messages.Message{
			{
				Content: question,
				Role:    roles.User,
			},
		}

		responses, err := dmToolsAgent.DetectParallelToolCalls(promptMessages, executeFunction)
		if err != nil {
			display.Errorf("Error calling tools agent: %v", err)
			continue
		}

		fmt.Println(responses.Results)

		if len(responses.Results) == 0 || responses.Results[0] != `{"map": "Displayed in terminal"}` {

			allContents := strings.Join(responses.Results, ",")

			display.Separator()

			//var responseBuilder strings.Builder
			//var reasoningBuilder strings.Builder

			// Chat with streaming and reasoning - no OpenAI types exposed
			_, err = dmBuddyAgent.GenerateStreamCompletionWithReasoning(
				[]messages.Message{
					{Role: roles.System, Content: allContents},
					{Role: roles.User, Content: "Make a report from the above data."},
				},
				func(reasoningChunk string, finishReason string) error {
					//reasoningBuilder.WriteString(reasoningChunk)
					display.Color(reasoningChunk, display.ColorBlue)
					if finishReason != "" {
						display.NewLine()
						display.KeyValue("Reasoning finish reason", finishReason)
						display.Separator()
					}
					return nil
				},
				func(responseChunk string, finishReason string) error {
					// responseBuilder.WriteString(responseChunk)
					// display.Color(responseChunk, display.ColorGreen)

					// Use markdown chunk parser for colorized streaming output
					if responseChunk != "" {
						display.MarkdownChunk(markdownParser, responseChunk)
					}
					if finishReason == "stop" {
						markdownParser.Flush()
						markdownParser.Reset()
						//markdownParser.Flush()
						display.NewLine()
						//fmt.Println()
					}


					if finishReason != "" {
						display.NewLine()
						display.KeyValue("Response finish reason", finishReason)
					}
					return nil
				},
			)

			if err != nil {
				display.Errorf("Error calling buddy agent: %v", err)
				continue
			}

		}

	}

}
