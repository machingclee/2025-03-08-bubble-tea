package cli

import (
	"fmt"
	"log"
	"os"
	"project_generator/internal/projgenerator"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	ProjectNamePage int = iota
	ConfigTypePage
	RouterTypePage
)

type View interface {
	View() string
	Update(msg tea.Msg, m *Model) tea.Cmd
}

type AppConfigurable interface {
	ToAppConfig() projgenerator.AppConfig
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

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println("cli update() msg:", msg)
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		}
	case createProjectMsg:
		log.Println("Createing project")
		appConfig := populateAppConfig(m.Views)
		go projgenerator.GenerateProject(appConfig, m.ProgressChannel)
		return m, listenToProgress(m.ProgressChannel)
	}
	log.Println(" m.CurrentViewIndex < len(m.Views)", m.CurrentViewIndex < len(m.Views))
	if m.CurrentViewIndex < len(m.Views) {
		cmd = m.Views[m.CurrentViewIndex].Update(msg, m)
	}
	return m, cmd
}

func (m *Model) View() string {
	log.Println("Cli View() > m.CreateProject: ", m.CreateProject)
	if m.Quitting {
		return "See you later!"
	} else if m.CreateProject {

	}
	log.Println("Cli View() m.", m.CurrentViewIndex)
	log.Println("m.Views[m.CurrentViewIndex]", m.Views[m.CurrentViewIndex])
	var results string
	for index := 0; index <= m.CurrentViewIndex; index++ {
		results += m.Views[index].View() + "\n"
	}
	return results
}

type createProjectMsg struct{}

func Start() {
	initialModel := Model{
		Views: []View{
			NewProjectNameView(),
			NewMultiChoiceView(
				"Read configuration settings from:",
				[]string{
					"Command-line flas",
					"Environment variables",
				},
				configTypeConfigurator,
			),
			NewMultiChoiceView(
				"Pick your preferred router:",
				[]string{
					"Gorilla Mux",
					"HttpRouter",
				},
				routerTypeConfigurator,
			),
		},
		CurrentViewIndex: 0,
		Progress:         progress.New(),
		ProgressChannel:  make(chan tea.Msg),
	}

	if _, err := tea.NewProgram(&initialModel).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start program: %v\n", err)
	}
}

func populateAppConfig(views []View) projgenerator.AppConfig {
	appConfig := projgenerator.AppConfig{}
	for _, view := range views {
		if appConfigurable, ok := view.(AppConfigurable); ok {
			appConfig = mergeAppConfigs(appConfig, appConfigurable.ToAppConfig())
		}
	}

	return appConfig
}

func mergeAppConfigs(dest, source projgenerator.AppConfig) projgenerator.AppConfig {
	if source.ProjectName != "" {
		dest.ProjectName = source.ProjectName
	}

	if source.UseRouter {
		dest.UseRouter = source.UseRouter
		dest.RouterType = source.RouterType
		dest.RouterImportPath = source.RouterImportPath
		dest.RouterConstructor = source.RouterConstructor
	}
	return source
}

func listenToProgress(progessChannel <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-progessChannel
		if !ok {
			return nil
		}
		return msg
	}
}
