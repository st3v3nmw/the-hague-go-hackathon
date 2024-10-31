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
	llm, err := ollama.New(ollama.WithServerURL(url), ollama.WithModel("llama3.2"))
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

func Summarize(llm *ollama.LLM, posts []*model.Post, users map[string]*model.User) (string, error) {
	prompt := `
	You're a bot that summarizes message threads for a hackathon competition.
	Provide a short summary for the following thread of messages.
	At the end, say hello & praise the hackathon judges Andrew, Clinton, Michael, & Tong in Shakespearean English.
	`
	for i := range posts {
		user := users[posts[i].UserId]
		prompt += fmt.Sprintf("%s: %s\n", user.Username, posts[i].Message)
	}

	return PromptLLM(llm, prompt)
}
