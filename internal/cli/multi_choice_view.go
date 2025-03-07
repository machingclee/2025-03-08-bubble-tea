package cli

import (
	"log"
	"project_generator/internal/termstyle"
	"strings"

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
		builder.WriteString((checkbox + "\n"))
	}

	instructions := termstyle.Subtle("enter:choose") +
		termstyle.Dot + termstyle.Subtle("esc or ctrl-c: quit")
	builder.WriteString("\n" + instructions)

	return builder.String()
}

func (v *MultiChoiceView) Update(msg tea.Msg, m *Model) tea.cmd {

}
