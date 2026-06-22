package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	targetDir := "./docs/commands"
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Fatal(err)
	}
	os.RemoveAll(targetDir)
	os.MkdirAll(targetDir, 0755)
	if err := generateDocs(cli.RootCmd, targetDir); err != nil {
		log.Fatal(err)
	}
	log.Printf("Generated hierarchical command documentation in %s", targetDir)
}

func generateDocs(cmd *cobra.Command, dir string) error {
	name := cmd.Name()
	var path string
	isParent := cmd.HasSubCommands() && cmd.Name() != "fix"

	if isParent {
		subDir := dir
		if cmd != cli.RootCmd {
			subDir = filepath.Join(dir, name)
		}
		if err := os.MkdirAll(subDir, 0755); err != nil {
			return err
		}
		path = filepath.Join(subDir, "index.md")
		for _, child := range cmd.Commands() {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			if err := generateDocs(child, subDir); err != nil {
				return err
			}
		}
	} else {
		path = filepath.Join(dir, name+".md")
	}

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

	// 1. Fix Links before splitting into lines
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
	// The structure is:
	// commands/index.md (djlt)
	// commands/auth.md (djlt_auth)
	// commands/playlist/index.md (djlt_playlist)
	// commands/playlist/fix.md (djlt_playlist_fix)

	isAtRoot := filepath.Dir(filePath) == "docs/commands"

	// 1. Fix root link
	if isAtRoot {
		content = strings.ReplaceAll(content, "(djlt.md)", "(index.md)")
	} else {
		content = strings.ReplaceAll(content, "(djlt.md)", "(../index.md)")
	}

	// 2. Fix subcommands
	// Standard cobra links are (djlt_subcommand.md)
	// We need to translate them.
	mappings := map[string]string{
		"djlt_auth.md":         "auth.md",
		"djlt_config.md":       "config.md",
		"djlt_folder.md":       "folder.md",
		"djlt_list.md":         "list.md",
		"djlt_metadata.md":     "metadata.md",
		"djlt_playlist.md":     "playlist/index.md",
		"djlt_playlist_fix.md": "playlist/fix.md",
		"djlt_stat.md":         "stat.md",
		"djlt_sync.md":         "sync.md",
	}

	// Adjust mappings if we are inside the playlist folder
	if !isAtRoot {
		mappings["djlt_playlist.md"] = "index.md"
		mappings["djlt_playlist_fix.md"] = "fix.md"
		// Everything else needs to go up one level
		for k, v := range mappings {
			if k != "djlt_playlist.md" && k != "djlt_playlist_fix.md" {
				mappings[k] = "../" + v
			}
		}
	}

	for old, new := range mappings {
		content = strings.ReplaceAll(content, "("+old+")", "("+new+")")
	}

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
