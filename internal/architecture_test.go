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

func TestCorePackagesPurity(t *testing.T) {
	forbidden := map[string]bool{
		"os":       true,
		"net":      true,
		"net/http": true,
		"syscall":  true,
		"io/fs":    true,
		"os/exec":  true,
	}

	out, err := exec.Command("go", "list", "-f", "{{.ImportPath}}\t{{join .Imports ` `}}",
		"github.com/llttlltt/dj-library-tools/internal/core/...").CombinedOutput()
	if err != nil {
		t.Fatalf("failed to list core packages: %v\nOutput: %s", err, string(out))
	}

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		pkg := parts[0]
		if len(parts) < 2 {
			continue
		}
		for _, imp := range strings.Fields(parts[1]) {
			if forbidden[imp] {
				t.Errorf("Core purity violation: %s imports forbidden package %q", pkg, imp)
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
