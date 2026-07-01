package cli

import (
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/providers/factory"
	"github.com/spf13/cobra"
)

func newProvidersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "providers",
		Short: "List registered providers and their capabilities",
		Run: func(cmd *cobra.Command, args []string) {
			infos := factory.ListProviders()

			headers := []string{"Provider", "Resources", "Write", "Meta", "Groups", "Cues", "Grid", "File"}
			var rows [][]string

			for _, info := range infos {
				caps := info.Capabilities

				var resNames []string
				for _, r := range info.Resources {
					resNames = append(resNames, r.Name)
				}

				rows = append(rows, []string{
					info.Name,
					strings.Join(resNames, ", "),
					fmtBool(caps.CanWrite),
					fmtBool(caps.CanUpdateMetadata),
					fmtBool(caps.CanManageGroups),
					fmtBool(caps.SupportsCues),
					fmtBool(caps.SupportsBeatgrids),
					fmtBool(caps.IsFileBased),
				})
			}

			renderTable(headers, rows)
		},
	}
}

func fmtBool(v bool) string {
	if v {
		return "YES"
	}
	return "NO"
}

func init() {
	// Re-sort alphabetically if desired, but factory.ListProviders already sorts.
}
