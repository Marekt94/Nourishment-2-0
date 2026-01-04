package AIClient

import (
	"context"
	"encoding/json"
	"nourishment_20/internal/logging"

	"github.com/revrost/go-openrouter"
)

type OpenRouterClient struct {
	ApiKey    string
	Model     string
	MaxTokens int
}

func (c *OpenRouterClient) ExecutePrompt(p string, s AIResponseSchemaIntf) (string, bool) {
	client := openrouter.NewClient(c.ApiKey)
	request := openrouter.ChatCompletionRequest{
		Model: c.Model,
		Messages: []openrouter.ChatCompletionMessage{
			{
				Role:    openrouter.ChatMessageRoleUser,
				Content: openrouter.Content{Text: p},
			},
		},

		MaxTokens: c.MaxTokens,
	}
	if s != nil {
		request.ResponseFormat = &openrouter.ChatCompletionResponseFormat{
			Type: openrouter.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openrouter.ChatCompletionResponseFormatJSONSchema{
				Name:   "meal_json_schema",
				Strict: true,
				Schema: s,
			},
		}
	}
	requestStr, err := json.MarshalIndent(request, "", "\t")
	if err != nil {
		logging.Global.Panicf("Error marshaling request: %v", err)
	}
	logging.Global.Debugf("Request: %s", requestStr)
	ctx := context.Background()
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		logging.Global.Panicf("Error while requesting openrouter: %v", err)
	}
	responseStr, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		logging.Global.Panicf("Error marshaling response: %v", err)
	}
	logging.Global.Debugf("Response: %s", responseStr)
	requestFinishedSuccessfully := response.Choices[0].FinishReason == openrouter.FinishReasonStop
	return response.Choices[0].Message.Content.Text, requestFinishedSuccessfully
}
