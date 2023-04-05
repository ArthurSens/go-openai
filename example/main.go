package main

import (
	"context"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	ctx := context.Background()
	client := openai.NewPprofClient(os.Getenv("OPENAI_API_KEY"))

	defer func() {
		err := client.WriteToFile("profile.pb.gz")
		if err != nil {
			log.Fatalf("Failed to write profile: %v", err)
		}
	}()
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Hi! I'm trying to create a custom Profiler for OpenAI API in Golang. I want to be able to tell how many tokens were used for each request.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "I've created a pprof.Profile object, and I'm useing profile.Add() to add the number of tokens used for each request. When I visualize the profile, the cumulative only shows me 1 instead of the actual number of tokens used. What am I doing wrong?",
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to complete: %v", err)
	}

	for _, choice := range resp.Choices {
		log.Printf("content=%q, tokensUsed=%d", choice.Message.Content, resp.Usage.TotalTokens)
	}
}
