package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tyspice/zhuzh/internal/models"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type responseMsg struct {
	Content string
}

type model struct {
	content    string
	ready      bool
	chatClient models.ChatClient
	viewport   viewport.Model
	textInput  textinput.Model
}

func (m model) Init() tea.Cmd {
	return waitForActivity(m.chatClient)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	updateTextInput := func() {
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	updateViewport := func() {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.content += "\n\n" + "====" + "\n\n"
			m.chatClient.Ask(m.textInput.Value())
			m.textInput.SetValue("")
			wrappedContent := wrapText(m.content, m.viewport.Width)
			m.viewport.SetContent(wrappedContent)
			m.viewport.GotoBottom()
		case "up", "down":
			updateViewport()
		default:
			updateTextInput()
		}

	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown {
			updateViewport()
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight + lipgloss.Height(m.textInput.View())

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent("")
			m.textInput = textinput.New()
			m.textInput.Focus()

			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
			wrappedContent := wrapText(m.content, m.viewport.Width)
			m.viewport.SetContent(wrappedContent)
		}

	case responseMsg:
		m.content += msg.Content
		wrappedContent := wrapText(m.content, m.viewport.Width)
		m.viewport.SetContent(wrappedContent)
		m.viewport.GotoBottom()
		cmds = append(cmds, waitForActivity(m.chatClient))
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		m.headerView(),
		m.viewport.View(),
		m.footerView(),
		m.textInput.View(),
	)
}

func (m model) headerView() string {
	title := titleStyle.Render("zhuzh")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func waitForActivity(c models.ChatClient) tea.Cmd {
	return func() tea.Msg {
		res, _ := c.Subscribe()
		next := <-res
		return responseMsg{Content: next}
	}
}

func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	var wrapped strings.Builder
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		if i > 0 {
			wrapped.WriteString("\n")
		}

		if len(line) <= width {
			wrapped.WriteString(line)
			continue
		}

		// Handle lines longer than width
		currentWidth := 0
		words := strings.Fields(line)
		for j, word := range words {
			wordLen := len(word)

			if j == 0 {
				wrapped.WriteString(word)
				currentWidth = wordLen
			} else if currentWidth+wordLen+1 > width {
				// Start a new line
				wrapped.WriteString("\n" + word)
				currentWidth = wordLen
			} else {
				// Add word to current line
				wrapped.WriteString(" " + word)
				currentWidth += wordLen + 1
			}
		}
	}

	return wrapped.String()
}
