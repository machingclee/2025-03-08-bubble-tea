package main

import (
	"log"
	"os"
	"project_generator/internal/cli"
)

const logFilePath = "project_craetor.log"

var cleanup func() error

func init() {

	// Open or create the log file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	// Set log output to the file
	log.SetOutput(logFile)

	// Cleanup function to close the file
	cleanup = func() error {
		return logFile.Close()
	}
}

func main() {
	defer func() {
		if err := cleanup(); err != nil {
			log.Printf("Error during cleanup: %v\n", err)
			os.Exit(1)
		}
	}()

	cli.Start()
}
