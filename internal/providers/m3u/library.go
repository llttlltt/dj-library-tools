package m3u

import (
	"github.com/llttlltt/dj-library-tools/internal/core/models"
)

type Library struct {
	tracks []models.Track
}

func NewLibrary(tracks []models.Track) *Library {
	return &Library{tracks: tracks}
}

func (l *Library) GetResources(kind string) []models.Resource {
	var results []models.Resource
	switch kind {
	case "track":
		for _, t := range l.tracks {
			results = append(results, t)
		}
	case "group":
		// M3U is usually a single playlist, so we represent it as one group
		// We'll handle the group discovery in the provider for now since M3U is file-specific
	}
	return results
}

func (l *Library) GetMembershipMap() map[string][]models.PlaylistMembership {
	return make(map[string][]models.PlaylistMembership)
}

func (l *Library) CreateGroup(parentID, name string, groupType models.GroupKind, position int) (models.ResourceGroup, error) {
	return models.ResourceGroup{}, nil
}

func (l *Library) DeleteGroup(groupID string, groupType models.GroupKind) error {
	return nil
}

func (l *Library) AddTracks(groupID string, trackIDs []string) (int, error) {
	// For M3U, we don't use trackIDs (strings) as easily because the "ID" is often the path.
	// We'll keep the track-based Add in the provider for now or expand this interface.
	return 0, nil
}

func (l *Library) RemoveTracks(groupID string, trackIDs []string) (int, error) {
	return 0, nil
}

func (l *Library) UpdateGroup(groupID string, trackIDs []string) error {
	return nil
}

func (l *Library) RenameGroup(groupID, newName string, groupType models.GroupKind) error {
	return nil
}

func (l *Library) MoveGroup(groupID string, groupType models.GroupKind, targetParentID string) error {
	return nil
}

func (l *Library) Save(path string) error {
	return nil
}
