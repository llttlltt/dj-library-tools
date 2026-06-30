package internal

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArchitectureBoundaries(t *testing.T) {
	root := "."
	
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }
		if !info.IsDir() || strings.Contains(path, "ui") || strings.HasPrefix(path, ".") || path == "scripts" || path == "cmd" || path == "tests" || path == "plan" || path == "docs" { return nil }

		// Check each package not under internal/ui
		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, path, nil, parser.ImportsOnly)
		if err != nil { return nil } // Skip non-go dirs

		for _, pkg := range pkgs {
			for _, file := range pkg.Files {
				for _, imp := range file.Imports {
					importPath := strings.Trim(imp.Path.Value, "\"")
					if strings.Contains(importPath, "github.com/llttlltt/dj-library-tools/internal/ui") {
						t.Errorf("Boundary violation: package %q imports %q", path, importPath)
					}
					
					// Core should not depend on other internal packages
					if strings.Contains(path, "core") && 
						strings.Contains(importPath, "github.com/llttlltt/dj-library-tools/internal/") &&
						!strings.Contains(importPath, "core") {
						t.Errorf("Core violation: package %q imports %q", path, importPath)
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}
