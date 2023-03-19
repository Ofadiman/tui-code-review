package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Khan/genqlient/graphql"
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

func init() {
}

type RouterModel struct {
	activeModel        string
	settingsScreen     SettingsScreenModel
	pullRequestsScreen PullRequestsScreenModel
}

func (r RouterModel) Init() tea.Cmd {
	return nil
}

func (r RouterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if r.activeModel == "settings" {
		_, cmd = r.settingsScreen.Update(msg)
	} else {
		_, cmd = r.pullRequestsScreen.Update(msg)
	}

	switch msg.(type) {
	case tea.KeyMsg:
		switch msg.(tea.KeyMsg).Type {
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
	}

	return r, cmd
}

func (r RouterModel) View() string {
	if r.activeModel == "settings" {
		return r.settingsScreen.View()
	} else {
		return r.pullRequestsScreen.View()
	}
}

type Settings struct {
	GithubToken string `json:"github_token,omitempty"`
}

const CONFIG_FILE_NAME = ".tui-code-review.json"

func main() {
	home, _ := os.UserHomeDir()
	configFilePath := home + "/" + CONFIG_FILE_NAME
	if _, err := os.Stat(configFilePath); err != nil {
		if os.IsNotExist(err) {
			if err = os.WriteFile(configFilePath, []byte("{}\n"), 0644); err != nil {
				debug.msg(debug.FileSystem(), "could not write configuration file")
				debug.msg(debug.Error(), err)
				panic(err)
			}
		} else {
			debug.msg(debug.FileSystem(), "could not stat configuration file")
			debug.msg(debug.Error(), err)
			panic(err)
		}
	}

	file, err := os.ReadFile(configFilePath)
	if err != nil {
		debug.msg(debug.FileSystem(), "could not read configuration file")
		debug.msg(debug.Error(), err)
		panic(err)
	}

	var settings Settings
	err = json.Unmarshal(file, &settings)
	if err != nil {
		debug.msg(debug.FileSystem(), "could not unmarshal configuration file")
		debug.msg(debug.Error(), err.Error())
		panic(err)
	}
	debug.msg(debug.FileSystem(), settings)

	var activeModel string
	if settings.GithubToken == "" {
		activeModel = "settings"
	} else {
		activeModel = "pull_requests"
	}

	httpClient := http.Client{
		Transport: &authedTransport{
			key:     settings.GithubToken,
			wrapped: http.DefaultTransport,
		},
	}
	graphqlClient := graphql.NewClient("https://api.github.com/graphql", &httpClient)

	var response *getRepositoryInfoResponse
	response, err = getRepositoryInfo(context.Background(), graphqlClient)
	if err != nil {
		debug.msg(debug.GraphQL(), "could not fetch data from github")
		debug.msg(debug.Error(), err.Error())
		panic(err)
	}
	debug.msg(debug.GraphQL(), fmt.Sprintf("%#v", response))

	model := RouterModel{
		activeModel:        activeModel,
		settingsScreen:     SettingsScreenModel{},
		pullRequestsScreen: PullRequestsScreenModel{},
	}
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
