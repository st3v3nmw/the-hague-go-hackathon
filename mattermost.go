package main

import (
	"context"
	"encoding/json"
	"log"
	"fmt"
	"github.com/mattermost/mattermost/server/public/model"
	"os"
	"strings"
)

// Provides a list of post messages from the same thread as `postId`, in reverse order.
func GetThread(bot *Bot, postId string) ([]*model.Post, error) {
	postList, _, err := bot.apiClient.GetPostThread(context.Background(), postId, "", false)

	if err != nil {
		return nil, err
	}

	return postList.ToSlice(), nil
}

func listenToEvents(bot *Bot) {
	for {
		bot.wsClient.Listen()

		for event := range bot.wsClient.EventChannel {
			go handleEvent(bot, event)
		}
	}
}

func handleEvent(bot *Bot, event *model.WebSocketEvent) {
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

	postList, err := GetThread(bot, post.Id)

	if err != nil {
		log.Println("Error getting thread")
		log.Println(err)
		return
	}

	userIds := make([]string, len(postList))

	for _, post := range postList {
		userIds = append(userIds, post.UserId)
	}

	users, err := getUsers(bot, userIds)

	if err != nil {
		log.Println("Error getting usernames")
		log.Println(err)
		return
	}

	log.Println("Sending to LLM")
	log.Println(postList)
	log.Println(users)
	// TODO: something with users and posts
	summary, err := Summarize(bot.llmClient, postList, users)

	if err != nil {
		log.Println("Error summarizing posts")
		log.Println(err)
		return
	}

	fmt.Println(summary)
}

func getUsers(bot *Bot, ids []string) (map[string]*model.User, error) {
	userList, _, err := bot.apiClient.GetUsersByIds(context.Background(), ids)

	if err != nil {
		return nil, err
	}

	userIds := make(map[string]*model.User)
	for _, user := range userList {
		userIds[user.Id] = user
	}

	return userIds, nil
}
