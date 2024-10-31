package main

import (
	"context"
	"github.com/mattermost/mattermost/server/public/model"
	"os"
)

func GetThread(postId string) ([]string, error) {
	// Provides a list of post messages from the same thread as `postId`, in reverse order.
	client := model.NewAPIv4Client("https://chat.canonical.com")
	client.SetToken(os.Getenv("MM_AUTHTOKEN"))

	postList, _, err := client.GetPostThread(context.Background(), postId, "", false)

	if err != nil {
		return nil, err
	}

	postMessages := make([]string, len(postList.Posts))
	for _, post := range postList.Posts {
		postMessages = append(postMessages, post.Message)
	}

	return postMessages, nil
}
