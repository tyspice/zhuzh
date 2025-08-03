package models

type ChatClient interface {
	Subscribe() (<-chan string, <-chan error)
	Ask(prompt string)
	Close()
}
