package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/xssdoctor/gofabric/chat"
	"github.com/xssdoctor/gofabric/db"
	"github.com/xssdoctor/gofabric/flags"
)


func copyToClipboard(message string) error {
	err := clipboard.WriteAll(message)
	if err != nil {
		return errors.New("could not copy to clipboard")
	}
	return nil
}

func latestPatterns(latestNumber int) error {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return errors.New("could not get home directory")
	}
	unique_patterns_file := filepath.Join(home_dir, ".config/fabric/unique_patterns.txt")
	contents, err := os.ReadFile(unique_patterns_file)
	if err != nil {
		return errors.New("could not read unique patterns file. Pleas run --updatepatterns")
	}
	unique_patterns := strings.Split(string(contents), "\n")
	for i := len(unique_patterns) - 1; i > len(unique_patterns)-latestNumber-1; i-- {
		fmt.Println(unique_patterns[i])
	}
	return nil

}

func createOutputFile(message string, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file")
	}
	defer file.Close()
	_, err = file.WriteString(message)
	if err != nil {
		fmt.Println("Error writing to file")
	}
}

func ContextAdd() error {
	var name, description, FileName string
	fmt.Println("Enter the name of the context")
	fmt.Scanln(&name)
	if strings.Contains(name, " ") {
		return errors.New("pattern name must be a single word")
	}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter your description, then press ENTER:")
	if scanner.Scan() {
		description = scanner.Text()
	}
	fmt.Println("Enter the file path of the context that you created")
	fmt.Scanln(&FileName)
	file, err := os.Open(FileName)
	if err != nil {
		return errors.New("could not open file")
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return errors.New("could not read file")
	}
	e := db.Entry{
		Name:        name,
		Description: description,
		Pattern:     string(content),
	}
	err = e.InsertContext()
	if err != nil {
		return err
	}
	fmt.Println("Context added successfully")
	return nil
}

func listAllPatterns() error {
	patterns, err := db.ListAllPatterns()
	if err != nil {
		return err
	}
	for _, pattern := range patterns {
		fmt.Println(pattern.Name, pattern.Description)
	}
	return nil
}

func listAllContexts() error {
	contexts, err := db.ListAllContexts()
	if err != nil {
		return err
	}
	for _, context := range contexts {
		fmt.Println(context.Name, context.Description)
	}
	return nil

}

func listAllModels() error {
	config, err := db.GetConfiguration()
	if err != nil {
		return err
	}
	ch := chat.Chat{
		OpenAIApiKey: config.Openai_api_key,
		AnthropicApiKey: config.Anthropic_api_key,
		OllamaUrl: config.Ollama_url,
		GroqApiKey: config.Groq_api_key,
		GoogleApiKey: config.Google_api_key,
	}
	models, _ := chat.ListAllModels(ch)
	for modelType, modelList := range models {
		fmt.Println(modelType)
		fmt.Println()
		for _, model := range modelList {
			fmt.Println(model)
		}
		fmt.Println()
	}
	return nil
}

func ListSessions() error {
	sessions, err := db.ListAllSessions()
	if err != nil {
		return err
	}
	for _, session := range sessions {
		fmt.Println(session.Name)
	}
	return nil
}

func initiateChat(flags flags.Flags) (string, error) {
	var activeModel string
	config, err := db.GetConfiguration()
	if err != nil {
		if err.Error() == "Could not get configuration from database" {
			return "", errors.New("could not get configuration from database. Please run -setup again and input your api keys, url and default model")
		}
		return "", err
	}
	if flags.Model == "" && config.Default_model == "" {
		activeModel = "gpt-4-turbo-preview"
	} else if flags.Model == "" {
		activeModel = config.Default_model
	} else {
		activeModel = flags.Model
	}
	if flags.Url == "" {
		flags.Url = config.Ollama_url
	}
	if flags.Pattern != "" {
		e := db.Entry{
			Name: flags.Pattern,
		}
		r, err := e.GetPatternByName()
		if err != nil {
			return "", errors.New("could not find pattern")
		}
		if r.Pattern != "" {
			flags.Pattern = r.Pattern
		}
	}
	var session []map[string]string
	if flags.Session != "" {
		ses, err := getSession(flags.Session)
		if err != nil {
			fmt.Println("Creating new session")
		} else {
			session = ses
		
		}

	}
	activeChat := chat.Chat{
		Message:          flags.Message,
		Pattern:          flags.Pattern,
		Context:          flags.Context,
		Model:            activeModel,
		Stream: 		 flags.Stream,
		OllamaUrl:        flags.Url,
		Temperature: 	  flags.Temperature,
		TopP:			flags.TopP,
		PresencePenalty: flags.PresencePenalty,
		FrequencyPenalty: flags.FrequencyPenalty,
		OpenAIApiKey: config.Openai_api_key,
		AnthropicApiKey: config.Anthropic_api_key,
		GroqApiKey: config.Groq_api_key,
		GoogleApiKey: config.Google_api_key,
		Session: session,
		ResponseChan: make(chan string),

	}
	message := ""
	if flags.Stream {
		go func() {
            _, err := activeChat.SendMessageToModel()
            if err != nil {
                activeChat.ResponseChan <- err.Error()
            }
        }()
        // fmt.printll evetying coming from the response channel
        for response := range activeChat.ResponseChan {
            message += response
            fmt.Print(response)
        }
	} else {
		message, err = activeChat.SendMessageToModel()
		if err != nil {
			return "", err
		}
	}
	if flags.Session != "" {
		err = UpdateSession(flags.Session, flags.Message, message)
		if err != nil {
			return "", err
		}
	}
	return message, err
}