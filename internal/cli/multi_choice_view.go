package cli

import (
	"log"
	"math"
	"project_generator/internal/termstyle"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type MultiChoiceView struct {
	Prompt   string
	Options  []string
	Selected int
}

func NewMultiChoiceView(prompt string, options []string) *MultiChoiceView {
	return &MultiChoiceView{
		Prompt:   prompt,
		Options:  options,
		Selected: 0,
	}
}

func (v *MultiChoiceView) View() string {
	log.Println("MultiChoiceView.View()")
	var builder strings.Builder
	builder.WriteString(v.Prompt + "\n\n")
	for index, option := range v.Options {
		checkbox := Checkbox(option, index == v.Selected)
		builder.WriteString(checkbox + "\n")
	}

	instructions := termstyle.Subtle("enter:choose") +
		termstyle.Dot + termstyle.Subtle("esc or ctrl-c: quit")
	builder.WriteString("\n" + instructions)

	return builder.String()
}

func (v *MultiChoiceView) Update(msg tea.Msg, m *Model) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			v.Selected = int(math.Min(float64(v.Selected+1), float64(len(v.Options)-1)))
		case tea.KeyUp:
			v.Selected = int(math.Max(float64(v.Selected-1), float64(0)))
		case tea.KeyEnter:
			if m.CurrentViewIndex == len(m.Views)-1 {
				m.CreateProject = true
				log.Println("Return createProjectMsg")
				cmd = createProject()
			} else {
				m.CurrentViewIndex++
			}
		}
	}
	return cmd
}

func createProject() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return createProjectMsg{}
	})
}
