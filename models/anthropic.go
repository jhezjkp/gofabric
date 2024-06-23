package models

import (
	"context"
	"errors"
	"io"

	"github.com/liushuangls/go-anthropic/v2"
	claude "github.com/potproject/claude-sdk-go"
)

type Anthropic struct {
	DefaultModel
}

func NewClaude(apiKey string, message string, pattern string, context string, model string, temperature float64, topP float64, session []map[string]string, responseChan chan string) *Anthropic {
	return &Anthropic{
		DefaultModel{
			Message: message,
			Pattern: pattern,
			Context: context,
			Model:   model,
			ApiKey:  apiKey,
			Temperature: temperature,
			TopP: topP,
			Session: session,
			ResponseChan: responseChan,
		},
	}

}

func (ant *Anthropic) SendMessage() (string, error) {
    if ant.Context != "" {
        ant.Context = "CONTEXT:\n" + ant.Context + "\n" //set context to CONTEXT\n[context]
    }
	messages := CreateClaudeMessage(ant)
	c := claude.NewClient(ant.ApiKey)
	m := claude.RequestBodyMessages{
		Model:     ant.Model,
		MaxTokens: 4096,
		Temperature: ant.Temperature,
		TopP: ant.TopP,
		System: ant.Context + ant.Pattern,
		Messages: messages,
	}
	ctx := context.Background()
	res, err := c.CreateMessages(ctx, m)
	if err != nil {
		return "", err
	}
	return res.Content[0].Text, nil
}

func (ant *Anthropic) StreamMessage() (error) {
	// streams message and also returns completed message for further functions
    if ant.Context != "" {
        ant.Context = "CONTEXT:\n" + ant.Context + "\n" //set context to CONTEXT:\n[context]
    }
	c := claude.NewClient(ant.ApiKey)
	m := claude.RequestBodyMessages{
		Model:     ant.Model,
		MaxTokens: 4096,
		Temperature: ant.Temperature,
		TopP: ant.TopP,
		System: ant.Context + ant.Pattern,
		Messages: []claude.RequestBodyMessagesMessages{
			{
				Role:    claude.MessagesRoleUser,
				Content: ant.Message,
			},
		},
	}
	ctx := context.Background()
	stream, err := c.CreateMessagesStream(ctx, m)
	if err != nil {
		return err
	}
	defer stream.Close()
	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			ant.ResponseChan <- "\n"
			close(ant.ResponseChan)
			return nil
		}
		if err != nil {
			return err
		}

		asText := res.Content[0].Text
		ant.ResponseChan <- asText
	}
}

func (ant *Anthropic) ListModels() ([]string, error) {
	// returns a list of models. I had to create it myself since the anthropic api doesn't have a ListModels function
	if ant.ApiKey == "" {
		return []string{}, errors.New("no claude api key")
	}
	return []string{anthropic.ModelClaude3Haiku20240307, anthropic.ModelClaude3Opus20240229, anthropic.ModelClaude2Dot0, anthropic.ModelClaude2Dot1, anthropic.ModelClaudeInstant1Dot2, "claude-3-5-sonnet-20240620"}, nil
}