package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/shared"
)

// MODEL_RUNNER_BASE_URL=http://localhost:11434 go run main.go
func main() {
	ctx := context.Background()

	engineBaseURL := os.Getenv("ENGINE_BASE_URL")
	if engineBaseURL == "" {
		engineBaseURL = "http://localhost:11434/v1/"

	}
	//modelId := "huggingface.co/menlo/jan-nano-gguf:q4_k_m"
	modelId := "qwen2:0.5b"

	client := openai.NewClient(
		option.WithBaseURL(engineBaseURL),
		option.WithAPIKey(""),
	)

	// TOOLS: Define the tools available to the model
	// Each tool must have a name, description, and parameters schema

	// TOOL:
	sayHelloTool := openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
		Name:        "say_hello",
		Description: openai.String("Say hello to the given person name"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]string{
					"type": "string",
				},
			},
			"required": []string{"name"},
		},
	})

	// TOOLS: used by the parameters request
	vulcanSaluteTool := openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
		Name:        "vulcan_salute",
		Description: openai.String("Give a vulcan salute to the given person name"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]string{
					"type": "string",
				},
			},
			"required": []string{"name"},
		},
	})

	tools := []openai.ChatCompletionToolUnionParam{
		sayHelloTool,
		vulcanSaluteTool,
	}

	// USER MESSAGE:
	userQuestion := openai.UserMessage(`
		Say hello to Jean-Luc Picard 
		and Say hello to James Kirk 
		and make a Vulcan salute to Spock
	`)

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			userQuestion,
		},
		ParallelToolCalls: openai.Bool(true),
		Tools:             tools,
		Model:             modelId,
		Temperature:       openai.Opt(0.0),
	}

	// Make [COMPLETION] request
	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	// TOOL CALLS: Extract tool calls from the response
	toolCalls := completion.Choices[0].Message.ToolCalls

	// Return early if there are no tool calls
	if len(toolCalls) == 0 {
		fmt.Println("😡 No function call")
		fmt.Println()
		return
	}

	// Display the tool calls
	for _, toolCall := range toolCalls {
		var args map[string]any

		switch toolCall.Function.Name {
		case "say_hello":
			args, _ = JsonStringToMap(toolCall.Function.Arguments)
			fmt.Println(sayHello(args))

		case "vulcan_salute":
			args, _ = JsonStringToMap(toolCall.Function.Arguments)
			fmt.Println(vulcanSalute(args))

		default:
			fmt.Println("Unknown function call:", toolCall.Function.Name)
		}
	}

}

func JsonStringToMap(jsonString string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func sayHello(arguments map[string]interface{}) string {

	if name, ok := arguments["name"].(string); ok {
		return "Hello " + name
	} else {
		return ""
	}
}

func vulcanSalute(arguments map[string]interface{}) string {
	if name, ok := arguments["name"].(string); ok {
		return "Live long and prosper " + name
	} else {
		return ""
	}
}
