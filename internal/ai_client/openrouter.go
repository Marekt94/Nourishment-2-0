package aiclient

import (
	"context"
	"encoding/json"
	"nourishment_20/internal/logging"

	"github.com/revrost/go-openrouter"
)

type OpenRouterClient struct {
	ApiKey string
	Model  string
}

func (c *OpenRouterClient)ExecutePrompt(p string) {
	client := openrouter.NewClient(c.ApiKey);
	request := openrouter.ChatCompletionRequest{
		Model: c.Model,
		Messages: []openrouter.ChatCompletionMessage{
			{
				Role:    openrouter.ChatMessageRoleUser,
				Content: openrouter.Content{Text: p},
			},
		},
		MaxTokens: 1000, //TODO: make it configurable
	}
	requestStr, err := json.MarshalIndent(request, "", "\t")
	if err != nil {
		logging.Global.Panicf("Error marshaling request: %v", err)
	}
	logging.Global.Debugf("Request: %s", requestStr)
	ctx := context.Background();
	response, err := client.CreateChatCompletion(ctx, request);
	if err != nil {
		logging.Global.Panicf("Error while requesting openrouter: %v", err)
	}
	logging.Global.Debugf("Response: %s", response)
}