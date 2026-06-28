package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	targetDir := "./docs/commands"
	os.RemoveAll(targetDir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Fatal(err)
	}

	if err := generateDocs(cli.RootCmd, targetDir); err != nil {
		log.Fatal(err)
	}
	log.Printf("Generated hierarchical command documentation in %s", targetDir)
}

func generateDocs(cmd *cobra.Command, dir string) error {
	name := cmd.Name()
	isRoot := cmd == cli.RootCmd
	isConfigChild := strings.Contains(cmd.CommandPath(), "djlt config")
	isPlexOrRB := name == "plex" || name == "rb"

	// 1. Root Command
	if isRoot {
		path := filepath.Join(dir, "index.md")
		if err := writeDocToFile(cmd, path); err != nil {
			return err
		}
		for _, child := range cmd.Commands() {
			if !child.IsAvailableCommand() {
				continue
			}
			if err := generateDocs(child, dir); err != nil {
				return err
			}
		}
		return nil
	}

	// 2. Special Case: Plex and RB within Config (Merge children)
	if isConfigChild && isPlexOrRB {
		subDir := filepath.Join(dir, name)
		os.MkdirAll(subDir, 0755)
		path := filepath.Join(subDir, "index.md")
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := writeCommandDoc(cmd, f, path); err != nil {
			return err
		}

		for _, child := range cmd.Commands() {
			if !child.IsAvailableCommand() {
				continue
			}
			fmt.Fprintf(f, "\n\n---\n\n")
			if err := writeCommandDoc(child, f, path); err != nil {
				return err
			}
		}
		return nil
	}

	// 3. Parent commands get a directory and an index.md
	if cmd.HasSubCommands() {
		subDir := filepath.Join(dir, name)
		os.MkdirAll(subDir, 0755)
		path := filepath.Join(subDir, "index.md")
		if err := writeDocToFile(cmd, path); err != nil {
			return err
		}
		for _, child := range cmd.Commands() {
			if !child.IsAvailableCommand() {
				continue
			}
			if err := generateDocs(child, subDir); err != nil {
				return err
			}
		}
		return nil
	}

	// 4. Leaf commands
	path := filepath.Join(dir, name+".md")
	return writeDocToFile(cmd, path)
}

func writeDocToFile(cmd *cobra.Command, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return writeCommandDoc(cmd, f, path)
}

func writeCommandDoc(cmd *cobra.Command, w io.Writer, filePath string) error {
	var buf strings.Builder
	if err := doc.GenMarkdownCustom(cmd, &buf, func(s string) string { return s }); err != nil {
		return err
	}
	content := buf.String()
	content = fixLinks(content, filePath)

	lines := strings.Split(content, "\n")
	var header, usage, options, inherited, longDesc, seeAlso []string
	var current *[]string = &header

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "## ") && current == &header {
			heading := strings.TrimPrefix(line, "## ")
			if strings.HasPrefix(heading, "djlt ") && heading != "djlt" {
				heading = strings.TrimPrefix(heading, "djlt ")
			}
			line = "# " + heading
		}

		switch {
		case line == "### Synopsis":
			current = &longDesc
			continue
		case strings.HasPrefix(line, "```") && current == &longDesc && len(usage) == 0:
			current = &usage
		case line == "### Options":
			current = &options
		case line == "### Options inherited from parent commands":
			line = "### Inherited Options"
			current = &inherited
		case line == "### SEE ALSO":
			line = "## See also"
			current = &seeAlso
		case strings.HasPrefix(line, "## ") && current != &header:
			current = &longDesc
		}

		*current = append(*current, line)

		if current == &usage && strings.HasPrefix(line, "```") && len(usage) > 1 {
			current = &longDesc
		}
	}

	longDesc = formatLongDesc(longDesc)

	var output []string
	output = append(output, header...)
	output = append(output, usage...)
	output = append(output, options...)
	output = append(output, inherited...)
	output = append(output, longDesc...)
	output = append(output, seeAlso...)

	final := strings.Join(output, "\n")
	footerStart := strings.Index(final, "###### Auto generated")
	if footerStart != -1 {
		final = final[:footerStart]
	}

	_, err := io.WriteString(w, strings.TrimSpace(final))
	return err
}

func fixLinks(content string, filePath string) string {
	commandsDir, _ := filepath.Abs("docs/commands")
	currentDir, _ := filepath.Abs(filepath.Dir(filePath))
	rootRel, _ := filepath.Rel(currentDir, filepath.Join(commandsDir, "index.md"))
	content = strings.ReplaceAll(content, "(djlt.md)", "("+rootRel+")")

	re := regexp.MustCompile(`\(djlt_([a-z0-9_]+)\.md\)`)
	content = re.ReplaceAllStringFunc(content, func(m string) string {
		match := re.FindStringSubmatch(m)[1]
		parts := strings.Split(match, "_")

		targetPath := filepath.Join(commandsDir, strings.Join(parts, "/"))
		if _, err := os.Stat(targetPath); err == nil {
			targetPath = filepath.Join(targetPath, "index.md")
		} else {
			targetPath += ".md"
		}

		rel, _ := filepath.Rel(currentDir, targetPath)
		return "(" + rel + ")"
	})
	return content
}

func formatLongDesc(lines []string) []string {
	var result []string
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if line == "" && (i == 0 || i == len(lines)-1) {
			continue
		}
		result = append(result, line)

		if strings.HasPrefix(line, "**") && i+1 < len(lines) {
			nextLine := lines[i+1]
			if strings.HasPrefix(nextLine, "  ") {
				result = append(result, "```bash")
				j := i + 1
				for ; j < len(lines); j++ {
					if lines[j] == "" || strings.HasPrefix(lines[j], "  ") {
						result = append(result, strings.TrimPrefix(lines[j], "  "))
					} else {
						break
					}
				}
				result = append(result, "```")
				i = j - 1
			}
		}
	}
	return result
}
