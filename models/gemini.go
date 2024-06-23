package models

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Gemini struct {
	DefaultModel
}

func NewGemini(apiKey string, message string, pattern string, context string, model string, temperature float64, topP float64, session []map[string]string, responseChan chan string) *Gemini {
	if pattern == "" {
		pattern = " "
	}
	if context == "" {
		context = " "
	}
	return &Gemini{
		DefaultModel{
			ApiKey: apiKey,
			Message: message,
			Pattern: pattern,
			Context: context,
			Model: model,
			Session: session,
			ResponseChan: responseChan,
		},
	}
}

func (gem *Gemini) SendMessage() (string, error) {
	finalResponse := ""
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(gem.ApiKey))
	if err != nil {
		return "", err
	}
	defer client.Close()
	model := client.GenerativeModel(gem.Model)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Part(genai.Text(gem.Context + gem.Pattern)),
		},
	}
	response, err := model.GenerateContent(ctx, genai.Text(gem.Message))
	if err != nil {
		return "", err
	}
	for _, cand := range response.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if text, ok := part.(genai.Text); ok {
					finalResponse += string(text)
				}
				
			}
		}
	}
	return finalResponse, nil
}

func (gem *Gemini) StreamMessage() (error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(gem.ApiKey))
	if err != nil {
		return err
	}
	defer client.Close()
	model := client.GenerativeModel(gem.Model)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Part(genai.Text(gem.Context + gem.Pattern)),
		},
	}
	iter := model.GenerateContentStream(ctx, genai.Text(gem.Message))
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			gem.ResponseChan <- "\n"
			close(gem.ResponseChan)
			return nil
		}
		if err != nil {
			return err
		}
		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					if text, ok := part.(genai.Text); ok {
						gem.ResponseChan <- string(text)
					}
				}
			}
		}
	}
}

func (gem *Gemini) ListModels() ([]string, error) {
	var finalList []string
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(gem.ApiKey))
	if err != nil {
		return []string{}, err
	}
	defer client.Close()
	iter := client.ListModels(ctx)
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []string{}, err
		}
		finalList = append(finalList, resp.Name)
	}
	return finalList, nil
}