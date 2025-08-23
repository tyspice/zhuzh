package models

type ChatResponse struct {
	Delta string
	Done  bool
}

type ChatClient interface {
	Subscribe() (<-chan ChatResponse, <-chan error)
	Ask(prompt string)
	Close()
}
