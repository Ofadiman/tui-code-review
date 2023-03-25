package globals

type window struct {
	Width  int
	Height int
}

func newWindow() *window {
	return &window{}
}
