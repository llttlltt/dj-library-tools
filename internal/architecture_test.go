package internal

import (
	"os/exec"
	"strings"
	"testing"
)

func TestArchitectureBoundaries(t *testing.T) {
	packages := []string{
		"github.com/llttlltt/dj-library-tools/internal/core/...",
		"github.com/llttlltt/dj-library-tools/internal/infra/...",
		"github.com/llttlltt/dj-library-tools/internal/providers/...",
		"github.com/llttlltt/dj-library-tools/internal/services/...",
	}

	for _, pkg := range packages {
		out, err := exec.Command("go", "list", "-deps", pkg).CombinedOutput()
		if err != nil {
			t.Fatalf("failed to list dependencies for %s: %v\nOutput: %s", pkg, err, string(out))
		}

		deps := strings.Split(string(out), "\n")
		for _, dep := range deps {
			if strings.Contains(dep, "github.com/llttlltt/dj-library-tools/internal/ui") {
				t.Errorf("Boundary violation: package in %q depends on %q", pkg, dep)
			}

			// Core should not depend on other internal packages
			if strings.Contains(pkg, "/core/") &&
				strings.Contains(dep, "github.com/llttlltt/dj-library-tools/internal/") &&
				!strings.Contains(dep, "/core/") {
				t.Errorf("Core violation: core package depends on %q", dep)
			}
		}
	}
}

func TestNoDirectOutput(t *testing.T) {
	dirs := []string{"core", "infra", "providers", "services"}

	for _, dir := range dirs {
		// Use grep to find direct stdout/stderr usage
		// Fails if fmt.Print, fmt.Printf, fmt.Println, or os.Stdout/os.Stderr are found.
		cmd := exec.Command("grep", "-rE", "fmt\\.Print|os\\.Stdout|os\\.Stderr", dir)
		out, _ := cmd.CombinedOutput()

		lines := strings.Split(string(out), "\n")
		var violations []string
		for _, line := range lines {
			if line == "" {
				continue
			}
			// Ignore test files
			if strings.Contains(line, "_test.go") {
				continue
			}
			violations = append(violations, line)
		}

		if len(violations) > 0 {
			t.Errorf("Direct output violation in %s:\n%s", dir, strings.Join(violations, "\n"))
		}
	}
}
