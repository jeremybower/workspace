package internal

type State struct {
	path string
}

func NewState(path string) *State {
	return &State{path}
}
