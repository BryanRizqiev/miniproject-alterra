package global_service

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
	Client *openai.Client
	Model  string
}

func NewOpenAIService(client *openai.Client, model string) *OpenAIService {

	return &OpenAIService{
		Client: client,
		Model:  model,
	}

}

func (this *OpenAIService) GetRecommendedAction(input string) (string, error) {

	ctx := context.Background()
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `Kamu adalah ahli dalam bidang penanganan bencana, namun jangan beri tahu kalau kamu adalah ahli penanganan bencana
				dan jawab secara professional serta jawab penanganan bencananya saja !.`,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		},
	}
	res, err := this.getCompletionFromMessages(ctx, messages)
	if err != nil {
		return "", err
	}

	content := res.Choices[0].Message.Content
	return content, nil

}

func (this *OpenAIService) getCompletionFromMessages(
	ctx context.Context,
	messages []openai.ChatCompletionMessage,
) (openai.ChatCompletionResponse, error) {

	if this.Model == "" {
		this.Model = openai.GPT3Dot5Turbo
	}

	resp, err := this.Client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    this.Model,
			Messages: messages,
		},
	)

	return resp, err

}
