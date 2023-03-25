package main

import (
	"fmt"
	"github.com/ofadiman/tui-code-review/log"
	"github.com/ofadiman/tui-code-review/settings"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type GlobalState struct {
	WindowWidth  int
	WindowHeight int
}

func NewGlobalState() *GlobalState {
	return &GlobalState{}
}

type RouterModel struct {
	activeModel string
	*SettingsScreenModel
	*PullRequestsScreenModel
	*GlobalState
	*settings.Settings
	*log.Logger
}

func NewRouterModel() *RouterModel {
	return &RouterModel{
		activeModel:             "",
		SettingsScreenModel:     nil,
		PullRequestsScreenModel: nil,
		GlobalState:             nil,
	}
}

func (r *RouterModel) WithGlobalState(globalState *GlobalState) *RouterModel {
	r.GlobalState = globalState

	return r
}

func (r *RouterModel) WithSettingsScreenModel(model *SettingsScreenModel) *RouterModel {
	r.SettingsScreenModel = model

	return r
}

func (r *RouterModel) WithPullRequestsScreenModel(model *PullRequestsScreenModel) *RouterModel {
	r.PullRequestsScreenModel = model

	return r
}

func (r *RouterModel) WithSettings(settings *settings.Settings) *RouterModel {
	r.Settings = settings

	return r
}

func (r *RouterModel) WithLogger(logger *log.Logger) *RouterModel {
	r.Logger = logger

	return r
}

func (r *RouterModel) Init() tea.Cmd {
	r.SettingsScreenModel.Init()
	r.PullRequestsScreenModel.Init()
	return nil
}

func (r *RouterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if r.activeModel == "settings" {
		_, cmd = r.SettingsScreenModel.Update(msg)
	} else {
		_, cmd = r.PullRequestsScreenModel.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			{
				r.activeModel = "settings"
				return r, cmd
			}
		case tea.KeyCtrlP:
			{
				r.activeModel = "pull_requests"
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

func (r *RouterModel) View() string {
	if r.activeModel == "settings" {
		return r.SettingsScreenModel.View()
	} else {
		return r.PullRequestsScreenModel.View()
	}
}

func main() {
	logger := log.NewLogger()

	settings_ := settings.NewSettings().WithLogger(logger)
	settings_.Load()

	gitHubGraphqlApi := NewGithubApi(settings_.GithubToken)

	globalState := NewGlobalState()

	settingsScreen := NewSettingsScreenModel().
		WithSettings(settings_).
		WithGlobalState(globalState).
		WithLogger(logger).
		WithGitHubGraphqlApi(gitHubGraphqlApi)

	pullRequestsScreenModel := NewPullRequestsScreenModel().
		WithSettings(settings_).
		WithGlobalState(globalState).
		WithLogger(logger).
		WithGitHubGraphqlApi(gitHubGraphqlApi)

	routerModel := NewRouterModel().
		WithSettings(settings_).
		WithGlobalState(globalState).
		WithSettingsScreenModel(settingsScreen).
		WithPullRequestsScreenModel(pullRequestsScreenModel).
		WithLogger(logger)

	routerModel.activeModel = "settings"

	program := tea.NewProgram(routerModel, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
