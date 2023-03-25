package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/ofadiman/tui-code-review/log"
	"github.com/ofadiman/tui-code-review/settings"
	"os/exec"
	"runtime"
)

type state string

const (
	ADD_GITHUB_TOKEN      state = "ADD_GITHUB_TOKEN"
	ADD_GITHUB_REPOSITORY state = "ADD_GITHUB_REPOSITORY"
	DEFAULT               state = "DEFAULT"
)

type SettingsScreen struct {
	TextInput               textinput.Model
	state                   state
	SelectedRepositoryIndex int
	*GlobalState
	*settings.Settings
	*log.Logger
	*GithubApi
}

func NewSettingsScreen(globalState *GlobalState, settings *settings.Settings, logger *log.Logger, gitHubApi *GithubApi) *SettingsScreen {
	textInput := textinput.New()
	textInput.Placeholder = "Type something..."
	textInput.CharLimit = 200
	textInput.Focus()
	textInput.Width = 50

	return &SettingsScreen{
		TextInput:               textInput,
		state:                   DEFAULT,
		SelectedRepositoryIndex: 0,
		GlobalState:             globalState,
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
			case "j":
				{
					if r.SelectedRepositoryIndex == len(r.Repositories)-1 {
						r.SelectedRepositoryIndex = 0
					} else {
						r.SelectedRepositoryIndex = r.SelectedRepositoryIndex + 1
					}
				}
			case "k":
				{
					if r.SelectedRepositoryIndex == 0 {
						r.SelectedRepositoryIndex = len(r.Repositories) - 1
					} else {
						r.SelectedRepositoryIndex = r.SelectedRepositoryIndex - 1
					}
				}
			case "esc":
				{
					if r.state == ADD_GITHUB_TOKEN || r.state == ADD_GITHUB_REPOSITORY {
						r.state = DEFAULT
						r.TextInput.Reset()
					}
				}
			case "ctrl+u":
				{
					r.state = ADD_GITHUB_TOKEN
				}
			case "ctrl+r":
				{
					r.state = ADD_GITHUB_REPOSITORY
				}
			case "delete":
				{
					r.Settings.DeleteRepositoryUrl(r.Settings.Repositories[r.SelectedRepositoryIndex])
				}
			case "enter":
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
					case ADD_GITHUB_TOKEN:
						{
							r.Logger.Info(fmt.Sprintf("current input value %v", r.TextInput.Value()))

							r.Settings.UpdateGitHubToken(r.TextInput.Value())
							r.GithubApi.UpdateClient(r.TextInput.Value())

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

	if r.state == ADD_GITHUB_TOKEN || r.state == ADD_GITHUB_REPOSITORY {
		r.TextInput, cmd = r.TextInput.Update(msg)
	}

	return r, cmd
}

var greyText = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
var whiteText = lipgloss.NewStyle().Foreground(lipgloss.Color("231"))

const DELIMITER = " • "

var HELP_QUIT = lipgloss.JoinHorizontal(lipgloss.Left, whiteText.Render("Ctrl+Q"), " ", greyText.Render("Quit program"))
var HELP_UPDATE_GITHUB_TOKEN = lipgloss.JoinHorizontal(lipgloss.Left, whiteText.Render("Ctrl+U"), " ", greyText.Render("Update GitHub token"))
var HELP_ADD_GITHUB_REPOSITORY = lipgloss.JoinHorizontal(lipgloss.Left, whiteText.Render("Ctrl+R"), " ", greyText.Render("Add GitHub repository"))
var HELP_DELETE_GITHUB_REPOSITORY = lipgloss.JoinHorizontal(lipgloss.Left, whiteText.Render("Delete"), " ", greyText.Render("Delete selected GitHub repository"))
var HELP_OPEN_GITHUB_REPOSITORY = lipgloss.JoinHorizontal(lipgloss.Left, whiteText.Render("Enter"), " ", greyText.Render("Open selected GitHub repository"))
var HELP_J = lipgloss.JoinHorizontal(lipgloss.Left, whiteText.Render("J"), " ", greyText.Render("Move down"))
var HELP_K = lipgloss.JoinHorizontal(lipgloss.Left, whiteText.Render("K"), " ", greyText.Render("Move up"))
var HELP = lipgloss.JoinHorizontal(lipgloss.Left, HELP_QUIT, DELIMITER, HELP_UPDATE_GITHUB_TOKEN, DELIMITER, HELP_ADD_GITHUB_REPOSITORY, DELIMITER, HELP_DELETE_GITHUB_REPOSITORY, DELIMITER, HELP_OPEN_GITHUB_REPOSITORY, DELIMITER, HELP_J, DELIMITER, HELP_K)

func (r *SettingsScreen) View() string {
	c := lipgloss.NewStyle().PaddingTop(1).PaddingBottom(1).PaddingLeft(2).PaddingRight(2)

	if r.state == ADD_GITHUB_TOKEN {
		return c.Render(fmt.Sprintf(
			"Paste your GitHub token here:\n\n%s\n\n%s",
			r.TextInput.View(),
			"(esc to quit)") + "\n")
	}

	if r.state == ADD_GITHUB_REPOSITORY {
		return c.Render(fmt.Sprintf(
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

	wrapper := wordwrap.NewWriter(r.GlobalState.WindowWidth - roundedBorder.GetLeftSize() - roundedBorder.GetRightSize())
	_, err := wrapper.Write([]byte(HELP))
	if err != nil {
		r.Logger.Error(err)
	}

	return c.Render(lipgloss.JoinVertical(lipgloss.Left, "Settings\n", repositories, wrapper.String()))
}
