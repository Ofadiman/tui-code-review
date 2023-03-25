package main

import (
	"fmt"
	"github.com/ofadiman/tui-code-review/globals"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func NewRouter(settingsScreen *SettingsScreen, pullRequestsScreen *PullRequestsScreen, globals *globals.Globals) *Router {
	return &Router{
		currentScreen:      "settings",
		SettingsScreen:     settingsScreen,
		PullRequestsScreen: pullRequestsScreen,
		Globals:            globals,
	}

}

type Router struct {
	currentScreen string
	*SettingsScreen
	*PullRequestsScreen
	*globals.Globals
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
			r.Globals.Info(fmt.Sprintf("window width is set to %v", strconv.Itoa(msg.Width)))
			r.Globals.Info(fmt.Sprintf("window height is set to %v", strconv.Itoa(msg.Height)))
			r.Globals.Height = msg.Height
			r.Globals.Width = msg.Width
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
	glob := globals.NewGlobals()
	settingsScreen := NewSettingsScreen(glob)
	pullRequestsScreen := NewPullRequestsScreen(glob)
	router := NewRouter(settingsScreen, pullRequestsScreen, glob)

	program := tea.NewProgram(router, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
