package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/snipwise/nova/nova-sdk/agents"
	"github.com/snipwise/nova/nova-sdk/agents/chat"
	"github.com/snipwise/nova/nova-sdk/messages"
	"github.com/snipwise/nova/nova-sdk/messages/roles"
	"github.com/snipwise/nova/nova-sdk/models"
	"github.com/snipwise/nova/nova-sdk/toolbox/files"
	"github.com/snipwise/nova/nova-sdk/ui/display"
	"github.com/snipwise/nova/nova-sdk/ui/prompt"
)

// startInteractiveChat starts an interactive roleplay session with the NPC
func startInteractiveChat(ctx context.Context, engineURL, npcAgentModelID, sheetFilePath string) error {

	characterSheetContent, npc, err := loadNPCSheetFromFile(sheetFilePath)
	if err != nil {
		display.Errorf("failed to load NPC character sheet: %v", err)
		return err
	}

	// === INTERACTIVE NPC CHAT ===
	display.Infof("💬 Starting interactive chat with NPC...")
	display.Separator()

	roleplaySystemInstructionsTemplate, err := files.ReadTextFile("./dnd.chat.system.instructions.md")
	if err != nil {
		display.Errorf("failed to read roleplay system instructions file: %v", err)
		return err
	}

	// Create system instructions for the roleplay agent
	roleplaySystemInstructions := fmt.Sprintf(roleplaySystemInstructionsTemplate,
		npc.FirstName,
		npc.FamilyName,
		npc.Gender,
		npc.Race,
		npc.Class,
		npc.SecretWord,
		characterSheetContent,
		npc.FirstName, // Always speak in first person as %s
	)

	//display.Colorln(roleplaySystemInstructions, display.ColorBrightRed)

	// === CREATE THE ROLEPLAY CHAT AGENT ===
	npcAgent, err := chat.NewAgent(
		ctx,
		agents.Config{
			Name:               "npc-roleplay-agent",
			EngineURL:          engineURL,
			SystemInstructions: roleplaySystemInstructions,
			KeepConversationHistory: true, // IMPORTANT: Keep history for better context 
		},
		models.Config{
			Name:        npcAgentModelID,
			Temperature: models.Float64(0.9), // High creativity for roleplay
			TopP:        models.Float64(0.95),
		},
	)
	if err != nil {
		display.Errorf("failed to create NPC agent: %v", err)
		return err
	}

	initialContextSize := npcAgent.GetContextSize()
	display.Colorf(display.ColorBrightRed, "Initial Context Size: %d characters\n", initialContextSize)
	display.Separator()

	npcName := fmt.Sprintf("%s %s", npc.FirstName, npc.FamilyName)

	for {

		input := prompt.NewWithColor("🤖 Ask me something? [" + npcName + "]")
		question, err := input.Run()

		if err != nil {
			display.Errorf("failed to get input: %v", err)
			return err
		}
		if strings.HasPrefix(question, "/bye") {
			display.Infof("👋 Goodbye!")
			break
		}
		if strings.HasPrefix(question, "/messages") {
			display.Infof("💬 Current conversation messages:")
			for i, msg := range npcAgent.GetMessages() {
				display.Infof("Message %d - Role: %s, Content: \n%s", i, msg.Role, msg.Content)
				display.Separator()
			}
			continue
		}

		// Generate response with streamin
		result, err := npcAgent.GenerateStreamCompletion(
			[]messages.Message{
				{
					Role:    roles.User,
					Content: question,
				},
			},
			func(chunk string, finishReason string) error {
				fmt.Print(chunk)
				if finishReason == "stop" {
					display.NewLine()
				}
				return nil
			},
		)
		if err != nil {
			display.Errorf("failed to get completion: %v", err)
			return err
		}
		display.NewLine()
		display.Separator()
		display.KeyValue("Finish reason", result.FinishReason)
		display.KeyValue("Context size", fmt.Sprintf("%d characters", npcAgent.GetContextSize()))
		display.Separator()

	}

	display.Infof("👋 Conversation ended.")
	return err
}
