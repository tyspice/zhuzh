package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

func Run() {
	raw, err := os.ReadFile("mock-data/viewport-text.md")
	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	pretty, err := glamour.Render(string(raw), "dark")
	if err != nil {
		fmt.Println("could not render markdown", err)
		os.Exit(1)
	}

	p := tea.NewProgram(
		model{content: pretty},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("Oopsie:", err)
		os.Exit(1)
	}
}
