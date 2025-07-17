package AIClient

import (
	"context"
	"encoding/json"
	"nourishment_20/internal/logging"

	"github.com/invopop/jsonschema"
	"github.com/revrost/go-openrouter"
)

type OpenRouterClient struct {
	ApiKey    string
	Model     string
	MaxTokens int
}

func (c *OpenRouterClient) ExecutePrompt(p string, s *jsonschema.Schema) (string, bool) {
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
				Name:   "my_json_schema",
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
	//TODO: dorobić zwrotke fals, gdy finish jest inny niz ok (np, tokens)
	return response.Choices[0].Message.Content.Text, true
}
