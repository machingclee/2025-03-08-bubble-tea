package cli

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectNameView struct {
	inputModel    textinput.Model
	endingMessage string
}

func NewProjectNameView() *ProjectNameView {
	inputModel := newTextInput("Project Name")
	return &ProjectNameView{
		inputModel:    inputModel,
		endingMessage: "",
	}
}

func newTextInput(prompt string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = prompt
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return ti
}

func (v *ProjectNameView) View() string {
	endingDisplay := func() string {
		if len(v.endingMessage) > 0 {
			return fmt.Sprintf("\n\nNice, the project name \"%v\" is well received.\n\n", v.inputModel.Value())
		}
		return ""
	}()
	return "Please input a project name: \n\n" + v.inputModel.View() + endingDisplay
}

func (v *ProjectNameView) Update(msg tea.Msg, m *Model) tea.Cmd {
	log.Println("ProjectNameView.Update() > msg:", msg)

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			v.endingMessage = v.inputModel.Value()
			m.CurrentViewIndex++
		}
	}
	v.inputModel, cmd = v.inputModel.Update(msg)
	return cmd
}
