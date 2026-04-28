package main

import (
	"context"
	"fmt"
	"log"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func main() {
		// Ollama Chat base URL
		baseURL := "http://localhost:11434/v1/"
		model := "qwen2:0.5b"
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(""),
	)

	ctx := context.Background()

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are an expert of medieval role playing games."),
		openai.UserMessage("[Brief] What is a dungeon crawler game?"),
	}

	param := openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       model,
		Temperature: openai.Opt(0.0),
	}

	stream := client.Chat.Completions.NewStreaming(ctx, param)

	for stream.Next() {
		chunk := stream.Current()
		// Stream each chunk as it arrives
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}

	if err := stream.Err(); err != nil {
		log.Fatalln("😡:", err)
	}
}
