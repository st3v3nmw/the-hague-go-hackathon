package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mattermost/mattermost/server/public/model"
	"os"
	"strings"
)

// Provides a list of post messages from the same thread as `postId`, in reverse order.
func GetThread(client *model.Client4, postId string) ([]*model.Post, error) {
	postList, _, err := client.GetPostThread(context.Background(), postId, "", false)

	if err != nil {
		return nil, err
	}

	return postList.ToSlice(), nil
}

func listenToEvents(wsClient *model.WebSocketClient) {
	for {
		wsClient.Listen()

		for event := range wsClient.EventChannel {
			go handleEvent(wsClient, event)
		}
	}
}

func handleEvent(wsClient *model.WebSocketClient, event *model.WebSocketEvent) {
	if event.EventType() != model.WebsocketEventPosted {
		return
	}

	post := &model.Post{}
	err := json.Unmarshal([]byte(event.GetData()["post"].(string)), &post)

	if err != nil {
		return
	}

	if post.UserId == os.Getenv("MM_USERID") {
		return
	}

	message := post.Message

	if strings.TrimSpace(message) != "!mm-bot tldr" {
		return
	}

	// TODO: get thread for post and summarize it
	fmt.Println(post)
}
