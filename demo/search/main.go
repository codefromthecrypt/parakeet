package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	//ollamaUrl := "http://bob.local:11434"
	
	var embeddingsModel = "magicoder:latest"
	var smallChatModel = "magicoder:latest"

	store := embeddings.BboltVectorStore{}
	store.Initialize("../embeddings.db")

	systemContent := `You are a Golang developer and an expert in computer programming.
	Please make friendly answer for the noobs. Use the provided context to answer.
	Add source code examples if you can.`

	// Question for the Chat system
	//userContent := `[Brief] How to create a stream completion with Parakeet?`
	userContent := `[Brief] How to create a stream chat completion with Parakeet?`


	// Create an embedding from the user question
	embeddingFromQuestion, err := embeddings.CreateEmbedding(
		ollamaUrl,
		llm.Query4Embedding{
			Model:  embeddingsModel,
			Prompt: userContent,
		},
		"question",
	)
	if err != nil {
		log.Fatalln("😡:", err)
	}
	fmt.Println("🔎 searching for similarity...")

	//similarity, _ := store.SearchMaxSimilarity(embeddingFromQuestion)

	similarities, _ := store.SearchSimilarities(embeddingFromQuestion, 0.0)

	documentsContent := "<context>"
	for _, similarity := range similarities {
		documentsContent += fmt.Sprintf("<doc>%s</doc>", similarity.Prompt)
	}
	documentsContent += "</context>"

	fmt.Println("🎉 similarities", similarities)

	query := llm.Query{
		Model: smallChatModel,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "system", Content: documentsContent},
			{Role: "user", Content: userContent},
		},
		Options: llm.Options{
			Temperature: 0.4,
			RepeatLastN: 2,
		},
		Stream: false,
	}

	fmt.Println("")
	fmt.Println("🤖 answer:")

	// Answer the question
	_, err = completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}

}
