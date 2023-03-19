package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"time"
)

var debug = Debug{
	now: time.Now().Unix(),
}

type Debug struct {
	now int64
}

func (r Debug) msg(label, msg string) {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", os.ModePerm); err != nil {
			panic(err)
		}
	}

	file, err := tea.LogToFile(fmt.Sprintf("logs/%v.log", r.now), "")
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString(label + ": " + msg + "\n")
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}
}

func (r Debug) GraphQL() string {
	return "graphql"
}

func (r Debug) Error() string {
	return "error"
}

func (r Debug) Environment() string {
	return "environment"
}

func (r Debug) UI() string {
	return "ui"
}
