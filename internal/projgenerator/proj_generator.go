package projgenerator

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	ProgressCreationInProgressMsg float64
	ProjectCreatedMsg             struct{}
)

const (
	webAppTemplateFileLocation         = "./templates/webapp.tmpl"
	progressAfterPrepOfOutputDirectory = 0.3
	progressAfterCreationOfMainFile    = 0.3
)

type AppConfig struct {
	ProjectName       string
	UseRouter         bool
	RouterType        string
	RouterImportPath  string
	RouterConstructor string
	ConfigSource      string
}

func GenerateProject(appConfig AppConfig, progressChannel chan tea.Msg) {
	defer close(progressChannel)

	outputPath := "../" + appConfig.ProjectName
	prepareOutputDirectory(outputPath)
	progressChannel <- ProgressCreationInProgressMsg(progressAfterPrepOfOutputDirectory)
	generateMainFile(appConfig, outputPath)
	log.Println("after generate main file")
	progressChannel <- ProgressCreationInProgressMsg(progressAfterCreationOfMainFile)
	initGoModuleAndDependencies(outputPath, appConfig)
	log.Println("after initGoModuleAndDependencies")
	log.Println("project created")
}

func prepareOutputDirectory(outputPath string) {
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		err := os.RemoveAll(outputPath)
		if err != nil {
			log.Fatalf("Failed to remove existing output directory: %v", err)
		}
	}

	err := os.MkdirAll(outputPath, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
}

func initGoModuleAndDependencies(outputPath string, config AppConfig) {
	projectName := config.ProjectName
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Dir = outputPath
	if err := cmd.Run(); err != nil {
		log.Fatalf("Faield in init go.mod: %v", err)
	}
	if config.UseRouter {
		cmd = exec.Command("go", "get", config.RouterImportPath)
		cmd.Dir = outputPath
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to download router package")
		}
	}
}

func generateMainFile(config AppConfig, outputPath string) {
	mainFilePath := filepath.Join(outputPath, "main.go")
	tmpl, err := template.ParseFiles(webAppTemplateFileLocation)
	if err != nil {
		log.Fatalf("Error opening template file: %v", err)
	}
	f, err := os.Create(mainFilePath)
	if err != nil {
		log.Fatalf("Failed to create main.go at %s: %v", mainFilePath, err)
	}
	err = tmpl.Execute(f, config)

	if err != nil {
		log.Fatalf("Failed")
	}

}
