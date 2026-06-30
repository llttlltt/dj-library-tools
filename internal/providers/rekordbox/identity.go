package rekordbox

import (
	"github.com/llttlltt/dj-library-tools/internal/core/models"
)

// Identify returns the authority-driven ID for a given group name.
// In Rekordbox, we currently use the name as the ID for simple lookups,
// but this package is the authority on how that ID is constructed.
func Identify(name string, groupType models.GroupKind) string {
	return name
}
