package models

import (
	"context"
	"errors"
	"fmt"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

// create Openai struct
type Openai struct {
	DefaultModel
}

func NewOpenai(apiKey string, message string, pattern string, context string, model string, temperature float64, topP float64, presencePenalty float64, FrequencyPenalty float64, session []map[string]string, responseChan chan string) *Openai {
	return &Openai{
		DefaultModel{
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
func (oai *Openai) SendMessage() (string, error){
	// If context is int the Openai struct, contextMessage will be CONTEXT:\n[context], otherwise contextMessage will be ""
    if oai.Context != "" {
        oai.Context = "CONTEXT:\n" + oai.Context + "\n" // set context to "CONTEXT:\n[context]"
    }
	// gives default values for Temperature, TopP, PresencePenalty and FrequencyPenalty if not mentioned
	client := openai.NewClient(oai.ApiKey)
	messages := CreateOaiMessage(oai)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: oai.Model,
			Temperature: float32(oai.Temperature),
			TopP: float32(oai.TopP),
			PresencePenalty: float32(oai.PresencePenalty),
			FrequencyPenalty:float32(oai.FrequencyPenalty),
			Messages: messages,
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// streams message AND yields a message and an error for futher processing if necessary
func (oai *Openai) StreamMessage() (error) {
	// If context is int the Openai struct, contextMessage will be CONTEXT:\n[context], otherwise contextMessage will be ""
    if oai.Context != "" {
        oai.Context = "CONTEXT:\n" + oai.Context + "\n" // set context to CONTEXT\n[context]
    }
	c := openai.NewClient(oai.ApiKey)
	messages := CreateOaiMessage(oai)
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model:     oai.Model,
		Temperature: float32(oai.Temperature),
		TopP: float32(oai.TopP),
		PresencePenalty: float32(oai.PresencePenalty),
		FrequencyPenalty: float32(oai.FrequencyPenalty),
		Messages: messages,
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
			oai.ResponseChan <- "\n"
			close(oai.ResponseChan)
			return nil
			
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return err
		}
		oai.ResponseChan <- response.Choices[0].Delta.Content
	}
}

	// returns a list of all available openai models
func (oai *Openai)ListModels() ([]string, error) {
	var modelList []string
	ctx := context.Background()
	client := openai.NewClient(oai.ApiKey)
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