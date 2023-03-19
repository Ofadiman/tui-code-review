package main

import tea "github.com/charmbracelet/bubbletea"

type PullRequestsScreenModel struct {
}

func (r *PullRequestsScreenModel) Init() tea.Cmd {
	return nil
}

func (r *PullRequestsScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.String() {
			case "s":
				{
					debug.msg(debug.KeyPressed(), "s")
					return r, nil
				}
			}
		}
	}

	return r, nil
}

func (r *PullRequestsScreenModel) View() string {
	return "renders pull requests screen"
}
