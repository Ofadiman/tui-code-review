package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ofadiman/tui-code-review/log"
	"github.com/ofadiman/tui-code-review/settings"
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
	*settings.Settings
	*log.Logger
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

func (r *SettingsScreenModel) WithSettings(settings *settings.Settings) *SettingsScreenModel {
	r.Settings = settings

	return r
}

func (r *SettingsScreenModel) WithLogger(logger *log.Logger) *SettingsScreenModel {
	r.Logger = logger

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
					r.Logger.KeyPress("escape")

					if r.state == ADD_GITHUB_TOKEN || r.state == ADD_GITHUB_REPOSITORY {
						r.state = DEFAULT
						r.TextInput.Reset()
					}
				}
			case tea.KeyCtrlU:
				{
					r.Logger.KeyPress("ctrl + u")

					r.state = ADD_GITHUB_TOKEN
				}
			case tea.KeyCtrlR:
				{
					r.Logger.KeyPress("ctrl + r")

					r.state = ADD_GITHUB_REPOSITORY
				}
			case tea.KeyEnter:
				{
					r.Logger.KeyPress("enter")

					switch r.state {
					case ADD_GITHUB_TOKEN:
						{
							r.Logger.Info(fmt.Sprintf("current input value %v", r.TextInput.Value()))

							r.Settings.UpdateGitHubToken(r.TextInput.Value())

							if r.TextInput.Value() != "" {
								r.TextInput.Reset()
							}

							r.state = DEFAULT
						}
					case ADD_GITHUB_REPOSITORY:
						{
							r.Logger.Info(fmt.Sprintf("current input value %v", r.TextInput.Value()))

							r.Settings.AddRepositoryUrl(r.TextInput.Value())

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

	return r, cmd
}

const HELP = "ctrl+q quit, ctrl+u update github token, ctrl+r add github repository"

func (r *SettingsScreenModel) View() string {
	if r.state == ADD_GITHUB_TOKEN || r.state == ADD_GITHUB_REPOSITORY {
		return fmt.Sprintf(
			"Paste your GitHub token here:\n\n%s\n\n%s",
			r.TextInput.View(),
			"(esc to quit)") + "\n"
	}

	return lipgloss.JoinVertical(lipgloss.Left, "renders settings screen", HELP)
}
