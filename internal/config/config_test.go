package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// overrideConfigDir redirects all config helpers to a temp directory for the
// duration of the test by setting XDG_CONFIG_HOME.
func overrideConfigDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)
	return filepath.Join(tmp, "djlt")
}

// ── Connection ────────────────────────────────────────────────────────────────

func TestConnection_Roundtrip(t *testing.T) {
	overrideConfigDir(t)

	s := Connection{
		ID:       NewConnectionID(),
		Name:     "Main Library",
		Provider: "rb",
		Config:   map[string]string{"file_path": "/Music/rekordbox.xml"},
	}

	require.NoError(t, SaveConnection(s))

	loaded, err := LoadConnections()
	require.NoError(t, err)
	require.Len(t, loaded, 1)
	assert.Equal(t, s.ID, loaded[0].ID)
	assert.Equal(t, s.Name, loaded[0].Name)
	assert.Equal(t, s.Provider, loaded[0].Provider)
	assert.Equal(t, s.Config["file_path"], loaded[0].Config["file_path"])
}

func TestConnection_Delete(t *testing.T) {
	overrideConfigDir(t)

	s := Connection{ID: NewConnectionID(), Name: "Tmp", Provider: "rb", Config: map[string]string{}}
	require.NoError(t, SaveConnection(s))

	require.NoError(t, DeleteConnection(s.ID))

	loaded, err := LoadConnections()
	require.NoError(t, err)
	assert.Empty(t, loaded)
}

func TestFindFirstConnection(t *testing.T) {
	overrideConfigDir(t)

	rb := Connection{ID: NewConnectionID(), Name: "RB", Provider: "rb", Config: map[string]string{"file_path": "/a.xml"}}
	px := Connection{ID: NewConnectionID(), Name: "Plex", Provider: "plex", Config: map[string]string{"host": "localhost"}}
	require.NoError(t, SaveConnection(rb))
	require.NoError(t, SaveConnection(px))

	found, err := FindFirstConnection("rb")
	require.NoError(t, err)
	assert.Equal(t, rb.ID, found.ID)

	_, err = FindFirstConnection("m3u")
	assert.Error(t, err)
}

// ── Workflow ──────────────────────────────────────────────────────────────────

func TestWorkflow_Roundtrip(t *testing.T) {
	overrideConfigDir(t)

	w := Workflow{
		ID:   NewWorkflowID(),
		Name: "Nightly Sync",
		Steps: []Step{
			{
				ID:   NewStepID(),
				Kind: "sync",
				Source: Endpoint{
					ConnectionID: NewConnectionID(),
					Resource:     "tracks",
					Query:        "playlists:Inbox",
				},
				Targets: []Endpoint{
					{ConnectionID: NewConnectionID(), Resource: "playlists", Query: "name:Target"},
				},
			},
		},
	}

	require.NoError(t, SaveWorkflow(w))

	loaded, err := LoadWorkflows()
	require.NoError(t, err)
	require.Len(t, loaded, 1)
	assert.Equal(t, w.ID, loaded[0].ID)
	assert.Equal(t, w.Name, loaded[0].Name)
	require.Len(t, loaded[0].Steps, 1)
	assert.Equal(t, w.Steps[0].ID, loaded[0].Steps[0].ID)
	assert.Equal(t, w.Steps[0].Kind, loaded[0].Steps[0].Kind)
	assert.Equal(t, w.Steps[0].Source.Query, loaded[0].Steps[0].Source.Query)
}

func TestWorkflow_Delete(t *testing.T) {
	overrideConfigDir(t)

	w := Workflow{ID: NewWorkflowID(), Name: "Tmp", Steps: []Step{}}
	require.NoError(t, SaveWorkflow(w))
	require.NoError(t, DeleteWorkflow(w.ID))

	loaded, err := LoadWorkflows()
	require.NoError(t, err)
	assert.Empty(t, loaded)
}

// ── PathMap ───────────────────────────────────────────────────────────────────

func TestPathMap_Roundtrip(t *testing.T) {
	overrideConfigDir(t)

	pm := PathMap{
		ID:            NewPathMapID(),
		ConnectionAID: NewConnectionID(),
		ConnectionBID: NewConnectionID(),
		Rules: []PathRule{
			{From: "/Volumes/Music/", To: "/media/music/"},
		},
	}

	require.NoError(t, SavePathMap(pm))

	loaded, err := LoadPathMaps()
	require.NoError(t, err)
	require.Len(t, loaded, 1)
	assert.Equal(t, pm.ID, loaded[0].ID)
	assert.Equal(t, pm.ConnectionAID, loaded[0].ConnectionAID)
	assert.Equal(t, pm.ConnectionBID, loaded[0].ConnectionBID)
	require.Len(t, loaded[0].Rules, 1)
	assert.Equal(t, pm.Rules[0].From, loaded[0].Rules[0].From)
	assert.Equal(t, pm.Rules[0].To, loaded[0].Rules[0].To)
}

func TestPathMap_Delete(t *testing.T) {
	overrideConfigDir(t)

	pm := PathMap{ID: NewPathMapID(), ConnectionAID: NewConnectionID(), ConnectionBID: NewConnectionID(), Rules: []PathRule{}}
	require.NoError(t, SavePathMap(pm))
	require.NoError(t, DeletePathMap(pm.ID))

	loaded, err := LoadPathMaps()
	require.NoError(t, err)
	assert.Empty(t, loaded)
}

func TestLoadConnections_EmptyDir(t *testing.T) {
	overrideConfigDir(t)
	// Ensure the connections dir exists but is empty
	dir, err := GetConnectionsDir()
	require.NoError(t, err)
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Empty(t, entries)

	connections, err := LoadConnections()
	require.NoError(t, err)
	assert.Nil(t, connections)
}
