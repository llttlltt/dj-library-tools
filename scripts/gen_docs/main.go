package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/ui/cli"
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
	isParent := cmd.HasSubCommands()

	var path string
	if isRoot {
		path = filepath.Join(dir, "index.md")
	} else if isParent {
		subDir := filepath.Join(dir, name)
		if err := os.MkdirAll(subDir, 0755); err != nil {
			return err
		}
		path = filepath.Join(subDir, "index.md")
	} else {
		path = filepath.Join(dir, name+".md")
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := writeCommandDoc(cmd, f, path); err != nil {
		return err
	}

	if isParent || isRoot {
		subDir := dir
		if !isRoot {
			subDir = filepath.Join(dir, name)
		}
		for _, child := range cmd.Commands() {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			if err := generateDocs(child, subDir); err != nil {
				return err
			}
		}
	}

	return nil
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

	// Also fix direct leaf links to parent commands
	reLeaf := regexp.MustCompile(`\(([a-z0-9_]+)\.md\)`)
	content = reLeaf.ReplaceAllStringFunc(content, func(m string) string {
		match := reLeaf.FindStringSubmatch(m)[1]
		// Special case for config, plex, rb which are now directories
		if match == "config" || match == "plex" || match == "rb" {
			return "(" + match + "/index.md)"
		}
		return m
	})

	return content
}

func formatLongDesc(lines []string) []string {
	var result []string
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if line == "" && (i == 0 || i == len(lines)-1) { continue }
		result = append(result, line)

		if strings.HasPrefix(line, "**") && i+1 < len(lines) {
			nextLine := lines[i+1]
			if strings.HasPrefix(nextLine, "  ") {
				result = append(result, "```bash")
				j := i + 1
				for ; j < len(lines); j++ {
					if lines[j] == "" || strings.HasPrefix(lines[j], "  ") {
						result = append(result, strings.TrimPrefix(lines[j], "  "))
					} else { break }
				}
				result = append(result, "```")
				i = j - 1
			}
		}
	}
	return result
}
