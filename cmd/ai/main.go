package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/memory"
)

func main() {
	ctx := context.Background()

	file, err := os.Open("./test.pdf")
	if err != nil {
		log.Panicln(err)
	}

	docs := documentloaders.NewPDF(file, 100)
	documents, err := docs.Load(ctx)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(documents)
	// AgentNode()

	// stringPromptTemplates()

	// OllamaPrompt()

	llm, err := ollama.New(
		ollama.WithModel("gemma3:270m"),
		ollama.WithServerURL("http://localhost:11434"),
	)
	if err != nil {
		log.Fatalln(err)
	}

	bufferMemory := memory.NewConversationBuffer()

	conversionChain := chains.NewConversation(llm, bufferMemory)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n 我：")

		userInput, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v\n", err)
		}

		userInput = strings.TrimSpace(userInput)

		if userInput == "" {
			fmt.Println("请输入问题")
			continue
		}

		fmt.Print("AI: ")
		fmt.Print("Thinking...")

		response, err := chains.Run(ctx, conversionChain, userInput)
		fmt.Print("\r AI: ")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Println(response)

	}
}
