package cli

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	ProjectNamePage int = iota
	ConfigTypePage
	RouterTypePage
)

type View interface {
	Upgrade(msg tea.Msg, m *Model) tea.Cmd
}

type Model struct {
	Views            []View
	CurrentViewIndex int
	CreateProject    bool
	Progress         progress.Model
	ProgressPercent  float64
	Quitting         bool
	ProgressChannel  chan tea.Msg
}

func Start() {
	initialModel := Model{
		Views: []View{
			NewProjectNameView(),
		},
		CurrentViewIndex: 0,
		Progress:         progress.New(),
		ProgressChannel:  make(chan tea.Msg),
	}
}
