package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"math"
	"os/exec"
	"runtime"
)

const (
	UPDATE_GITHUB_TOKEN       string = "UPDATE_GITHUB_TOKEN"
	ADD_GITHUB_REPOSITORY_URL string = "ADD_GITHUB_REPOSITORY_URL"
	DEFAULT                   string = "DEFAULT"
)

var SETTINGS_HELP = []Help{helpUp, helpDown, helpQuit, helpAddGitHubRepositoryUrl, helpDeleteGitHubRepositoryUrl, helpOpenGitHubRepositoryUrl, helpUpdateGithubToken, helpSwitchToPullRequestsScreen}

type SettingsScreen struct {
	TextInput               textinput.Model
	state                   string
	SelectedRepositoryIndex int
	*Window
	*Settings
	*Logger
	*GithubApi
}

func NewSettingsScreen(globalState *Window, settings *Settings, logger *Logger, gitHubApi *GithubApi) *SettingsScreen {
	textInput := textinput.New()
	textInput.Placeholder = "Type something..."
	textInput.CharLimit = 200
	textInput.Focus()
	textInput.Width = 50

	return &SettingsScreen{
		TextInput:               textInput,
		state:                   DEFAULT,
		SelectedRepositoryIndex: 0,
		Window:                  globalState,
		Settings:                settings,
		Logger:                  logger,
		GithubApi:               gitHubApi,
	}
}

func (r *SettingsScreen) Init() tea.Cmd {
	return nil
}

func (r *SettingsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			r.Logger.KeyPress(msg.String())

			switch msg.String() {
			case helpDown.Shortcut:
				{
					if r.SelectedRepositoryIndex == len(r.Repositories)-1 {
						r.SelectedRepositoryIndex = 0
					} else {
						r.SelectedRepositoryIndex = r.SelectedRepositoryIndex + 1
					}
				}
			case helpUp.Shortcut:
				{
					if r.SelectedRepositoryIndex == 0 {
						r.SelectedRepositoryIndex = len(r.Repositories) - 1
					} else {
						r.SelectedRepositoryIndex = r.SelectedRepositoryIndex - 1
					}
				}
			case helpEscape.Shortcut:
				{
					if r.state == UPDATE_GITHUB_TOKEN || r.state == ADD_GITHUB_REPOSITORY_URL {
						r.state = DEFAULT
						r.TextInput.Reset()
					}
				}
			case helpUpdateGithubToken.Shortcut:
				{
					r.state = UPDATE_GITHUB_TOKEN
				}
			case helpAddGitHubRepositoryUrl.Shortcut:
				{
					r.state = ADD_GITHUB_REPOSITORY_URL
				}
			case helpDeleteGitHubRepositoryUrl.Shortcut:
				{
					r.Settings.DeleteRepositoryUrl(r.Settings.Repositories[r.SelectedRepositoryIndex])
					r.SelectedRepositoryIndex = int(math.Max(float64(r.SelectedRepositoryIndex-1), float64(0)))
				}
			case helpOpenGitHubRepositoryUrl.Shortcut:
				{
					switch r.state {
					case DEFAULT:
						{
							selectedRepository := r.Settings.Repositories[r.SelectedRepositoryIndex]

							r.Logger.Info(fmt.Sprintf("opening a default browser on %v page", selectedRepository))

							var err error
							switch runtime.GOOS {
							case "linux":
								{
									err = exec.Command("xdg-open", selectedRepository).Start()
								}
							case "windows":
								{
									err = exec.Command("rundll32", "url.dll,FileProtocolHandler", selectedRepository).Start()
								}
							case "darwin":
								{
									err = exec.Command("open", selectedRepository).Start()
								}
							default:
								err = fmt.Errorf("unsupported platform %v", runtime.GOOS)
							}

							if err != nil {
								panic(err)
							}
						}
					case UPDATE_GITHUB_TOKEN:
						{
							r.Logger.Info(fmt.Sprintf("current input value %v", r.TextInput.Value()))

							r.Settings.UpdateGitHubToken(r.TextInput.Value())
							r.GithubApi.UpdateClient(r.TextInput.Value())

							if r.TextInput.Value() != "" {
								r.TextInput.Reset()
							}

							r.state = DEFAULT
						}
					case ADD_GITHUB_REPOSITORY_URL:
						{
							r.Logger.Info(fmt.Sprintf("current input value %v", r.TextInput.Value()))

							r.Settings.AddRepositoryUrl(r.TextInput.Value())

							if r.TextInput.Value() != "" {
								r.SelectedRepositoryIndex = len(r.Repositories) - 1
								r.TextInput.Reset()
							}

							r.state = DEFAULT
						}
					}

				}
			}
		}
	}

	if r.state == UPDATE_GITHUB_TOKEN || r.state == ADD_GITHUB_REPOSITORY_URL {
		r.TextInput, cmd = r.TextInput.Update(msg)
	}

	return r, cmd
}

func (r *SettingsScreen) View() string {

	if r.state == UPDATE_GITHUB_TOKEN {
		return StyledMain.Render(fmt.Sprintf(
			"Paste your GitHub token here:\n\n%s\n\n%s",
			r.TextInput.View(),
			"(esc to quit)") + "\n")
	}

	if r.state == ADD_GITHUB_REPOSITORY_URL {
		return StyledMain.Render(fmt.Sprintf(
			"Paste your repository URL here:\n\n%s\n\n%s",
			r.TextInput.View(),
			"(esc to quit)") + "\n")
	}

	repositories := ""

	s := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	x := lipgloss.NewStyle().Underline(true)
	for index, url := range r.Settings.Repositories {
		if index == r.SelectedRepositoryIndex {
			repositories += x.Render(url)
			repositories += "\n"
		} else {
			repositories += s.Render(url)
			repositories += "\n"
		}
	}

	helpString := ""
	for _, help := range SETTINGS_HELP {
		helpString += lipgloss.JoinHorizontal(lipgloss.Left, StyledHelpShortcut.Render(help.Display), " ", StyledHelpDescription.Render(help.Description), "   ")
	}

	wrapper := wordwrap.NewWriter(r.Window.Width - StyledMain.GetHorizontalPadding())
	wrapper.Breakpoints = []rune{' '}
	_, err := wrapper.Write([]byte(helpString))
	if err != nil {
		r.Logger.Error(err)
	}

	return StyledMain.Render(lipgloss.JoinVertical(lipgloss.Left, StyledHeader.Render("Settings"), repositories, wrapper.String()))
}
