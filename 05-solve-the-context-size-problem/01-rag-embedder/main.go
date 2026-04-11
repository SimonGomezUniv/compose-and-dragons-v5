package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/option"
)

func main() {
	ctx := context.Background()

	engineBaseURL := os.Getenv("ENGINE_BASE_URL")
	if engineBaseURL == "" {
		engineBaseURL = "http://localhost:12434/engines/v1/"

	}
	modelId := "ai/mxbai-embed-large:latest"

	oaiPlugin := &openai.OpenAI{
		APIKey: "I💙DockerModelRunner",
		Opts: []option.RequestOption{
			option.WithBaseURL(engineBaseURL),
		},
	}
	genKitInstance := genkit.Init(ctx, genkit.WithPlugins(oaiPlugin))

	//embedder := compatPlugin.DefineEmbedder("dmr", "openai/ai/nomic-embed-text-v1.5:latest", nil)

	/*
		# Embedders
		An embedder is a function that takes content (text, images, audio, etc.)
		and creates a numeric vector that encodes the semantic meaning of the original content.
		As mentioned above, embedders are leveraged as part of the process of indexing.
		However, they can also be used independently to create embeddings without an index.
	*/

	// NOTE: Embedder definition/creation
	// you don't need to specify the provider again, it's already set in the plugin 🤔
	// == you don't need to prefix the model name with the provider
	embedder := oaiPlugin.DefineEmbedder(modelId, nil)

	chunks := []string{
		"Squirrels run in the forest",
		"Birds fly in the sky",
		"Frogs swim in the pond",
		"Fishes swim in the sea",
	}

	resp, err := genkit.Embed(ctx, genKitInstance,
		ai.WithEmbedder(embedder),
		ai.WithTextDocs(chunks...),
	)

	if err != nil {
		log.Fatal(err)
	}

	for i, emb := range resp.Embeddings {
		fmt.Printf("Chunk %d (%s) \nembedding: %v\n", i, chunks[i], emb)
		fmt.Println(strings.Repeat("=", 80))
	}

	respUserQuestion, err := genkit.Embed(ctx, genKitInstance,
		ai.WithEmbedder(embedder),
		ai.WithTextDocs("Which animals swim?"),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Question (Which animals swim?) \nembedding: %v\n", respUserQuestion.Embeddings[0])
	// you can now use respUserQuestion.Embeddings[0] to find the most similar chunks in resp.Embeddings
	// and return them as context to a LLM prompt



}
