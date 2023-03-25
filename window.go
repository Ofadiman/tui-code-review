package main

type Window struct {
	Width  int
	Height int
}

func NewWindow() *Window {
	return &Window{}
}
