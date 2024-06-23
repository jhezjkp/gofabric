package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Ollama struct {
	DefaultModel
}

type ResponseData struct {
	Model      string `json:"model"`
	CreatedAt  string `json:"created_at"`
	Message    MessageData `json:"message"`
	Done       bool `json:"done"`
}

type MessageData struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaModels struct {
    Models []OllamaModelsInner `json:"models"`
}

type OllamaModelsInner struct {
    Name string `json:"name"`
    Model string `json:"model"`
    Modified_at string `json:"modified_at"`
    Size int64 `json:"size"`
    Digest string `json:"digest"`
    Details OllamaDetails `json:"details"`
}

type OllamaDetails struct {
    Parent_model string `json:"parent_model"`
    Format string `json:"format"`
    Family string `json:"family"`
    Families []string `json:"families"`
    Parameter_size string `json:"parameter_size"`
    Quantization_level string `json:"quantization_level"`
}

func NewOllama(url string, message string, pattern string, context string, model string, temperature float64, topP float64, presencePenalty float64, FrequencyPenalty float64, session []map[string]string, responseChan chan string) *Ollama{
    return &Ollama{
        DefaultModel{
            Message: message,
            Pattern: pattern,
            Context: context,
            Model: model,
            Url: url,
            Temperature: temperature,
            TopP: topP,
            PresencePenalty: presencePenalty,
            FrequencyPenalty: FrequencyPenalty,
            Session: session,
            ResponseChan: responseChan,

        },
    }
}

// returns the message or an error
func (ollama *Ollama) SendMessage() (string, error) {
    if ollama.Context != "" {
        ollama.Context = "CONTEXT:\n" + ollama.Context + "\n" // sets context to CONTEXT:\n[context]
    }
    finalMessage := ""
    url := ollama.Url + "/api/chat"
    messages := CreateOllamaMessages(ollama)

    payload := map[string]interface{}{
        "model": ollama.Model,
        "messages": messages,
        "options": map[string]float64{
            "temperature":       ollama.Temperature,
            "presence_penalty":  ollama.PresencePenalty,
            "frequency_penalty": ollama.FrequencyPenalty,
            "top_p":             ollama.TopP,
        },
    }

    requestBody, err := json.Marshal(payload)
    if err != nil {
        return "", err
    }

    client := &http.Client{}
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
    if err != nil {
        return "", err
    }

    req.Header.Add("Content-Type", "application/json")
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    // Preprocess response body to form a valid JSON array
    validJson := "[" + strings.Replace(string(responseBody), "}\n{", "},{", -1) + "]"

    var responses []ResponseData
    err = json.Unmarshal([]byte(validJson), &responses)
    if err != nil {
        return "", fmt.Errorf("json unmarshaling error: %v", err)
    }

    for _, response := range responses {
        finalMessage += response.Message.Content
    }

    return finalMessage, nil
}

func (ollama *Ollama) StreamMessage() (error) {

    if ollama.Context!= "" {
        ollama.Context = "CONTEXT:\n" + ollama.Context + "\n"
    }
    url := ollama.Url + "/api/chat"
    messages := CreateOllamaMessages(ollama)
    payload := map[string]interface{}{
        "model": ollama.Model,
        "messages": messages,
        "options": map[string]float64{
            "temperature":       ollama.Temperature,
            "presence_penalty":  ollama.PresencePenalty,
            "frequency_penalty": ollama.FrequencyPenalty,
            "top_p":             ollama.TopP,
        },
    }
    requestBody, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    client := &http.Client{}
    req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
    if err != nil {
        return err
    }

    req.Header.Add("Content-Type", "application/json")
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    reader := bufio.NewReader(resp.Body)
    var buffer bytes.Buffer
    var responses []ResponseData

    for {
        line, err := reader.ReadBytes('\n') // Assumes that each JSON object ends with a newline
        if err == io.EOF {
            ollama.ResponseChan <- "\n"
            close(ollama.ResponseChan)
            return nil
        }
        if err != nil {
            return err // Handle other errors
        }

        // Attempt to unmarshal each line as a JSON object
        buffer.Write(line)
        validJson := "[" + strings.Replace(string(buffer.String()), "}\n{", "},{", -1) + "]"
        if err := json.Unmarshal([]byte(validJson), &responses); err != nil {
            continue // If unmarshaling fails, it might be due to incomplete data, continue reading
        }

        // Process the successfully unmarshaled messages
        for _, response := range responses {
           ollama.ResponseChan <- response.Message.Content
        }
        buffer.Reset() // Clear the buffer once the data is processed
    }
}

func (ollama Ollama) ListModels()([]string, error) {
    var finalModels []string
	url := ollama.Url + "/api/tags"
	resp, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
        var jsonAsJson OllamaModels
		jsonString := scanner.Text()
        err := json.Unmarshal([]byte(jsonString), &jsonAsJson)
        if err != nil {
            return []string{}, err
        }
        for _, model := range jsonAsJson.Models {
            finalModels = append(finalModels, model.Model)
        }
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading the response body: %s\n", err)
	}
	return finalModels, nil
}