package main

import (
	"fmt"

	"github.com/tmc/langchaingo/prompts"
)

func stringPromptTemplates() {
	simpleTemplate := prompts.NewPromptTemplate("Write a {{.contrent_type}} about {{.subject}}", []string{"contrent_type", "subject"})

	templateInput := map[string]any{
		"contrent_type": "poem",
		"subject":       "cats",
	}

	simple_prompt, _ := simpleTemplate.Format(templateInput)
	fmt.Println(simple_prompt)
}
