package main

import (
	"context"
	"fmt"

	"github.com/snipwise/nova/nova-sdk/agents"
	"github.com/snipwise/nova/nova-sdk/agents/chat"
	"github.com/snipwise/nova/nova-sdk/messages"
	"github.com/snipwise/nova/nova-sdk/messages/roles"
	"github.com/snipwise/nova/nova-sdk/models"
	"github.com/snipwise/nova/nova-sdk/ui/display"
)

func main() {
	ctx := context.Background()

	// Create a simple agent without exposing OpenAI SDK types
	agent, err := chat.NewAgent(
		ctx,
		agents.Config{
			EngineURL:          "http://localhost:11434/v1/",
			SystemInstructions: "You are an expert of medieval role playing games.",
			KeepConversationHistory: true,
		},
		models.Config{
			Name:        "qwen2:0.5b",
			Temperature: models.Float64(0.8),
		},
	)
	if err != nil {
		panic(err)
	}

	result, err := agent.GenerateStreamCompletion(
		[]messages.Message{
			{Role: roles.User, Content: "[Brief] What is a dungeon crawler game?"},
		},
		func(chunk string, finishReason string) error {
			fmt.Print(chunk)
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	display.NewLine()
	display.Separator()
	display.KeyValue("Finish reason", result.FinishReason)
	display.KeyValue("Context size", fmt.Sprintf("%d characters", agent.GetContextSize()))
	display.Separator()

}
