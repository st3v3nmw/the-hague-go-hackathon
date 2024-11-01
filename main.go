package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/tmc/langchaingo/llms/ollama"
)

type Bot struct {
	apiClient *model.Client4
	wsClient  *model.WebSocketClient
	llmClient *ollama.LLM
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	authToken := os.Getenv("MM_AUTHTOKEN")
	mmURL := os.Getenv("MM_URL")
	apiClient := model.NewAPIv4Client(fmt.Sprintf("https://%s", mmURL))
	apiClient.SetToken(authToken)

	wsClient, err := model.NewWebSocketClient4(fmt.Sprintf("wss://%s", mmURL), authToken)

	if err != nil {
		log.Fatal("Error connecting websocket client")
	}

	llmClient, err := newLLMClient()

	if err != nil {
		log.Fatal("Error creating LLM client")
	}

	bot := Bot{
		apiClient,
		wsClient,
		llmClient,
	}

	listenToEvents(&bot)
}
