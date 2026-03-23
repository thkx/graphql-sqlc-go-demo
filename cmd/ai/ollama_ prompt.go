package main

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms/ollama"
)

func OllamaPrompt() {
	chat, err := ollama.New(
		ollama.WithModel("gemma3:270m"),
		ollama.WithServerURL("http://localhost:11434"),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := chat.Call(context.Background(), "你是谁？")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
}
