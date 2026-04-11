package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)


func main() {
	ctx := context.Background()

	engineBaseURL := os.Getenv("ENGINE_BASE_URL")
	if engineBaseURL == "" {
		engineBaseURL = "http://localhost:12434/engines/v1/"

	}
	modelId := "huggingface.co/menlo/jan-nano-gguf:q4_k_m"
	//modelId := "ai/qwen2.5:0.5B-F16"

	client := openai.NewClient(
		option.WithBaseURL(engineBaseURL),
		option.WithAPIKey(""),
	)

	// Get a list of countries somewhere in the world
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"countries": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"name": map[string]any{
							"type": "string",
						},
						"capital": map[string]any{
							"type": "string",
						},
						"languages": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "string",
							},
						},
					},
					"required": []string{"name", "capital", "languages"},
				},
			},
		},
		"required": []string{"countries"},
	}

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "List of countries",
		Description: openai.String("List of countries in the world"),
		Schema:      schema,
		Strict:      openai.Bool(true),
	}

	userQuestion := openai.UserMessage("List of 10 countries in Europe")

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			userQuestion,
		},
		Model:       modelId,
		Temperature: openai.Opt(0.0),
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
	}

	// Make completion request
	completion, err := client.Chat.Completions.New(ctx, params)

	if err != nil {
		panic(err)
	}

	data := completion.Choices[0].Message.Content

	var countriesList map[string][]any

	err = json.Unmarshal([]byte(data), &countriesList)

	if err != nil {
		panic(err)
	}
	fmt.Println("Countries List:")
	for idx, country := range countriesList["countries"] {
		fmt.Println(idx, ".", country)
	}

}
