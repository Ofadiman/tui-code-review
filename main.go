package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kelseyhightower/envconfig"

	"github.com/joho/godotenv"

	"github.com/Khan/genqlient/graphql"
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

type Env struct {
	GitHubToken string `envconfig:"GITHUB_TOKEN" required:"true"`
}

var env Env

func init() {
	err := godotenv.Load()
	if err != nil {
		debug.msg(debug.GraphQL(), err.Error())
	}

	err = envconfig.Process("", &env)
	if err != nil {
		debug.msg(debug.Error(), err.Error())
	}

	debug.msg(debug.Environment(), fmt.Sprintf("%#v", env))

	httpClient := http.Client{
		Transport: &authedTransport{
			key:     env.GitHubToken,
			wrapped: http.DefaultTransport,
		},
	}
	graphqlClient := graphql.NewClient("https://api.github.com/graphql", &httpClient)

	var response *getRepositoryInfoResponse
	response, err = getRepositoryInfo(context.Background(), graphqlClient)
	if err != nil {
		panic(err)
	}
	debug.msg(debug.GraphQL(), fmt.Sprintf("%#v", response))
}

func main() {
	model := &Model{
		activeColumn: Waiting,
		waiting: list.New([]list.Item{
			&PullRequest{
				id:    "db9eede3-0c80-456a-b323-e8c302506950",
				title: "implement feature 1",
			},
			&PullRequest{
				id:    "4c8e08dc-92da-4397-839a-2cad98706d3a",
				title: "implement feature 2",
			},
			&PullRequest{
				id:    "c0a9ebb2-9ff9-449e-86c0-c7313acc2591",
				title: "implement feature 4",
			},
		}, list.NewDefaultDelegate(), 0, 0),
		checked: list.New([]list.Item{
			&PullRequest{
				id:    "42660d9c-cabd-4372-968d-e68087e42c65",
				title: "implement feature 3",
			},
			&PullRequest{
				id:    "14d32ef4-043a-41ad-bcad-28b695248b3a",
				title: "implement feature 5",
			},
		}, list.NewDefaultDelegate(), 0, 0),
	}

	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
