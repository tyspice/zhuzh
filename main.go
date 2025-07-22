package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return "\nPress q to quit.\n"
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Uh Oh! There was an error: %v", err)
		os.Exit(1)
	}
}
