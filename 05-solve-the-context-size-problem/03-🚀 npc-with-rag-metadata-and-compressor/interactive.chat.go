package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/snipwise/nova/nova-sdk/agents"
	"github.com/snipwise/nova/nova-sdk/agents/chat"
	"github.com/snipwise/nova/nova-sdk/agents/compressor"
	"github.com/snipwise/nova/nova-sdk/agents/rag"
	"github.com/snipwise/nova/nova-sdk/messages"
	"github.com/snipwise/nova/nova-sdk/messages/roles"
	"github.com/snipwise/nova/nova-sdk/models"
	"github.com/snipwise/nova/nova-sdk/toolbox/files"
	"github.com/snipwise/nova/nova-sdk/ui/display"
	"github.com/snipwise/nova/nova-sdk/ui/prompt"
	"github.com/snipwise/nova/nova-sdk/ui/spinner"
)

// startInteractiveChat starts an interactive roleplay session with the NPC
func startInteractiveChat(ctx context.Context, engineURL, npcAgentModelID, sheetJsonFilePath string, ragAgent *rag.Agent, compressorAgent *compressor.Agent) error {

	npc, err := loadNPCSheetFromJsonFile(sheetJsonFilePath)
	if err != nil {
		display.Errorf("failed to load NPC character sheet: %v", err)
		return err
	}

	// === INTERACTIVE NPC CHAT ===
	display.Infof("💬 Starting interactive chat with NPC...")
	display.Separator()

	// Create system instructions for the roleplay agent
	// dnd.chat.system.instructions.md
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
		npc.FirstName, // Always speak in first person as %s

	)
	// NOTE: this time we DO NOT inject the character sheet content directly
	// into the system instructions, as we will use RAG to provide context

	// === CREATE THE ROLEPLAY CHAT AGENT ===
	npcAgent, err := chat.NewAgent(
		ctx,
		agents.Config{
			Name:                    "npc-roleplay-agent",
			EngineURL:               engineURL,
			SystemInstructions:      roleplaySystemInstructions,
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
		// === CONTEXT PACKING WITH COMPRESSOR AGENT ===
		if npcAgent.GetContextSize() > 5000 { // NOTE: for testing, set a low limit

			sp := spinner.NewWithColor("").SetSuffix("compressing...").
				SetFrames(spinner.FramesDots).
				SetDelay(80*time.Millisecond).
				SetColors(
					spinner.ColorBrightYellow, // Prefix color
					spinner.ColorBrightCyan,   // Frame color
					spinner.ColorBrightGreen,  // Suffix color
				)

			display.Infof("🗜️ Context size (%d) exceeded limit, compressing conversation history...", npcAgent.GetContextSize())

			sp.Start()

			compressionResult, err := compressorAgent.CompressContext(npcAgent.GetMessages())
			if err != nil {
				sp.Error("❌ Error during compression")
				display.Errorf("failed to compress conversation history: %v", err)
			} else {
				// Reset conversation with compressed history
				npcAgent.ResetMessages()

				npcAgent.AddMessage(
					roles.System,
					compressionResult.CompressedText,
				)
				sp.Success("✅ Compression successful")
				display.Infof("✅ Conversation history compressed. New context size: %d", npcAgent.GetContextSize())
			}
			display.Separator()
		}

		input := prompt.NewWithColor("🤖 Ask me something? [" + npcName + "]").SetCursorStyle(prompt.CursorBlockBlink)
		question, err := input.RunWithEdit()

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

		// === SIMILARITY SEARCH IN RAG STORE ===
		//similarRecords, err := ragAgent.SearchTopN(question, 0.45, 3)
		similarRecords, err := ragAgent.SearchTopN(question, 0.4, 7) 
		// Tell me about your background
		// Tell me about your background and your family

		similarityContext := ""
		if err != nil {
			display.Errorf("failed to search RAG agent: %v", err)
		} else {
			if len(similarRecords) > 0 {
				display.Infof("📚 Retrieved %d relevant context pieces from RAG store.", len(similarRecords))
				for i, record := range similarRecords {
					display.Separator()
					display.Infof("📄 Context Piece %d (Score: %.4f):\n%s", i+1, record.Similarity, record.Prompt)
					similarityContext += "\n" + record.Prompt
				}
				display.Separator()
			} else {
				display.Infof("📚 No relevant context found in RAG store.")
			}
		}

		var messagesWithContext []messages.Message
		if similarityContext != "" {
			messagesWithContext = []messages.Message{
				{
					Role: roles.System,
					//Content: fmt.Sprintf("Here is some context that might help you answer better:%s\n", similarityContext),
					Content: fmt.Sprintf("# YOUR CHARACTER SHEET%s\n", similarityContext),
				},
				{
					Role:    roles.User,
					Content: fmt.Sprintf("Now, please answer the following question:\n%s", question),
				},
			}
		} else {
			messagesWithContext = []messages.Message{
				{
					Role:    roles.User,
					Content: question,
				},
			}
		}

		// Generate response with streamin
		result, err := npcAgent.GenerateStreamCompletion(
			messagesWithContext,
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
