package models

import (
	claude "github.com/potproject/claude-sdk-go"
	openai "github.com/sashabaranov/go-openai"
)

func CreateOllamaMessages(model *Ollama) []map[string]string {
	// Initialize a slice of map[string]string
	messageList := []map[string]string{}

	// Check if Pattern or Context is not empty
	if model.Pattern != "" || model.Context != "" {
		// Construct the map for the system role
		systemMap := map[string]string{
			"role":    "system",
			"content": model.Context + model.Pattern,
		}
		// Append the map to the slice
		messageList = append(messageList, systemMap)
	}

	// Check if Session is not nil
	if model.Session != nil {
		// Append the Session map directly as it matches the expected type
		messageList = append(messageList, model.Session...)
	}
	userMessage := map[string]string{
		"role":    "user",
		"content": model.Message,
	}
	messageList = append(messageList, userMessage)

	return messageList
}

func CreateGroqMessage(grok *Groq) []openai.ChatCompletionMessage {
	messageList := []openai.ChatCompletionMessage{}

	if grok.Pattern != "" || grok.Context != "" {
		systemMap := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: grok.Context + grok.Pattern,
		}
		messageList = append(messageList, systemMap)
	}

	if grok.Session != nil {
		sessionAsCompletionMessage := []openai.ChatCompletionMessage{}
		for _, sess := range grok.Session{
		for role, content := range sess {
			sessionAsCompletionMessage = append(sessionAsCompletionMessage, openai.ChatCompletionMessage{
				Role:    role,
				Content: content,
			})
		}}
		messageList = append(messageList, sessionAsCompletionMessage...)
		messageList = append(messageList, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleSystem,
			Content: grok.Context + grok.Pattern,
		})
	}

	userMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: grok.Message,
	}
	messageList = append(messageList, userMessage)

	return messageList
}


func CreateOaiMessage(oai *Openai) []openai.ChatCompletionMessage {
	messageList := []openai.ChatCompletionMessage{}

	if oai.Pattern != "" || oai.Context != "" {
		systemMap := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: oai.Context + oai.Pattern,
		}
		messageList = append(messageList, systemMap)
	}

	if oai.Session != nil {
		sessionAsCompletionMessage := []openai.ChatCompletionMessage{}
		for _, sess := range oai.Session{
		for range sess {
			sessionAsCompletionMessage = append(sessionAsCompletionMessage, openai.ChatCompletionMessage{
				Role:    sess["Role"],
				Content: sess["Content"],
			})
		
		}}
		messageList = append(messageList, sessionAsCompletionMessage...)
		messageList = append(messageList, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleSystem,
			Content: oai.Context + oai.Pattern,
		})
	}

	userMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: oai.Message,
	}
	messageList = append(messageList, userMessage)

	return messageList
}

func CreateClaudeMessage(ant *Anthropic) []claude.RequestBodyMessagesMessages {
	messageList := []claude.RequestBodyMessagesMessages{}

	if ant.Session != nil {
		for _, ses := range ant.Session {
		for role, content := range ses {
			messageList = append(messageList, claude.RequestBodyMessagesMessages{
				Role:    role,
				Content: content,
			})
		}
	}}

	userMessage := claude.RequestBodyMessagesMessages{
		Role:    claude.MessagesRoleUser,
		Content: ant.Message,
	}
	messageList = append(messageList, userMessage)

	return messageList
}