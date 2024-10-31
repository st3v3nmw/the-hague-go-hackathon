package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func newLLMClient() (*ollama.LLM, error) {
	url := os.Getenv("OLLAMA_URL")
	llm, err := ollama.New(ollama.WithServerURL(url))
	if err != nil {
		return nil, err
	}

	return llm, nil
}

func PromptLLM(llm *ollama.LLM, prompt string) (string, error) {
	ctx := context.Background()
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		return "", err
	}

	return completion, nil
}

func Summarize(llm *ollama.LLM, posts []model.Post, users []model.User) (string, error) {
	prompt := "Summarize the following thread of messages: \n"

	for i := range len(posts) {
		prompt += fmt.Sprintf("%s: %s\n", users[i].Username, posts[i].Message)
	}

	return PromptLLM(llm, prompt)
}
