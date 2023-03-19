package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"time"
)

var now int64 = time.Now().Unix()

func debug(label, msg string) {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", os.ModePerm); err != nil {
			panic(err)
		}
	}

	file, err := tea.LogToFile(fmt.Sprintf("logs/%v.log", now), "")
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
