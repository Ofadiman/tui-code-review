package main

import (
	"fmt"
	"github.com/ofadiman/tui-code-review/log"
	"github.com/ofadiman/tui-code-review/settings"
	"net/http"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type authedTransport struct {
	key     string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "bearer "+t.key)
	return t.wrapped.RoundTrip(req)
}

type PullRequest struct {
	id    string
	title string
}

func (p PullRequest) FilterValue() string {
	return p.title
}

func (p PullRequest) Title() string {
	return p.title
}

func (p PullRequest) Description() string {
	return "desc"
}

type Column string

const (
	Waiting Column = "waiting"
	Checked Column = "checked"
)

type Model struct {
	activeColumn Column
	waiting      list.Model
	checked      list.Model
}

func (m *Model) ToggleActiveColumn() {
	if m.activeColumn == Checked {
		m.activeColumn = Waiting
	} else {
		m.activeColumn = Checked
	}
}

func (m *Model) GetActiveList() *list.Model {
	if m.activeColumn == Waiting {
		return &m.waiting
	} else {
		return &m.checked
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

var defaultColumnStyle = lipgloss.
	NewStyle().
	Border(lipgloss.HiddenBorder())

var activeColumnStyle = lipgloss.
	NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("63"))

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.String() {

			case "ctrl+c", "q":
				{
					return m, tea.Quit
				}

			case "up", "k":
				{
					if m.activeColumn == Waiting {
						updatedList, cmd := m.waiting.Update(msg)
						m.waiting = updatedList
						return m, cmd
					} else {
						updatedList, cmd := m.checked.Update(msg)
						m.checked = updatedList
						return m, cmd
					}
				}

			case "down", "j":
				{
					if m.activeColumn == Waiting {
						updatedList, cmd := m.waiting.Update(msg)
						m.waiting = updatedList
						return m, cmd
					} else {
						updatedList, cmd := m.checked.Update(msg)
						m.checked = updatedList
						return m, cmd
					}
				}

			case "left", "h":
				{
					m.ToggleActiveColumn()
					return m, nil
				}

			case "right", "l":
				{
					m.ToggleActiveColumn()
					return m, nil
				}

			case "enter", " ":
				{
				}
			}
		}
	case tea.WindowSizeMsg:
		{
			horizontalFrameSize, verticalFrameSize := defaultColumnStyle.GetFrameSize()
			defaultColumnStyle.Width(msg.Width/2 - horizontalFrameSize)
			activeColumnStyle.Width(msg.Width/2 - horizontalFrameSize)

			m.waiting.SetSize(msg.Width/2-horizontalFrameSize, msg.Height-verticalFrameSize)
			m.waiting.SetShowHelp(false)

			m.checked.SetSize(msg.Width/2-horizontalFrameSize, msg.Height-verticalFrameSize)
			m.checked.SetShowHelp(false)

		}
	}

	var cmd tea.Cmd
	m.waiting, cmd = m.waiting.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	if m.activeColumn == Waiting {
		return lipgloss.JoinHorizontal(lipgloss.Left, activeColumnStyle.Render(m.waiting.View()), defaultColumnStyle.Render(m.checked.View()))
	}

	if m.activeColumn == Checked {
		return lipgloss.JoinHorizontal(lipgloss.Left, defaultColumnStyle.Render(m.waiting.View()), activeColumnStyle.Render(m.checked.View()))
	}

	panic("invalid activeColumn value")
}

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

	//httpClient := http.Client{
	//	Transport: &authedTransport{
	//		key:     settings.GithubToken,
	//		wrapped: http.DefaultTransport,
	//	},
	//}
	//graphqlClient := graphql.NewClient("https://api.github.com/graphql", &httpClient)

	//var response *getRepositoryInfoResponse
	//response, err = getRepositoryInfo(context.Background(), graphqlClient)
	//if err != nil {
	//	debug.msg(debug.GraphQL(), "could not fetch data from github")
	//	debug.msg(debug.Error(), err.Error())
	//	panic(err)
	//}
	//debug.msg(debug.GraphQL(), fmt.Sprintf("%#v", response))

	globalState := NewGlobalState()
	settingsScreen := NewSettingsScreenModel().
		WithSettings(settings_).
		WithGlobalState(globalState).
		WithLogger(logger)
	pullRequestsScreenModel := NewPullRequestsScreenModel().
		WithSettings(settings_).
		WithGlobalState(globalState).
		WithLogger(logger)
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
