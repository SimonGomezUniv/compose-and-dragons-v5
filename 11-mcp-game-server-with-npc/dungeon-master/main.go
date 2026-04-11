package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/snipwise/nova/nova-sdk/agents"
	"github.com/snipwise/nova/nova-sdk/agents/chat"
	"github.com/snipwise/nova/nova-sdk/agents/remote"
	"github.com/snipwise/nova/nova-sdk/agents/tools"
	"github.com/snipwise/nova/nova-sdk/mcptools"
	"github.com/snipwise/nova/nova-sdk/messages"
	"github.com/snipwise/nova/nova-sdk/messages/roles"
	"github.com/snipwise/nova/nova-sdk/models"
	"github.com/snipwise/nova/nova-sdk/toolbox/conversion"
	"github.com/snipwise/nova/nova-sdk/toolbox/env"
	"github.com/snipwise/nova/nova-sdk/toolbox/files"
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
	// [MCP CLIENT] Initialize
	mcpGatewayURL := env.GetEnvOrDefault("MCP_GATEWAY_URL", "http://localhost:9011")
	mcpClient, err := mcptools.NewStreamableHttpMCPClient(ctx, mcpGatewayURL)

	if err != nil {
		display.Errorf("failed to create MCP client: %v", err)
		return
		//panic(err)
	}

	// Dungeaon [MCP TOOLS] are now available via mcpClient.GetTools()
	dungeonTools := mcpClient.GetTools()

	// Print available tools
	display.Title("Tools:")
	for _, tool := range dungeonTools {
		//println("Tool:", tool.Name, "-", tool.Description)
		//display.NewLine()
		display.KeyValue(tool.Name, tool.Description)

	}

	display.Separator()
	//fmt.Println(strings.Repeat("=", 50))

	engineUrl := env.GetEnvOrDefault("ENGINE_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	toolsModelId := env.GetEnvOrDefault("TOOLS_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")
	buddyModelId := env.GetEnvOrDefault("BUDDY_MODEL", "huggingface.co/menlo/lucy-gguf:q4_k_m")

	// AGENT: Dungeon Master Buddy Agent
	// dm.buddy.system.instructions.md
	dmBuddyAgentSystemInstructions, err := files.ReadTextFile("./docs/dm.buddy.system.instructions.md")
	if err != nil {
		display.Errorf("failed to read dmBuddyAgent system instructions file: %v", err)
		return
		//panic(err)
	}

	dmBuddyAgent, err := chat.NewAgent(
		ctx,
		agents.Config{
			EngineURL:               engineUrl,
			SystemInstructions:      dmBuddyAgentSystemInstructions,
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
		display.Errorf("failed to create dmBuddyAgent: %v", err)
		return
		//panic(err)
	}

	// AGENT: Dungeon Master Tools Agent
	// dm.tools.system.instructions.md
	dmToolsAgentSystemInstructions, err := files.ReadTextFile("./docs/dm.tools.system.instructions.md")
	if err != nil {
		display.Errorf("failed to read dmToolsAgent system instructions file: %v", err)
		return
		//panic(err)
	}

	dmToolsAgent, err := tools.NewAgent(
		ctx,
		agents.Config{
			EngineURL:          engineUrl,
			SystemInstructions: dmToolsAgentSystemInstructions,
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
		display.Errorf("failed to create dmToolsAgent: %v", err)
		return
		//panic(err)
	}

	// [REMOTE] AGENTS:
	dwarfAgentURL := env.GetEnvOrDefault("DWARF_AGENT_URL", "http://localhost:9091")
	elfAgentURL := env.GetEnvOrDefault("ELF_AGENT_URL", "http://localhost:9092")
	humanAgentURL := env.GetEnvOrDefault("HUMAN_AGENT_URL", "http://localhost:9093")
	sphinxAgentURL := env.GetEnvOrDefault("SPHINX_AGENT_URL", "http://localhost:9094")

	dwarfAgent, err := remote.NewAgent(
		ctx,
		agents.Config{
			Name: "Interactive Remote Dwarf NPC Client",
		},
		dwarfAgentURL,
	)
	if err != nil {
		display.Errorf("failed to create dwarfAgent: %v", err)
		return
	}
	display.Colorf(display.ColorCyan, "🌐 Connected to remote agent at %s\n", dwarfAgentURL)
	display.Colorf(display.ColorCyan, "Agent: %s\n", dwarfAgent.GetName())
	display.Colorf(display.ColorCyan, "Model: %s\n\n", dwarfAgent.GetModelID())

	elfAgent, err := remote.NewAgent(
		ctx,
		agents.Config{
			Name: "Interactive Remote Elf NPC Client",
		},
		elfAgentURL,
	)
	if err != nil {
		display.Errorf("failed to create elfAgent: %v", err)
		return
	}
	display.Colorf(display.ColorCyan, "🌐 Connected to remote agent at %s\n", elfAgentURL)
	display.Colorf(display.ColorCyan, "Agent: %s\n", elfAgent.GetName())
	display.Colorf(display.ColorCyan, "Model: %s\n\n", elfAgent.GetModelID())

	humanAgent, err := remote.NewAgent(
		ctx,
		agents.Config{
			Name: "Interactive Remote Human NPC Client",
		},
		humanAgentURL,
	)
	if err != nil {
		display.Errorf("failed to create humanAgent: %v", err)
		return
	}
	display.Colorf(display.ColorCyan, "🌐 Connected to remote agent at %s\n", humanAgentURL)
	display.Colorf(display.ColorCyan, "Agent: %s\n", humanAgent.GetName())
	display.Colorf(display.ColorCyan, "Model: %s\n\n", humanAgent.GetModelID())

	sphinxAgent, err := remote.NewAgent(
		ctx,
		agents.Config{
			Name: "Interactive Remote Sphinx NPC Client",
		},
		sphinxAgentURL,
	)
	if err != nil {
		display.Errorf("failed to create sphinxAgent: %v", err)
		return
	}
	display.Colorf(display.ColorCyan, "🌐 Connected to remote agent at %s\n", sphinxAgentURL)
	display.Colorf(display.ColorCyan, "Agent: %s\n", sphinxAgent.GetName())
	display.Colorf(display.ColorCyan, "Model: %s\n\n", sphinxAgent.GetModelID())

	npcRemoteAgents := map[string]*remote.Agent{
		"dwarf":  dwarfAgent,
		"elf":    elfAgent,
		"human":  humanAgent,
		"sphinx": sphinxAgent,
	}

	display.Separator()

	// -------------------------------------------------
	// FUNCTION TO EXECUTE TOOL CALLS
	// ------------------------------------------------
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

	// -------------------------------------------------
	// INTERACTIVE GAME LOOP
	// ------------------------------------------------
	for {

		markdownParser := display.NewMarkdownChunkParser()

		input := prompt.NewWithColor("🤖 Ask me something?")
		question, err := input.RunWithEdit()

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
			display.Errorf("err: %v", err)
			continue
		}

		//fmt.Println(responses.Results)

		if len(responses.Results) == 0 {
			display.Errorf("no results from the dungeon tools agent")
			continue
		}

		resultMap, err := conversion.JsonStringToMap(responses.Results[0])
		if err != nil {
			display.Errorf("err: %v", err)
			continue
		}
		display.Color(fmt.Sprintf("%+v\n", resultMap), display.ColorGray)

		createStoryFromContent := func() {
			// Aggregate all results into a single string
			allContents := strings.Join(responses.Results, ",")

			display.Separator()

			// NOTE: dmBuddyAgent generates a fancy report with the results of the tools calls
			// Chat with streaming and reasoning
			_, err = dmBuddyAgent.GenerateStreamCompletionWithReasoning(
				[]messages.Message{
					{Role: roles.System, Content: allContents},
					{Role: roles.User, Content: "Make a report from the above data."},
				},
				// [REASONING]
				func(reasoningChunk string, finishReason string) error {
					display.Color(reasoningChunk, display.ColorBlue)
					if finishReason != "" {
						display.NewLine()
						display.KeyValue("Reasoning finish reason", finishReason)
						display.Separator()
					}
					return nil
				},
				// [RESPONSE]
				func(responseChunk string, finishReason string) error {
					// Use markdown chunk parser for colorized streaming output
					if responseChunk != "" {
						display.MarkdownChunk(markdownParser, responseChunk)
					}
					if finishReason == "stop" {
						markdownParser.Flush()
						markdownParser.Reset()
						display.NewLine()
					}

					if finishReason != "" {
						display.NewLine()
						display.KeyValue("Response finish reason", finishReason)
					}
					return nil
				},
			)

			if err != nil {
				display.Errorf("err: %v", err)
			}
		}

		switch {

		case resultMap["map"] == "Displayed in terminal":
			// Map already displayed
			display.Separator()

		case resultMap["conversations"] != nil:

			if success, ok := resultMap["success"].(bool); ok && success {
				conversations := resultMap["conversations"].([]interface{})
				if len(conversations) > 0 {
					conv := conversations[0].(map[string]interface{})

					kind := conv["npc_type"].(string)
					npcId := conv["npc_id"].(string)
					dialogue := conv["dialogue"].(string)
					display.Colorf(display.ColorBrightRed, "%s NPC: %s - %s\n", kind, npcId, dialogue)

					selectedNPCAgent, exists := npcRemoteAgents[kind]
					if !exists {
						display.Errorf("No remote agent found for NPC type: %s", kind)
						continue
					}

					for {

						input := prompt.NewWithColor("😃["+npcId+"] Ask me something?").
							SetMessageColor(prompt.ColorBrightCyan).
							SetInputColor(prompt.ColorBrightWhite)

						question, err := input.Run()
						if err != nil {
							log.Fatal(err)
						}

						if strings.HasPrefix(question, "/bye") {
							fmt.Println("Goodbye!")
							break
						}
						display.NewLine()

						_, err = selectedNPCAgent.GenerateStreamCompletion(
							[]messages.Message{
								{
									Role:    roles.User,
									Content: question,
								},
							},
							func(chunk string, finishReason string) error {

								// Use markdown chunk parser for colorized streaming output
								if chunk != "" {
									fmt.Print(chunk)
								}
								if finishReason == "stop" {
									fmt.Println()
								}
								return nil
							},
						)
						if err != nil {
							display.Errorf("failed to chat with NPC agent: %v", err)
							return
							//panic(err)
						}
						display.NewLine()
						// display.Separator()
						// display.KeyValue("Finish reason", result.FinishReason)
						// display.KeyValue("Context size", fmt.Sprintf("%d characters", selectedNPCAgent.GetContextSize()))
						// display.Separator()

					}
					// if conv["npc_type"] == "elf" {
					// 	npcId := conv["npc_id"].(string)
					// 	dialogue := conv["dialogue"].(string)
					// 	display.Colorf(display.ColorBrightRed, "Elf NPC: %s - %s\n", npcId, dialogue)
					// }
					// if conv["npc_type"] == "dwarf" {
					// 	npcId := conv["npc_id"].(string)
					// 	dialogue := conv["dialogue"].(string)
					// 	display.Colorf(display.ColorBrightRed, "Dwarf NPC: %s - %s\n", npcId, dialogue)
					// }
					// if conv["npc_type"] == "human" {
					// 	npcId := conv["npc_id"].(string)
					// 	dialogue := conv["dialogue"].(string)
					// 	display.Colorf(display.ColorBrightRed, "Human NPC: %s - %s\n", npcId, dialogue)
					// }
					// if conv["npc_type"] == "sphinx" {
					// 	npcId := conv["npc_id"].(string)
					// 	dialogue := conv["dialogue"].(string)
					// 	display.Colorf(display.ColorBrightRed, "Sphinx NPC: %s - %s\n", npcId, dialogue)
					// }

				}
			} else {
				createStoryFromContent()
			}

		case responses.Results[0] == ``:
			display.Warningf("⚠️ No response from tools")
		default:
			createStoryFromContent()

		}
	}

}
