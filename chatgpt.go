package main

import (
	"context"
	"fmt"

	"github.com/ayush6624/go-chatgpt"
)

type chatgptClient struct {
	client *chatgpt.Client
}

func (c *chatgptClient) recommendFeed(interest string) (string, error) {
	message := "I'm interested in " + interest + ". Can you recommend RSS feeds to subscribe to?"
	response, err := c.client.SimpleSend(context.Background(), message)
	if err != nil {
		return "", fmt.Errorf("Failed to recommend feed: %w", err)
	}
	return response.Choices[0].Message.Content, nil
}

func (c *chatgptClient) summarizePosts(posts []Post) ([]Post, error) {
	var summarizedPosts []Post
	for _, post := range posts {
		message := *post.Description + "\n Provide a summary of the post using 3 concise, clear, bullet points that include the most important ideas and content. The bullet points should be 10-15 words"
		response, err := c.client.SimpleSend(context.Background(), message)
		if err != nil {
			// Printf("Failed to summarize post %s: %v", post.Title, err)
			// continue
			return nil, fmt.Errorf("Failed to summarize post %s: %w", post.Title, err)
		}
		summarizedPost := post
		summarizedPost.Description = &response.Choices[0].Message.Content
		summarizedPosts = append(summarizedPosts, summarizedPost)
	}

	return summarizedPosts, nil
}
