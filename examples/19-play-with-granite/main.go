package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/prompt"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	model := "granite-code:3b"

	info, status, err := llm.ShowModelInformation(ollamaUrl, model)
	if err != nil {
		if status == 404 {
			fmt.Println("✋ we need to pull the model")
			result, status, errPull := llm.PullModel(ollamaUrl, model)
			if errPull != nil {
				log.Fatal(errPull, status)
			}
			fmt.Println(result)
		}
		log.Fatal(err)
	}

	fmt.Println("📝 Model information:", info.Details.Family)

	systemContent := `You are an expert in computer programming.`

	userContent := prompt.Brief(`can you generate an "hello world" program in Golang`)

	options := llm.Options{
		Temperature: 0.0, // default (0.8)
		RepeatLastN: 2,   // default (64) the default value will "freeze" deepseek-coder
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	_, err = completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}
}
