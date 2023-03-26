package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

const SCREEN_SETTINGS = "settings"
const SCREEN_PULL_REQUESTS = "pull_requests"

func NewRouter(settingsScreen *SettingsScreen, pullRequestsScreen *PullRequestsScreen, globalState *Window, settings *Settings, logger *Logger) *Router {
	return &Router{
		currentScreen:      SCREEN_PULL_REQUESTS,
		SettingsScreen:     settingsScreen,
		PullRequestsScreen: pullRequestsScreen,
		Window:             globalState,
		Settings:           settings,
		Logger:             logger,
	}

}

type Router struct {
	currentScreen string
	*SettingsScreen
	*PullRequestsScreen
	*Window
	*Settings
	*Logger
}

func (r *Router) Init() tea.Cmd {
	r.SettingsScreen.Init()
	r.PullRequestsScreen.Init()
	return nil
}

func (r *Router) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if r.currentScreen == SCREEN_SETTINGS {
		_, cmd = r.SettingsScreen.Update(msg)
	}

	if r.currentScreen == SCREEN_PULL_REQUESTS {
		_, cmd = r.PullRequestsScreen.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			{
				r.currentScreen = "settings"
				return r, cmd
			}
		case tea.KeyCtrlP:
			{
				r.currentScreen = "pull_requests"
				return r, cmd
			}
		case tea.KeyCtrlQ:
			{
				return r, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		{
			r.Logger.Info(fmt.Sprintf("window width is set to %v", strconv.Itoa(msg.Width)))
			r.Logger.Info(fmt.Sprintf("window height is set to %v", strconv.Itoa(msg.Height)))
			r.Window.Height = msg.Height
			r.Window.Width = msg.Width

			StyledHeader.Width(msg.Width - lipgloss.RoundedBorder().GetLeftSize() - lipgloss.RoundedBorder().GetRightSize() - StyledMain.GetPaddingLeft() - StyledMain.GetPaddingRight())
		}
	}

	return r, cmd
}

func (r *Router) View() string {
	if r.currentScreen == SCREEN_SETTINGS {
		return r.SettingsScreen.View()
	}

	if r.currentScreen == SCREEN_PULL_REQUESTS {
		return r.PullRequestsScreen.View()
	}

	panic(fmt.Sprintf("incorrect screen name %v", r.currentScreen))
}

func main() {
	logger := NewLogger()

	settingsInstance := NewSettings(logger)
	settingsInstance.Load()

	gitHubApi := NewGithubApi(settingsInstance.GithubToken)

	globalState := NewWindow()

	settingsScreen := NewSettingsScreen(globalState, settingsInstance, logger, gitHubApi)

	pullRequestsScreen := NewPullRequestsScreen(globalState, settingsInstance, logger, gitHubApi)

	router := NewRouter(settingsScreen, pullRequestsScreen, globalState, settingsInstance, logger)

	program := tea.NewProgram(router, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
