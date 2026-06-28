package main

import (
	"log"
	"os"

	"github.com/llttlltt/dj-library-tools/internal/cli"
	"github.com/spf13/cobra/doc"
)

func main() {
	targetDir := "./docs/commands"
	// Clean the directory
	os.RemoveAll(targetDir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Fatal(err)
	}

	// Generate standard Cobra markdown documentation
	if err := doc.GenMarkdownTree(cli.RootCmd, targetDir); err != nil {
		log.Fatal(err)
	}

	// Rename djlt.md to index.md for MkDocs
	if err := os.Rename(targetDir+"/djlt.md", targetDir+"/index.md"); err != nil {
		log.Fatal(err)
	}

	log.Printf("Generated native command documentation in %s", targetDir)
}
