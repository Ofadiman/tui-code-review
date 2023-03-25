package main

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/wordwrap"
	"strings"
)

var roundedBorder = lipgloss.RoundedBorder()
var columnStyle = lipgloss.NewStyle().Border(roundedBorder).BorderForeground(lipgloss.Color("63"))

type PullRequestsScreen struct {
	*GlobalState
	*Settings
	*Logger
	*GithubApi
}

func NewPullRequestsScreen(globalState *GlobalState, settings *Settings, logger *Logger, githubApi *GithubApi) *PullRequestsScreen {
	return &PullRequestsScreen{
		GlobalState: globalState,
		Settings:    settings,
		Logger:      logger,
		GithubApi:   githubApi,
	}
}

func (r *PullRequestsScreen) Init() tea.Cmd {
	if r.Settings.GithubToken == "" {
		return nil
	}

	var response *getRepositoryInfoResponse
	var err error
	response, err = getRepositoryInfo(context.Background(), *r.GithubApi.client, "Ofadiman", "tui-code-review")
	if err != nil {
		r.Logger.Error(err)

		if strings.Contains(err.Error(), "401") {
			r.Settings.UpdateGitHubToken("")
		}
	}

	r.Logger.Struct(response)

	return nil
}

func (r *PullRequestsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.String() {
			case "s":
				{
					r.Logger.KeyPress("s")
					return r, nil
				}
			}
		}
	}

	return r, nil
}

func (r *PullRequestsScreen) View() string {
	columnStyle.Width(r.GlobalState.WindowWidth - roundedBorder.GetLeftSize() - roundedBorder.GetRightSize())
	header := columnStyle.Render("renders pull requests screen")
	lorem := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."

	f := wordwrap.NewWriter(r.GlobalState.WindowWidth - roundedBorder.GetLeftSize() - roundedBorder.GetRightSize())
	f.Breakpoints = []rune{' '}
	_, err := f.Write([]byte(lorem))
	if err != nil {
		r.Logger.Error(err)
	}

	help := []struct {
		shortcut    string
		description string
	}{
		{
			shortcut:    "ctrl + q",
			description: "quit",
		},
		{
			shortcut:    "ctrl + s",
			description: "settings",
		},
		{
			shortcut:    "ctrl + p",
			description: "pull requests",
		},
	}

	description := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	shortcut := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))

	mapped := make([]string, len(help))
	for i, value := range help {
		mapped[i] = padding.String(fmt.Sprintf("%v %v", shortcut.Render(value.shortcut), description.Render(value.description)), 20)
	}

	return lipgloss.JoinVertical(0, header, columnStyle.Render(f.String()), strings.Join(mapped, ""))
}
