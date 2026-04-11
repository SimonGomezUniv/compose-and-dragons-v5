package main

import (
	"context"
	"demo-rag-retriever/rag"
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

	// -------------------------------------------------
	// Create a "in memory" vector store
	// -------------------------------------------------
	store := rag.MemoryVectorStore{
		Records: make(map[string]rag.VectorRecord),
	}

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
		// Store the embedding in the vector store
		record, errSave := store.Save(rag.VectorRecord{
			Prompt:    chunks[i],
			Embedding: emb.Embedding,
		})
		if errSave != nil {
			log.Fatal(errSave)
		}
		fmt.Println("Saved record:", record.Prompt, record.Id)
	}

	respUserQuestion, err := genkit.Embed(ctx, genKitInstance,
		ai.WithEmbedder(embedder),
		ai.WithTextDocs("Which animals swim?"),
	)
	_ = respUserQuestion
	if err != nil {
		log.Fatal(err)
	}

	// -------------------------------------------------
	// Use the custom retriever to find similar documents
	// -------------------------------------------------

	// Create the memory vector retriever
	memoryRetriever := rag.DefineMemoryVectorRetriever(genKitInstance, &store, embedder)

	// Create a query document from the user question
	queryDoc := ai.DocumentFromText("Which animals swim?", nil)

	// Create a retriever request with custom options
	request := &ai.RetrieverRequest{
		Query: queryDoc,
		Options: rag.MemoryVectorRetrieverOptions{
			Limit:      0.5, // Lower similarity threshold to get more results
			MaxResults: 3,   // Return top 3 results
		},
	}

	// Use the memory vector retriever to find similar documents
	retrieveResponse, err := memoryRetriever.Retrieve(ctx, request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("\nFound %d similar documents:\n", len(retrieveResponse.Documents))
	for i, doc := range retrieveResponse.Documents {
		similarity := doc.Metadata["cosine_similarity"]
		id := doc.Metadata["id"]
		fmt.Printf("%d. ID: %s, Similarity: %.4f\n", i+1, id, similarity)
		fmt.Printf("   Content: %s\n\n", doc.Content[0].Text)
	}

}
