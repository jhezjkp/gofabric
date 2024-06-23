package chat

// this interfact allows any of the models to be passed to the chat instance
type Model interface {
	SendMessage() (string, error)
	StreamMessage() (error)
	ListModels() ([]string, error)
}