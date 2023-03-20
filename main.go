package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
			debug.msg(debug.UI(), strconv.Itoa(msg.Width))

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
	*Settings
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

func (r *RouterModel) WithSettings(settings *Settings) *RouterModel {
	r.Settings = settings

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
			debug.msg(debug.UI(), fmt.Sprintf("window width is set to %v\n", strconv.Itoa(msg.Width)))
			debug.msg(debug.UI(), fmt.Sprintf("window height is set to %v\n", strconv.Itoa(msg.Height)))
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

type Settings struct {
	GithubToken    string   `json:"github_token,omitempty"`
	Repositories   []string `json:"repositories,omitempty"`
	ConfigFilePath string
}

func NewSettings() *Settings {
	home, _ := os.UserHomeDir()

	return &Settings{
		ConfigFilePath: home + "/" + ".tui-code-review.json",
	}
}

func (r *Settings) Load() {
	_, err := os.Stat(r.ConfigFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			r.Save()
		} else {
			debug.msg(debug.FileSystem(), "could not stat configuration file")
			debug.msg(debug.Error(), err.Error())
			panic(err)
		}
	}

	bytes, err := os.ReadFile(r.ConfigFilePath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, r)
	if err != nil {
		debug.msg(debug.FileSystem(), "could not unmarshal configuration file")
		debug.msg(debug.Error(), err.Error())
		panic(err)
	}

	debug.msg(debug.FileSystem(), r)
}

func (r *Settings) Save() {
	bytes, err := json.Marshal(r)
	if err != nil {
		debug.msg(debug.FileSystem(), "could not stat configuration file")
		debug.msg(debug.Error(), err)
		panic(err)
	}

	err = os.WriteFile(r.ConfigFilePath, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	settings := NewSettings()
	settings.Load()

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
		WithSettings(settings).
		WithGlobalState(globalState)
	pullRequestsScreenModel := NewPullRequestsScreenModel().
		WithSettings(settings).
		WithGlobalState(globalState)
	routerModel := NewRouterModel().
		WithSettings(settings).
		WithGlobalState(globalState).
		WithSettingsScreenModel(settingsScreen).
		WithPullRequestsScreenModel(pullRequestsScreenModel)
	routerModel.activeModel = "settings"

	program := tea.NewProgram(routerModel, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
