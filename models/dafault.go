package models

// the default struct that the models will be based on
type DefaultModel struct {
	Message string
	Pattern string
    ApiKey string
	Context string
	Session []map[string]string
	Model   string
	Url     string
    Temperature float64
	TopP float64
	PresencePenalty float64
	FrequencyPenalty float64
	ResponseChan chan string
}