package main

import (
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func debug(msg string) {
	file, err := tea.LogToFile("debug.log", "")
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString(msg + "\n")
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}
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

type Model struct {
	waiting list.Model
	checked list.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

var columnStyle = lipgloss.
	NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("63"))

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				}

			case "down", "j":
				{
				}

			case "enter", " ":
				{
				}
			}
		}
	case tea.WindowSizeMsg:
		{
			horizontalFrameSize, verticalFrameSize := columnStyle.GetFrameSize()
			columnStyle.Width(msg.Width/2 - horizontalFrameSize)
			debug(strconv.Itoa(msg.Width))

			m.waiting = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
			m.waiting.SetItems([]list.Item{
				&PullRequest{
					id:    "db9eede3-0c80-456a-b323-e8c302506950",
					title: "implement feature 1",
				},
				&PullRequest{
					id:    "4c8e08dc-92da-4397-839a-2cad98706d3a",
					title: "implement feature 2",
				},
			})
			m.waiting.SetSize(msg.Width/2-horizontalFrameSize, msg.Height-verticalFrameSize)
			m.waiting.SetShowHelp(false)

			m.checked = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
			m.checked.SetItems([]list.Item{
				&PullRequest{
					id:    "42660d9c-cabd-4372-968d-e68087e42c65",
					title: "implement feature 3",
				},
			})
			m.checked.SetSize(msg.Width/2-horizontalFrameSize, msg.Height-verticalFrameSize)
			m.checked.SetShowHelp(false)

		}
	}

	var cmd tea.Cmd
	m.waiting, cmd = m.waiting.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Left, columnStyle.Render(m.waiting.View()), columnStyle.Render(m.checked.View()))
}

func main() {
	program := tea.NewProgram(Model{}, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
