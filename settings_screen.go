package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state string

const (
	ADD_GITHUB_TOKEN      state = "ADD_GITHUB_TOKEN"
	ADD_GITHUB_REPOSITORY state = "ADD_GITHUB_REPOSITORY"
	DEFAULT               state = "DEFAULT"
)

type SettingsScreenModel struct {
	TextInput textinput.Model
	state     state
	*GlobalState
	*Settings
}

func NewSettingsScreenModel() *SettingsScreenModel {
	textInput := textinput.New()
	textInput.Placeholder = "Type something..."
	textInput.CharLimit = 200
	textInput.Focus()
	textInput.Width = 40

	return &SettingsScreenModel{
		TextInput: textInput,
	}
}

func (r *SettingsScreenModel) WithGlobalState(globalState *GlobalState) *SettingsScreenModel {
	r.GlobalState = globalState

	return r
}

func (r *SettingsScreenModel) WithSettings(settings *Settings) *SettingsScreenModel {
	r.Settings = settings

	return r
}

func (r *SettingsScreenModel) Init() tea.Cmd {
	return nil
}

func (r *SettingsScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.Type {
			case tea.KeyEscape:
				{
					debug.msg(debug.KeyPressed(), "escape")

					if r.state == ADD_GITHUB_TOKEN || r.state == ADD_GITHUB_REPOSITORY {
						r.state = DEFAULT
						r.TextInput.Reset()
					}
				}
			case tea.KeyCtrlU:
				{
					debug.msg(debug.KeyPressed(), "ctrl + U")

					r.state = ADD_GITHUB_TOKEN
				}
			case tea.KeyCtrlR:
				{
					debug.msg(debug.KeyPressed(), "ctrl + R")

					r.state = ADD_GITHUB_REPOSITORY
				}
			case tea.KeyEnter:
				{
					debug.msg(debug.KeyPressed(), "enter")

					switch r.state {
					case ADD_GITHUB_TOKEN:
						{
							debug.msg(debug.UI(), "current state is ADD_GITHUB_TOKEN")
							if r.TextInput.Value() != "" {
								r.TextInput.Reset()
							}

							r.state = DEFAULT
						}
					case ADD_GITHUB_REPOSITORY:
						{
							debug.msg(debug.UI(), "current state is ADD_GITHUB_REPOSITORY")
							if r.TextInput.Value() != "" {
								r.TextInput.Reset()
							}

							r.state = DEFAULT
						}
					}

				}
			}
		}
	}

	r.TextInput, cmd = r.TextInput.Update(msg)

	debug.msg(debug.UI(), r.TextInput.Value())

	return r, cmd
}

const HELP = "ctrl+q quit, ctrl+g update github token, ctrl+r add github repository"

func (r *SettingsScreenModel) View() string {
	if r.state == ADD_GITHUB_TOKEN || r.state == ADD_GITHUB_REPOSITORY {
		return fmt.Sprintf(
			"Paste your GitHub token here:\n\n%s\n\n%s",
			r.TextInput.View(),
			"(esc to quit)") + "\n"
	}

	return lipgloss.JoinVertical(lipgloss.Left, "renders settings screen", HELP)
}
