package main

type GlobalState struct {
	WindowWidth  int
	WindowHeight int
}

func NewGlobalState() *GlobalState {
	return &GlobalState{}
}
