package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SettingsScreenModel struct {
	WindowWidth  int
	WindowHeight int
}

func (r *SettingsScreenModel) Init() tea.Cmd {
	return nil
}

func (r *SettingsScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.String() {
			case "p":
				{
					debug.msg(debug.KeyPressed(), "p")
					return r, nil
				}
			}
		}
	}

	return r, nil
}

func (r *SettingsScreenModel) View() string {
	return "renders settings screen"
}
