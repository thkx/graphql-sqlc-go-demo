package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/documentloaders"
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

	// OllamaChat()
}
