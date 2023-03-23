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
	TextInput               textinput.Model
	state                   state
	SelectedRepositoryIndex int
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

			switch msg.String() {
			case "j":
				{
					r.Logger.KeyPress("j")

					if r.SelectedRepositoryIndex == len(r.Repositories)-1 {
						r.SelectedRepositoryIndex = 0
					} else {
						r.SelectedRepositoryIndex = r.SelectedRepositoryIndex + 1
					}
				}
			case "k":
				{
					r.Logger.KeyPress("k")

					if r.SelectedRepositoryIndex == 0 {
						r.SelectedRepositoryIndex = len(r.Repositories) - 1
					} else {
						r.SelectedRepositoryIndex = r.SelectedRepositoryIndex - 1
					}
				}
			case "escape":
				{
					r.Logger.KeyPress("escape")

					if r.state == ADD_GITHUB_TOKEN || r.state == ADD_GITHUB_REPOSITORY {
						r.state = DEFAULT
						r.TextInput.Reset()
					}
				}
			case "ctrl+u":
				{
					r.Logger.KeyPress("ctrl + u")

					r.state = ADD_GITHUB_TOKEN
				}
			case "ctrl+r":
				{
					r.Logger.KeyPress("ctrl + r")

					r.state = ADD_GITHUB_REPOSITORY
				}
			case "delete":
				{
					r.Settings.DeleteRepositoryUrl(r.Settings.Repositories[r.SelectedRepositoryIndex])
				}
			case "enter":
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

const HELP = "ctrl+q quit, ctrl+u update github token, ctrl+r add github repository delete delete selected repository"

func (r *SettingsScreenModel) View() string {
	if r.state == ADD_GITHUB_TOKEN || r.state == ADD_GITHUB_REPOSITORY {
		return fmt.Sprintf(
			"Paste your GitHub token here:\n\n%s\n\n%s",
			r.TextInput.View(),
			"(esc to quit)") + "\n"
	}

	repositories := ""

	s := lipgloss.NewStyle()
	x := lipgloss.NewStyle().Inherit(s).Underline(true)
	for index, url := range r.Settings.Repositories {
		if index == r.SelectedRepositoryIndex {
			repositories += x.Render(url)
			repositories += "\n"
		} else {
			repositories += s.Render(url)
			repositories += "\n"
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, "renders settings screen\n", repositories, HELP)
}
