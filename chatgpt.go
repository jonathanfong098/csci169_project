package main

// import (
// 	"github.com/ayush6624/go-chatgpt"
// )

// // key := os.Getenv("OPENAI_KEY")

// // client, err := chatgpt.NewClient(key)
// // if err != nil {
// // 	log.Fatal(err)
// // }

// type chatgptClient struct {
// 	client *chatgpt.Client
// }

// func (client *chatgptClient) recommendFeed(interest: string) string {
// 	message := "I'm interested in " + interestString + ". Can you recommend RSS feeds to subscribe to?"
// 	res, err := c.SimpleSend(ctx, message)
// 	if err != nil {
// 		return "Sorry, I couldn't recommend a feed."
// 	}
// 	return response.Choices[0].Text
// }
// func summarizePost(posts: string[]) string {
// 	result = ""
// 	for post in posts {
// 		message := "Summarize the post in 3 bullet points " + post
// 		res, err := c.SimpleSend(ctx, message)
// 		if err != nil {
// 			return "Sorry, I couldn't summarize that post."
// 		}
// 		result += response.Choices[0].Text
// 	}
// 	return result
// }
