package models

import (
	"context"
	"errors"
	"fmt"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

// create Openai struct
type Groq struct {
	DefaultModel
}

func NewGroq(apiKey string, message string, pattern string, context string, model string, temperature float64, topP float64, presencePenalty float64, FrequencyPenalty float64, session []map[string]string, responseChan chan string) *Groq {
	return &Groq{
		DefaultModel: DefaultModel{
			Message: message,
			Pattern: pattern,
			Context: context,
			Model: model,
			ApiKey: apiKey,
			Temperature: temperature,
			TopP: topP,
			PresencePenalty: presencePenalty,
			FrequencyPenalty: FrequencyPenalty,
			Session: session,
			ResponseChan: responseChan,
		},

	}
}

// creates a Sendmessage method which yields the message or an error
func (Groq *Groq) SendMessage() (string, error){
	// If context is int the Openai struct, contextMessage will be CONTEXT:\n[context], otherwise contextMessage will be ""
    if Groq.Context != "" {
        Groq.Context = "CONTEXT:\n" + Groq.Context + "\n" // set context to "CONTEXT:\n[context]"
    }
	// gives default values for Temperature, TopP, PresencePenalty and FrequencyPenalty if not mentioned
    config := openai.DefaultConfig(Groq.ApiKey)
    config.BaseURL = "https://api.groq.com/openai/v1"
	client := openai.NewClientWithConfig(config)
	messages := CreateGroqMessage(Groq)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: Groq.Model,
			Temperature: float32(Groq.Temperature),
			TopP: float32(Groq.TopP),
			PresencePenalty: float32(Groq.PresencePenalty),
			FrequencyPenalty:float32(Groq.FrequencyPenalty),
			Messages: messages,
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// streams message AND yields a message and an error for futher processing if necessary
func (Groq Groq) StreamMessage() (error) {
	// If context is int the Openai struct, contextMessage will be CONTEXT:\n[context], otherwise contextMessage will be ""
    if Groq.Context != "" {
        Groq.Context = "CONTEXT:\n" + Groq.Context + "\n" // set context to CONTEXT\n[context]
    }
	config := openai.DefaultConfig(Groq.ApiKey)
    config.BaseURL = "https://api.groq.com/openai/v1"
	c := openai.NewClientWithConfig(config)
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model:     Groq.Model,
		Temperature: float32(Groq.Temperature),
		TopP: float32(Groq.TopP),
		PresencePenalty: float32(Groq.PresencePenalty),
		FrequencyPenalty: float32(Groq.FrequencyPenalty),
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: Groq.Context + Groq.Pattern, // if no context or pattern are passed as arguments, this will evaluate to ""
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: Groq.Message,
			},
		},
		Stream: true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return err
	}
	defer stream.Close()
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			Groq.ResponseChan <- "\n"
			close(Groq.ResponseChan)
			return nil
			
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return err
		}
		Groq.ResponseChan <- response.Choices[0].Delta.Content
	}
}

	// returns a list of all available openai models
func (Groq Groq)ListModels() ([]string, error) {
	var modelList []string
	ctx := context.Background()
    config := openai.DefaultConfig(Groq.ApiKey)
    config.BaseURL = "https://api.groq.com/openai/v1"
	client := openai.NewClientWithConfig(config)
	modelsTemp, err := client.ListModels(ctx)
	if err != nil {
		return []string{}, err
	}
	model := modelsTemp.Models
	for _, mod := range model {
		modelList = append(modelList, mod.ID)
	}
	return modelList, nil
}