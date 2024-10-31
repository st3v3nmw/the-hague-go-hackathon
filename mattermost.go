package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
)

const TDLR_COMMAND = "!mm-bot tldr"

func IsTDLRCommand(message string) bool {
	return strings.TrimSpace(message) == TDLR_COMMAND
}

// Provides a list of post messages from the same thread as `postId`, in reverse order.
func GetThread(bot *Bot, postId string) ([]*model.Post, error) {
	postList, _, err := bot.apiClient.GetPostThread(context.Background(), postId, "", false)

	if err != nil {
		return nil, err
	}

	postsSlice := postList.ToSlice()
	postMessages := make([]*model.Post, 0, len(postsSlice))

	for _, post := range postsSlice {
		if !IsTDLRCommand(post.Message) {
			postMessages = append(postMessages, post)
		}
	}
	return postMessages, nil
}

func listenToEvents(bot *Bot) {
	for {
		bot.wsClient.Listen()
		log.Println("Listening for Mattermost events")

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

	message := post.Message

	if !IsTDLRCommand(message) {
		return
	}

	postList, err := GetThread(bot, post.Id)

	rootID := postList[0].RootId

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
	sendMessage(bot.apiClient, rootID, post.ChannelId, summary)
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

func sendMessage(client *model.Client4, rootPostID string, channelID string, message string) {
	ctx := context.Background()
	reply := model.Post{ChannelId: channelID, RootId: rootPostID, Message: message}
	// post, resp, err := client.CreatePost(ctx, &reply)
	client.CreatePost(ctx, &reply)
	// fmt.Println(post)
	// fmt.Println(resp)
	// fmt.Println(err)
}
