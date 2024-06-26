package chat

import (
	"errors"
	"sync"

	"github.com/xssdoctor/gofabric/models"
	"github.com/xssdoctor/gofabric/utils"
)

// chat struct. This is passed to the SendMessageToModel function to interact with the models
type Chat struct {
	Message          string
	Pattern          string
	OpenAIApiKey     string
	AnthropicApiKey  string
	GroqApiKey       string
	GoogleApiKey     string
	Session          []map[string]string
	Context          string
	Model            string
	OllamaUrl        string
	Temperature      float64
	TopP             float64
	PresencePenalty  float64
	FrequencyPenalty float64
	Stream           bool
	ResponseChan     chan string
}

// the following functions are meant to be used with the Model interface in order to interact with any of the three models
func SendMessage(model Model) (string, error) {
	return model.SendMessage()
}

func StreamMessage(model Model) error {
	return model.StreamMessage()
}

func ListModels(model Model) ([]string, error) {
	return model.ListModels()
}

// ListAllModels returns a map of all models and any errors that occurred Uses concurrence to make it faster
func ListAllModels(chat Chat) (map[string][]string, []error) {
	var wg sync.WaitGroup
	// make the channels
	openaiModelsChan := make(chan []string, 1)
	claudeModelsChan := make(chan []string, 1)
	ollamaModelsChan := make(chan []string, 1)
	groqModelsChan := make(chan []string, 1)
	googleModelsChan := make(chan []string, 1)
	errorsChan := make(chan error, 5)

	wg.Add(5) // We have five concurrent operations.
	errs := make([]error, 0)
	// create a map to store the models. this is used to check if the model is in the list of available models
	modelsMap := make(map[string][]string, 4)
	// add the api keys to the structs, these values come from the chat struct
	openai := &models.Openai{}
	openai.ApiKey = chat.OpenAIApiKey

	claude := &models.Anthropic{}
	claude.ApiKey = chat.AnthropicApiKey

	ollama := &models.Ollama{}
	ollama.Url = chat.OllamaUrl

	groq := &models.Groq{}
	groq.ApiKey = chat.GroqApiKey

	google := &models.Gemini{}
	google.ApiKey = chat.GoogleApiKey

	// create goroutines to list the models for each of the services. function is defined below
	createGoroutines(&wg, openai, errorsChan, openaiModelsChan)

	createGoroutines(&wg, claude, errorsChan, claudeModelsChan)

	createGoroutines(&wg, ollama, errorsChan, ollamaModelsChan)

	createGoroutines(&wg, groq, errorsChan, groqModelsChan)

	createGoroutines(&wg, google, errorsChan, googleModelsChan)

	wg.Wait() // Wait for all goroutines to finish
	close(errorsChan)

	for err := range errorsChan {
		errs = append(errs, err)
	}
	// get the models from the channels
	openaiModels := <-openaiModelsChan
	if openaiModels != nil {
		modelsMap["openai"] = openaiModels
	}

	claudeModels := <-claudeModelsChan
	if claudeModels != nil {
		modelsMap["claude"] = claudeModels
	}

	ollamaModels := <-ollamaModelsChan
	if ollamaModels != nil {
		modelsMap["ollama"] = ollamaModels
	}

	groqModels := <-groqModelsChan
	if groqModels != nil {
		modelsMap["groq"] = groqModels
	}

	googleModels := <-googleModelsChan
	if googleModels != nil {
		modelsMap["google"] = googleModels
	}
	return modelsMap, errs
}

// this is the main function of the app. it takes a chat struct and sends the message to the model with the correct parameters
func (chat Chat) SendMessageToModel() (string, error) {
	stream := chat.Stream
	var activeModel Model
	modelsMap, _ := ListAllModels(chat)
	openAiModels := modelsMap["openai"]
	claudeModels := modelsMap["claude"]
	ollamaModels := modelsMap["ollama"]
	groqModels := modelsMap["groq"]
	googleModels := modelsMap["google"]

	// check if the model is in the list of available models, if so, create a new instance of the model. thi is how the app knows which api to use based on the users choice of model
	if utils.ExistsInArray(chat.Model, openAiModels) {
		activeModel = models.NewOpenai(chat.OpenAIApiKey, chat.Message, chat.Pattern, chat.Context, chat.Model, chat.Temperature, chat.TopP, chat.PresencePenalty, chat.FrequencyPenalty, chat.Session, chat.ResponseChan)
	} else if utils.ExistsInArray(chat.Model, claudeModels) {
		activeModel = models.NewClaude(chat.AnthropicApiKey, chat.Message, chat.Pattern, chat.Context, chat.Model, chat.Temperature, chat.TopP, chat.Session, chat.ResponseChan)
	} else if utils.ExistsInArray(chat.Model, ollamaModels) {
		activeModel = models.NewOllama(chat.OllamaUrl, chat.Message, chat.Pattern, chat.Context, chat.Model, chat.Temperature, chat.TopP, chat.PresencePenalty, chat.FrequencyPenalty, chat.Session, chat.ResponseChan)
	} else if utils.ExistsInArray(chat.Model, groqModels) {
		activeModel = models.NewGroq(chat.GroqApiKey, chat.Message, chat.Pattern, chat.Context, chat.Model, chat.Temperature, chat.TopP, chat.PresencePenalty, chat.FrequencyPenalty, chat.Session, chat.ResponseChan)
	} else if utils.ExistsInArray(chat.Model, googleModels) {
		activeModel = models.NewGemini(chat.GoogleApiKey, chat.Message, chat.Pattern, chat.Context, chat.Model, chat.Temperature, chat.TopP, chat.Session, chat.ResponseChan)
	} else {
		return "", errors.New("Model not found")
	}
	if stream {
		err := StreamMessage(activeModel)
		if err != nil {
			chat.ResponseChan <- err.Error()
		}
		return "", nil

	} else {
		return SendMessage(activeModel)
	}

}

// helper fnction which creates goroutines to list the models for each of the services
func createGoroutines(wg *sync.WaitGroup, model Model, errorsChan chan error, modelChan chan []string) {
	go func() {
		defer wg.Done()
		models, err := ListModels(model)
		if err != nil {
			errorsChan <- err
			//fmt.Println("Error listing models", reflect.TypeOf(model).Elem().Name())
			modelChan <- nil
		} else {
			modelChan <- models
		}
	}()
}
