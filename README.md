# 🦜🪺 Parakeet
<!-- meta-data 
topic: this an introduction on what is parakeet
key-words: parakeet genai apps
-->

Parakeet is the simplest Go library to create **GenAI apps** with **[Ollma](https://ollama.com/)**.

> A GenAI app is an application that uses generative AI technology. Generative AI can create new text, images, or other content based on what it's been trained on. So a GenAI app could help you write a poem, design a logo, or even compose a song! These are still under development, but they have the potential to be creative tools for many purposes. - [Gemini](https://gemini.google.com)

> ✋ Parakeet is only for creating GenAI apps generating **text** (not image, music,...).

<!--split-->

<!-- meta-data 
topic: how to install parakeet
-->
## Install

```bash
go get github.com/parakeet-nest/parakeet
```

<!--split-->

<!-- meta-data 
topic: how to do a simple completion with parakeet
key-words: completion parakeet
-->
## Simple completion

The simple completion can be used to generate a response for a given prompt with a provided model.

```golang
package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"fmt"
	"log"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "tinydolphin"

	options := llm.Options{
		Temperature: 0.5,  // default (0.8)
	}

	question := llm.Query{
		Model: model,
		Prompt: "Who is James T Kirk?",
		Options: options,
	}

	answer, err := completion.Generate(ollamaUrl, question)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Response)
}
```

<!--split-->

<!-- meta-data 

-->
### Simple completion with stream

```golang
package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"fmt"
	"log"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "tinydolphin"

	options := llm.Options{
		Temperature: 0.5, // default (0.8)
	}

	question := llm.Query{
		Model: model,
		Prompt: "Who is James T Kirk?",
		Options: options,
	}
	
	answer, err := completion.GenerateStream(ollamaUrl, question,
		func(answer llm.Answer) error {
			fmt.Print(answer.Response)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}
}
```

<!--split-->

<!-- meta-data 

-->
### Completion with context
> see: https://github.com/ollama/ollama/blob/main/docs/api.md#generate-a-completion

> The context can be used to keep a short conversational memory for the next completion.

```golang
package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "tinydolphin"

	options := llm.Options{
		Temperature: 0.5, // default (0.8)
	}

	firstQuestion := llm.Query{
		Model: model,
		Prompt: "Who is James T Kirk?",
		Options: options,
	}

	answer, err := completion.Generate(ollamaUrl, firstQuestion)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Response)

	fmt.Println()

	secondQuestion := llm.Query{
		Model: model,
		Prompt: "Who is his best friend?",
		Context: answer.Context,
		Options: options,
	}

	answer, err = completion.Generate(ollamaUrl, secondQuestion)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Response)
}
```

<!--split-->

<!-- meta-data 

-->
## Chat completion

The chat completion can be used to generate a conversational response for a given set of messages with a provided model.

```golang
package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "deepseek-coder"

	systemContent := `You are an expert in computer programming.
	Please make friendly answer for the noobs.
	Add source code examples if you can.`

	userContent := `I need a clear explanation regarding the following question:
	Can you create a "hello world" program in Golang?
	And, please, be structured with bullet points`

	options := llm.Options{
		Temperature: 0.5, // default (0.8)
		RepeatLastN: 2, // default (64)
		RepeatPenalty: 2.0, // default (1.1)
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
		Stream: false,
	}

	answer, err := completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Message.Content)
}
```

✋ **To keep a conversational memory** for the next chat completion, update the list of messages with the previous question and answer.
> I plan to add the support of [bbolt](https://github.com/etcd-io/bbolt) in the incoming v0.0.1 of Parakeet to store the conversational memory.

<!--split-->

<!-- meta-data 

-->
### Chat completion with stream

```golang
package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "deepseek-coder"

	systemContent := `You are an expert in computer programming.
	Please make friendly answer for the noobs.
	Add source code examples if you can.`

	userContent := `I need a clear explanation regarding the following question:
	Can you create a "hello world" program in Golang?
	And, please, be structured with bullet points`

	options := llm.Options{
		Temperature: 0.5, // default (0.8)
		RepeatLastN: 2, // default (64) 
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
		Stream:  false,
	}

	_, err := completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}
}
```

<!--split-->

<!-- meta-data 

-->
## Chat completion with conversational memory

### In memory history

To store the messages in memory, use `history.MemoryMessages`

```golang
package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/history"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "tinydolphin" // fast, and perfect answer (short, brief)

	conversation := history.MemoryMessages{
		Messages: make(map[string]llm.MessageRecord),
	}

	systemContent := `You are an expert with the Star Trek series. use the history of the conversation to answer the question`

	userContent := `Who is James T Kirk?`

	options := llm.Options{
		Temperature: 0.5,
		RepeatLastN: 2,  
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	// Ask the question
	answer, err := completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		},
	)
	if err != nil {
		log.Fatal("😡:", err)
	}

	// Save the conversation
	_, err = conversation.SaveMessage("1", llm.Message{
		Role:    "user",
		Content: userContent,
	})
	if err != nil {
		log.Fatal("😡:", err)
	}

	_, err = conversation.SaveMessage("2", llm.Message{
		Role:    "system",
		Content: answer.Message.Content,
	})

	if err != nil {
		log.Fatal("😡:", err)
	}

	// New question
	userContent = `Who is his best friend ?`

	previousMessages, _ := conversation.GetAllMessages()

	// (Re)Create the conversation
	conversationMessages := []llm.Message{}
	// instruction
	conversationMessages = append(conversationMessages, llm.Message{Role: "system", Content: systemContent})
	// history
	conversationMessages = append(conversationMessages, previousMessages...)
	// last question
	conversationMessages = append(conversationMessages, llm.Message{Role: "user", Content: userContent})

	query = llm.Query{
		Model:    model,
		Messages: conversationMessages,
		Options:  options,
	}

	answer, err = completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		},
	)
	fmt.Println()
	if err != nil {
		log.Fatal("😡:", err)
	}

}
```

<!--split-->

<!-- meta-data 

-->
### Bbolt history

**[Bbolt](https://github.com/etcd-io/bbolt)** is an embedded key/value database for Go.

To store the messages in a bbolt bucket, use `history.BboltMessages`

```golang
conversation := history.BboltMessages{}
conversation.Initialize("../conversation.db")
```

> 👀 you will find a complete example in `examples/11-chat-conversational-bbolt`
> - `examples/11-chat-conversational-bbolt/begin`: start a conversation and save the history
> - `examples/11-chat-conversational-bbolt/resume`: load the messages from the history bucket and resue the conversation

<!--split-->

<!-- meta-data 

-->
## Embeddings

### Create embeddings

```golang
embedding, err := embeddings.CreateEmbedding(
	ollamaUrl,
	llm.Query4Embedding{
		Model:  "all-minilm",
		Prompt: "Jean-Luc Picard is a fictional character in the Star Trek franchise.",
	},
	"Picard", // identifier
)
```

<!--split-->

<!-- meta-data 

-->
## Vector stores

A vector store allows to store and search for embeddings in an efficient way.

### In memory vector store

**Create a store**:
```golang
store := embeddings.MemoryVectorStore{
	Records: make(map[string]llm.VectorRecord),
}
```

**Save embeddings**:
```golang
store.Save(embedding)
```

**Search embeddings**:
```golang
embeddingFromQuestion, err := embeddings.CreateEmbedding(
	ollamaUrl,
	llm.Query4Embedding{
		Model:  "all-minilm",
		Prompt: "Who is Jean-Luc Picard?",
	},
	"question",
)
// find the nearest vector
similarity, _ := store.SearchMaxSimilarity(embeddingFromQuestion)

documentsContent := `<context><doc>` + similarity.Prompt + `</doc></context>`
```

> 👀 you will find a complete example in `examples/08-embeddings`

<!--split-->

<!-- meta-data 

-->
### Bbolt vector store

**[Bbolt](https://github.com/etcd-io/bbolt)** is an embedded key/value database for Go.

**Create a store, and open an existing store**:
```golang
store := embeddings.BboltVectorStore{}
store.Initialize("../embeddings.db")
```

> 👀 you will find a complete example in `examples/09-embeddings-bbolt`
> - `examples/09-embeddings-bbolt/create-embeddings`: create and populate the vector store
> - `examples/09-embeddings-bbolt/use-embeddings`: search similarities in the vector store


<!--split-->

<!-- meta-data 

-->
## Create embeddings from text files and Similarity search

### Create embeddings
```golang
ollamaUrl := "http://localhost:11434"
embeddingsModel := "all-minilm"

store := embeddings.BboltVectorStore{}
store.Initialize("../embeddings.db")

// Parse all golang source code of the examples
// Create embeddings from documents and save them in the store
counter := 0
_, err := content.ForEachFile("../../examples", ".go", func(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fmt.Println("📝 Creating embedding from:", path)
	counter++
	embedding, err := embeddings.CreateEmbedding(
		ollamaUrl,
		llm.Query4Embedding{
			Model:  embeddingsModel,
			Prompt: string(data),
		},
		strconv.Itoa(counter), // don't forget the id (unique identifier)
	)
	fmt.Println("📦 Created: ", len(embedding.Embedding))

	if err != nil {
		fmt.Println("😡:", err)
	} else {
		_, err := store.Save(embedding)
		if err != nil {
			fmt.Println("😡:", err)
		}
	}
	return nil
})
if err != nil {
	log.Fatalln("😡:", err)
}
```

<!--split-->

<!-- meta-data 

-->
### Similarity search

```golang
ollamaUrl := "http://localhost:11434"
embeddingsModel := "all-minilm"
chatModel := "magicoder:latest"

store := embeddings.BboltVectorStore{}
store.Initialize("../embeddings.db")

systemContent := `You are a Golang developer and an expert in computer programming.
Please make friendly answer for the noobs. Use the provided context and doc to answer.
Add source code examples if you can.`

// Question for the Chat system
userContent := `How to create a stream chat completion with Parakeet?`

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

similarities, _ := store.SearchSimilarities(embeddingFromQuestion, 0.3)

// Generate the context from the similarities
// This will generate a string with a content like this one:
// `<context><doc>...<doc><doc>...<doc></context>`
documentsContent := embeddings.GenerateContextFromSimilarities(similarities)

fmt.Println("🎉 similarities", len(similarities))

query := llm.Query{
	Model: chatModel,
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
```

<!--split-->

<!-- meta-data 

-->
## Function Calling

What is **"Function Calling"**? First, it's not a feature where a LLM can call and execute a function. "Function Calling" is the ability for certain LLMs to provide a specific output with the same format (we could say: "a predictable output format").

So, the principle is simple:

- You (or your GenAI application) will create a prompt with a delimited list of tools (the functions) composed by name, descriptions, and parameters: `SayHello`, `AddNumbers`, etc.
- Then, you will add your question ("Hey, say 'hello' to Bob!") to the prompt and send all of this to the LLM.
- If the LLM "understand" that the `SayHello` function can be used to say "hello" to Bob, then the LLM will answer with only the name of the function with the parameter(s). For example: `{"name":"SayHello","arguments":{"name":"Bob"}}`.

Then, it will be up to you to implement the call of the function.

The [latest version (v0.3) of Mistral 7b](https://ollama.com/library/mistral:7b) supports function calling and is available for Ollama.

### Define a list of tools

First, you have to provide the LLM with a list of tools with the following format:

```golang
toolsList := []llm.Tool{
	{
		Type: "function",
		Function: llm.Function{
			Name:        "hello",
			Description: "Say hello to a given person with his name",
			Parameters: llm.Parameters{
				Type: "object",
				Properties: map[string]llm.Property{
					"name": {
						Type:        "string",
						Description: "The name of the person",
					},
				},
				Required: []string{"name"},
			},
		},
	},
	{
		Type: "function",
		Function: llm.Function{
			Name:        "addNumbers",
			Description: "Make an addition of the two given numbers",
			Parameters: llm.Parameters{
				Type: "object",
				Properties: map[string]llm.Property{
					"a": {
						Type:        "number",
						Description: "first operand",
					},
					"b": {
						Type:        "number",
						Description: "second operand",
					},
				},
				Required: []string{"a", "b"},
			},
		},
	},
}
```

### Generate a prompt from the tools list and the user instructions

The `tools.GenerateContent` method generates a string with the tools in JSON format surrounded by `[AVAILABLE_TOOLS]` and `[/AVAILABLE_TOOLS]`:
```golang
toolsContent, err := tools.GenerateContent(toolsList)
if err != nil {
	log.Fatal("😡:", err)
}
```


The `tools.GenerateInstructions` method generates a string with the user instructions surrounded by `[INST]` and `[/INST]`:
```golang
userContent := tools.GenerateInstructions(`say "hello" to Bob`)
```

Then, you can add these two strings to the messages list:
```golang
messages := []llm.Message{
	{Role: "system", Content: toolsContent},
	{Role: "user", Content: userContent},
}
```

### Send the prompt (messages) to the LLM

It's important to set the `Temperature` to `0.0`:
```golang
options := llm.Options{
	Temperature:   0.0,
	RepeatLastN:   2,
	RepeatPenalty: 2.0,
}

You must set the `Format` to `json` and `Raw` to `true`:
query := llm.Query{
	Model: model,
	Messages: messages,
	Options: options,
	Format:  "json",
	Raw:     true,
}
```
> When building the payload to be sent to Ollama, we need to set the `Raw` field to true, thanks to that, no formatting will be applied to the prompt (we override the prompt template of Mistral), and we need to set the `Format` field to `"json"`.

No you can call the `Chat` method. The answer of the LLM will be in JSON format:
```golang
answer, err := completion.Chat(ollamaUrl, query)
if err != nil {
	log.Fatal("😡:", err)
}
// PrettyString is a helper that prettyfies the JSON string
result, err := gear.PrettyString(answer.Message.Content)
if err != nil {
	log.Fatal("😡:", err)
}
fmt.Println(result)
```

You should get this answer:
```json
{
  "name": "hello",
  "arguments": {
    "name": "Bob"
  }
}
```

You can try with the other tool (or function):
```golang
userContent := tools.GenerateInstructions(`add 2 and 40`)
```

You should get this answer:
```json
{
  "name": "addNumbers",
  "arguments": {
    "a": 2,
    "b": 40
  }
}
```

> **Remark**: always test the format of the output, even if Mistral is trained for "function calling", the result are not entirely predictable.

Look at this sample for a complete sample: [examples/15-mistral-function-calling](examples/15-mistral-function-calling)


## Function Calling with LLMs that do not implement Function Calling

It is possible to reproduce this feature with some LLMs that do not implement the "Function Calling" feature natively, but we need to supervise them and explain precisely what we need. The result (the output) will be less predictable, so you will need to add some tests before using the output, but with some "clever" LLMs, you will obtain correct results. I did my experiments with **[phi3:mini](https://ollama.com/library/phi3:mini)**.

The trick is simple:

Add this message at the begining of the list of messages:
```golang
systemContentIntroduction := `You have access to the following tools:`
```

Add this message at the end of the list of messages, just before the user message:
```golang
systemContentInstructions := `If the question of the user matched the description of a tool, the tool will be called.
To call a tool, respond with a JSON object with the following structure: 
{
	"name": <name of the called tool>,
	"arguments": {
	<name of the argument>: <value of the argument>
	}
}

search the name of the tool in the list of tools with the Name field
`
```

At the end, you will have this:
```golang
messages := []llm.Message{
	{Role: "system", Content: systemContentIntroduction},
	{Role: "system", Content: toolsContent},
	{Role: "system", Content: systemContentInstructions},
	{Role: "user", Content: `say "hello" to Bob`},
}
```

Look at this sample for a complete sample: [examples/17-fake-function-calling](examples/17-fake-function-calling)

<!--split-->

<!-- meta-data 

-->
## Wasm plugins

The release `0.0.6` of Parakeet brings the support of **WebAssembly** thanks to the **[Extism project](https://extism.org/)**. That means you can write your own wasm plugins for Parakeet to add new features (for example, a chunking helper for doing RAG) with various languages (Rust, Go, C, ...).

Or you can use the Wasm plugins with the "Function Calling" feature, which is implemented in Parakeet.

You can find an example of "Wasm Function Calling" in [examples/18-call-functions-for-real](examples/18-call-functions-for-real) - the wasm plugin is located in the `wasm` folder and it is built with **[TinyGo](https://tinygo.org/)**.

🚧 more samples to come.
<!--split-->

<!-- meta-data 

-->
## Parakeet Demos

- https://github.com/parakeet-nest/parakeet-demo
- https://github.com/parakeet-nest/tiny-genai-stack

## Blog Posts

- [Parakeet, an easy way to create GenAI applications with Ollama and Golang](https://k33g.hashnode.dev/parakeet-an-easy-way-to-create-genai-applications-with-ollama-and-golang)
- [Understand RAG with Parakeet](https://k33g.hashnode.dev/understand-rag-with-parakeet)
-[Function Calling with Ollama, Mistral 7B, Bash and Jq](https://k33g.hashnode.dev/function-calling-with-ollama-mistral-7b-bash-and-jq)
- [Function Calling with Ollama and LLMs that do not support function calling](https://k33g.hashnode.dev/function-calling-with-ollama-and-llms-that-do-not-support-function-calling)

