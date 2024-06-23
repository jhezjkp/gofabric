package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xssdoctor/gofabric/chat"
)

var DB *sql.DB

type Entry struct {
	Id int
	Name string
	Description string
	Pattern string
	Context string
	Session string
	Openai_api_key string
	Anthropic_api_key string
	Groq_api_key string
	Ollama_url string
	Google_api_key string
	Default_model string
}

// this function always runs when the program is started. It creates the patterns directory if not there and downloads the folders from the internet
func InitDB() error {
    var err error
    home_dir, _ := os.UserHomeDir()
    fabric_config := home_dir + "/.config/fabric"
	patterns_dir := fabric_config + "/patterns"

    // Create the directory if it doesn't exist
    err = os.MkdirAll(fabric_config, os.ModePerm) // creates the directory ~/.confib/fabric if it doesn't exist
    if err != nil {
        return fmt.Errorf("failed to create directory: %v", err)
    }
	err = os.MkdirAll(patterns_dir, os.ModePerm) // creates the directory ~/.confib/fabric/patterns if it doesn't exist
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	_, folderError := os.Stat(patterns_dir)
    if os.IsNotExist(folderError) { // checks if the database file exists, if not, it runs the following code
        createTables() // creates the tables in the database. this function is in db_functions.go
        err := InitialRun() // runs the InitialRun function in setup.go. this is defined below
        if err != nil {
            return err
        }
        err = PopulateDB() // populates the database with patterns that are downloaded from the internet. this function is in setup.go
        if err != nil {
            return err
        }
    }

    return nil
}

// runs the initial setup of the program. this includes entering the api keys and default models
func InitialRun() (error) {
	var openAi, anthropic, ollama, model, groq, google string
	// enters the OpenAI, Anthropic, Groq api keys, and ollama base url
	e := Entry{}
	fmt.Println("Enter your OpenAI API key: (Leave blank if you don't have one)")
	fmt.Scanln(&openAi)
	e.Openai_api_key = strings.TrimRight(openAi, "\n")
	fmt.Println("Enter your Anthropic API key: (Leave blank if you don't have one)")
	fmt.Scanln(&anthropic)
	e.Anthropic_api_key = strings.TrimRight(anthropic, "\n")
	fmt.Println("Enter your google API key: (Leave blank if you don't have one)")
	fmt.Scanln(&google)
	e.Google_api_key = strings.TrimRight(google, "\n")
	fmt.Println("Enter your Groq API key: (Leave blank if you don't have one)")
	fmt.Scanln(&groq)
	e.Groq_api_key = strings.TrimRight(groq, "\n")
	fmt.Println("Enter your Ollama URL: (leave blank if you don't have one or if you want the default of localhost:11434)")
	fmt.Scanln(&ollama)
	e.Ollama_url = strings.TrimRight(ollama, "\n")
	if e.Ollama_url == "" {
		e.Ollama_url = "http://127.0.0.1:11434" // this is the default value for the ollama url
	
	}
	fmt.Println()
	fmt.Println()
	chatInstance := chat.Chat{ // creates a blank chat instance with the api keys, the purpose of this is to list all the models
		OpenAIApiKey: e.Openai_api_key,
		AnthropicApiKey: e.Anthropic_api_key,
		OllamaUrl: e.Ollama_url,
		GoogleApiKey: e.Google_api_key,
	}
	models, _ := chat.ListAllModels(chatInstance) // lists all the models for each of the three services, returns a map[string]string of the models
	for key, value := range models {
		fmt.Println(key)
		fmt.Println()
		for _, model := range value {
			fmt.Println(model)
		}
		fmt.Println()
	
	}
	// from the listed models, the user is prompted to choose a default model
	fmt.Println("Enter your default model: Choose from the above options")
	fmt.Scanln(&model)
	e.Default_model = strings.TrimRight(model, "\n")
	err := e.InsertConfiguration() // takes the Entry struct which includes all relivant api keys and default models and inserts it into the database
	if err != nil {
		return errors.New("could not insert configuration")
	}
	return nil
}

