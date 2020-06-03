package main

import (
	"log"
	"os"

	"github.com/samuelstevens/goimglabeler/api"
	"github.com/samuelstevens/goimglabeler/project"
)

func loadRoot() string {
	rootDir := os.Getenv("PROJECT_DIR")

	if rootDir == "" {
		log.Fatal("Need a project directory to scrape!")
	}

	return rootDir
}

func main() {

	project := project.New(loadRoot(), api.New())

	err := project.UpdateImages()

	if err != nil {
		log.Fatal(err.Error())
	}
}
