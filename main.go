package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	textInput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	m := NewModel()

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()

	if err != nil {
		log.Fatal(err)
	}
}

type Model struct {
	input      textInput.Model
	definition string

	width  int
	height int
}

type Definition struct {
	Definition string `json:"definition"`
}

type UrbanResponse struct {
	List []Definition `json:"list"`
}

func NewModel() Model {
	input := textInput.New()
	input.Placeholder = "Press Enter to search, Esc to quit"
	input.Focus()

	return Model{
		input: input,
	}
}

func (m Model) Init() tea.Cmd {
	return textInput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			query := m.input.Value()
			return m, HandleSearch(query)

		case tea.KeyEscape:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case UrbanResponse:
		if len(msg.List) > 0 {
			m.definition = msg.List[0].Definition
		} else {
			m.definition = "No definition found."
		}
	}

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	inputArea := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		PaddingLeft(1).
		PaddingRight(1).
		MarginLeft(1).
		MarginRight(1).
		Width(m.width - 4).
		Render(m.input.View())

	definitionArea := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		PaddingLeft(1).
		PaddingRight(1).
		MarginLeft(1).
		MarginRight(1).
		Height(m.height - 8).
		Width(m.width - 4).
		Render(m.definition)

	footer := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		PaddingLeft(1).
		PaddingRight(1).
		MarginLeft(1).
		MarginRight(1).
		AlignHorizontal(lipgloss.Center).
		Width(m.width - 4).
		Render("Developed and maintained by https://anav.dev")

	return lipgloss.JoinVertical(lipgloss.Top, inputArea, definitionArea, footer)
}

func HandleSearch(query string) tea.Cmd {
	return func() tea.Msg {
		escapedQuery := url.QueryEscape(query)

		res, err := http.Get("https://api.urbandictionary.com/v0/define?term=" + escapedQuery)

		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)

		var urbanRes UrbanResponse

		json.Unmarshal(body, &urbanRes)

		return urbanRes
	}
}
