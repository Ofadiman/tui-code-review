package globals

import (
	"fmt"
	"os"
	"time"
)

type logger struct {
	file string
}

func newLogger() *logger {
	logger := &logger{
		file: fmt.Sprintf("logs/%v.log", time.Now().Unix()),
	}

	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err = os.Mkdir("logs", os.ModePerm); err != nil {
			panic(err)
		}
	}

	return logger
}

func (r *logger) append(msg string) {
	file, _ := os.OpenFile(r.file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	file.WriteString(fmt.Sprintf("%v\n", msg))
}

func (r *logger) Error(error error) {
	r.append(fmt.Sprintf("error: %v", error))
}

func (r *logger) Info(msg string) {
	r.append(fmt.Sprintf("info: %v", msg))
}

func (r *logger) Struct(msg any) {
	r.append(fmt.Sprintf("struct: %#v", msg))
}

func (r *logger) KeyPress(msg string) {
	r.append(fmt.Sprintf("key_press: %v", msg))
}
