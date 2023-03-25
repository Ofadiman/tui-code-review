package main

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func NewRouter(settingsScreen *SettingsScreen, pullRequestsScreen *PullRequestsScreen, globalState *GlobalState, settings *Settings, logger *Logger) *Router {
	return &Router{
		currentScreen:      "settings",
		SettingsScreen:     settingsScreen,
		PullRequestsScreen: pullRequestsScreen,
		GlobalState:        globalState,
		Settings:           settings,
		Logger:             logger,
	}

}

type Router struct {
	currentScreen string
	*SettingsScreen
	*PullRequestsScreen
	*GlobalState
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
	if r.currentScreen == "settings" {
		_, cmd = r.SettingsScreen.Update(msg)
	} else {
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
			r.GlobalState.WindowHeight = msg.Height
			r.GlobalState.WindowWidth = msg.Width
		}
	}

	return r, cmd
}

func (r *Router) View() string {
	if r.currentScreen == "settings" {
		return r.SettingsScreen.View()
	} else {
		return r.PullRequestsScreen.View()
	}
}

func main() {
	logger := NewLogger()

	settingsInstance := NewSettings(logger)
	settingsInstance.Load()

	gitHubApi := NewGithubApi(settingsInstance.GithubToken)

	globalState := NewGlobalState()

	settingsScreen := NewSettingsScreen(globalState, settingsInstance, logger, gitHubApi)

	pullRequestsScreen := NewPullRequestsScreen(globalState, settingsInstance, logger, gitHubApi)

	router := NewRouter(settingsScreen, pullRequestsScreen, globalState, settingsInstance, logger)

	program := tea.NewProgram(router, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
