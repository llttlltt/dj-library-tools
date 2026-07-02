// Package gui exposes the Wails application bindings. App is the only type
// the frontend imports; all methods on App are callable as Go bindings.
package gui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers/factory"
	"github.com/llttlltt/dj-library-tools/internal/providers/plex"
	"github.com/llttlltt/dj-library-tools/internal/services/orchestrator"
	"github.com/llttlltt/dj-library-tools/internal/services/workflow"
	"github.com/llttlltt/dj-library-tools/internal/ui/gui/update"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	// Provider registrations.
	_ "github.com/llttlltt/dj-library-tools/internal/providers/m3u"
	_ "github.com/llttlltt/dj-library-tools/internal/providers/plex"
	_ "github.com/llttlltt/dj-library-tools/internal/providers/rekordbox"
)

// Version is the current application version, injected at build time via ldflags.
var Version = "v0.0.0-dev"

// App is the Wails application object. Its exported methods become TypeScript
// bindings in the frontend. All methods return JSON-serialisable types.
type App struct {
	ctx    context.Context
	orch   *orchestrator.Orchestrator
	engine *workflow.Engine
}

// NewApp constructs the App, resolving the primary rekordbox Connection at startup.
// If no Connection is found the orchestrator is still created; commands that
// require a file will fail with a clear error at call time.
func NewApp() *App {
	var primaryPath string
	if conn, err := config.FindFirstConnection("rb"); err == nil {
		primaryPath = config.ResolveProviderOptions(*conn).FilePath
	}

	orch := orchestrator.New(nil, orchestrator.Options{
		RekordboxPrimaryPath: primaryPath,
	})

	return &App{
		orch:   orch,
		engine: workflow.New(orch),
	}
}

// Startup is called by Wails when the application starts.
func (a *App) Startup(ctx context.Context) { a.ctx = ctx }

// ── System / Settings ────────────────────────────────────────────────────────

// GetVersion returns the current version of the application.
func (a *App) GetVersion() string {
	return Version
}

// GetUpdateConfig returns the current update-related configuration.
func (a *App) GetUpdateConfig() (config.UpdateConfig, error) {
	cfg, err := config.LoadAppConfig()
	if err != nil {
		return config.UpdateConfig{}, err
	}
	if cfg.Updates.CheckIntervalHour == 0 {
		cfg.Updates.CheckIntervalHour = 168 // Default 1 week
	}
	return cfg.Updates, nil
}

// SetUpdateInterval updates how frequently the app checks for updates.
func (a *App) SetUpdateInterval(hours int) error {
	cfg, err := config.LoadAppConfig()
	if err != nil {
		return err
	}
	cfg.Updates.CheckIntervalHour = hours
	return config.SaveAppConfig(cfg)
}

// CheckForUpdate queries GitHub for updates. If manual is false, it respects
// the check interval.
func (a *App) CheckForUpdate(manual bool) (*update.UpdateInfo, error) {
	cfg, err := config.LoadAppConfig()
	if err != nil {
		return nil, err
	}

	if !manual {
		interval := time.Duration(cfg.Updates.CheckIntervalHour) * time.Hour
		if interval == 0 {
			interval = 168 * time.Hour
		}

		lastCheck, _ := time.Parse(time.RFC3339, cfg.Updates.LastCheckAt)
		if time.Since(lastCheck) < interval {
			return &update.UpdateInfo{Available: false, Current: Version}, nil
		}
	}

	info, err := update.Check(Version)
	if err != nil {
		return nil, err
	}

	// Persist last check time
	cfg.Updates.LastCheckAt = time.Now().Format(time.RFC3339)
	_ = config.SaveAppConfig(cfg)

	return info, nil
}

// InstallUpdate downloads and applies the update.
func (a *App) InstallUpdate() error {
	return update.Apply(Version)
}

// GetPermissionStatus checks write access to the application bundle.
func (a *App) GetPermissionStatus() string {
	return update.GetPermissionStatus()
}

// FixPermissions triggers platform-specific escalation to fix permissions.
func (a *App) FixPermissions() error {
	return update.FixPermissions()
}

// ListProviders returns static metadata for all registered providers.
func (a *App) ListProviders() []factory.ProviderInfo {
	return factory.ListProviders()
}

// ── Connections ──────────────────────────────────────────────────────────────

// ListConnections returns all configured Connections.
func (a *App) ListConnections() ([]config.Connection, error) {
	return config.LoadConnections()
}

// CreateConnection generates a new UUID, saves the Connection file, and returns it.
func (a *App) CreateConnection(name, prov string, cfg map[string]string) (config.Connection, error) {
	if name == "" {
		return config.Connection{}, fmt.Errorf("connection name is required")
	}
	if prov == "" {
		return config.Connection{}, fmt.Errorf("connection provider is required")
	}
	s := config.Connection{
		ID:       config.NewConnectionID(),
		Name:     name,
		Provider: prov,
		Config:   cfg,
	}
	if err := config.SaveConnection(s); err != nil {
		return config.Connection{}, err
	}
	return s, nil
}

// DeleteConnection removes the Connection with the given ID.
func (a *App) DeleteConnection(id string) error {
	return config.DeleteConnection(id)
}

// UpdateConnection overwrites an existing Connection file (identified by s.ID).
func (a *App) UpdateConnection(s config.Connection) error {
	return config.SaveConnection(s)
}

// ── Workflows ─────────────────────────────────────────────────────────────────

// ListWorkflows returns all configured Workflows.
func (a *App) ListWorkflows() ([]config.Workflow, error) {
	return config.LoadWorkflows()
}

// GetWorkflow returns the Workflow with the given ID.
func (a *App) GetWorkflow(id string) (config.Workflow, error) {
	wfs, err := config.LoadWorkflows()
	if err != nil {
		return config.Workflow{}, err
	}
	for _, w := range wfs {
		if w.ID == id {
			return w, nil
		}
	}
	return config.Workflow{}, fmt.Errorf("workflow %q not found", id)
}

// SaveWorkflow assigns a UUID if the Workflow has no ID, saves it, and returns
// the saved Workflow (with ID populated).
func (a *App) SaveWorkflow(wf config.Workflow) (config.Workflow, error) {
	if wf.ID == "" {
		wf.ID = config.NewWorkflowID()
	}
	for i := range wf.Steps {
		if wf.Steps[i].ID == "" {
			wf.Steps[i].ID = config.NewStepID()
		}
	}
	if err := config.SaveWorkflow(wf); err != nil {
		return config.Workflow{}, err
	}
	return wf, nil
}

// DeleteWorkflow removes the Workflow with the given ID.
func (a *App) DeleteWorkflow(id string) error {
	return config.DeleteWorkflow(id)
}

// ── Execution ─────────────────────────────────────────────────────────────────

// PreviewWorkflow executes the Workflow with Apply=false and returns per-Step
// results showing what would change.
func (a *App) PreviewWorkflow(id string) (workflow.WorkflowResult, error) {
	wf, err := a.GetWorkflow(id)
	if err != nil {
		return workflow.WorkflowResult{}, err
	}
	return a.engine.Execute(a.ctx, wf, false)
}

// RunWorkflow executes the Workflow with Apply=true, committing all changes.
func (a *App) RunWorkflow(id string) (workflow.WorkflowResult, error) {
	wf, err := a.GetWorkflow(id)
	if err != nil {
		return workflow.WorkflowResult{}, err
	}
	return a.engine.Execute(a.ctx, wf, true)
}

// ── Diff ──────────────────────────────────────────────────────────────────────

// TrackRow is a lightweight track summary used in StepDiff lists.
type TrackRow struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	BPM      string `json:"bpm"`
	Location string `json:"location"`
}

// StepDiff holds per-target diff data for one Step.
type StepDiff struct {
	StepID     string     `json:"step_id"`
	TargetName string     `json:"target_name"`
	Current    []TrackRow `json:"current"`
	Added      []TrackRow `json:"added"`
	Removed    []TrackRow `json:"removed"`
	Unchanged  []TrackRow `json:"unchanged"`
}

// GetWorkflowDiff returns per-target, track-level diff data for every sync
// Step in the Workflow without applying any changes.
func (a *App) GetWorkflowDiff(id string) ([]StepDiff, error) {
	wf, err := a.GetWorkflow(id)
	if err != nil {
		return nil, err
	}

	runOpts := orchestrator.RunOptions{Apply: false}
	var out []StepDiff

	for _, step := range wf.Steps {
		kind := strings.ToLower(step.Kind)
		if kind != "sync" && kind != "add" && kind != "remove" && kind != "m3u_export" {
			continue
		}

		srcLoc, err := connectionProviderLoc(step.Source)
		if err != nil {
			return nil, fmt.Errorf("step %s source: %w", step.ID, err)
		}

		targetLocs := []string{}
		appendOnly := false

		switch kind {
		case "sync":
			for _, tgt := range step.Targets {
				tgtLoc, err := connectionProviderLoc(tgt)
				if err != nil {
					return nil, fmt.Errorf("step %s target: %w", step.ID, err)
				}
				if tgt.Query != "" {
					tgtLoc = tgtLoc + " " + tgt.Query
				}
				targetLocs = append(targetLocs, tgtLoc)
			}
		case "add":
			for _, tgt := range step.Targets {
				tgtLoc, err := connectionProviderLoc(tgt)
				if err != nil {
					return nil, fmt.Errorf("step %s target: %w", step.ID, err)
				}
				// For add, the query is the name of the new resource
				if tgt.Query != "" {
					tgtLoc = tgtLoc + " " + tgt.Query
				}
				targetLocs = append(targetLocs, tgtLoc)
			}
			appendOnly = true
		case "remove":
			for _, tgt := range step.Targets {
				tgtLoc, err := connectionProviderLoc(tgt)
				if err != nil {
					return nil, fmt.Errorf("step %s target: %w", step.ID, err)
				}
				if tgt.Query != "" {
					tgtLoc = tgtLoc + " " + tgt.Query
				}
				targetLocs = append(targetLocs, tgtLoc)
			}
		case "m3u_export":
			path, _ := step.Options["path"].(string)
			if path != "" {
				targetLocs = append(targetLocs, "m3u://"+path+"/playlists")
			}
			if a, ok := step.Options["append"].(bool); ok {
				appendOnly = a
			}
		}

		for _, tgtLoc := range targetLocs {
			diffs, err := a.orch.GetSyncDiff(a.ctx, srcLoc, tgtLoc, step.Source.Query, runOpts, appendOnly)
			if err != nil {
				return nil, fmt.Errorf("step %s diff: %w", step.ID, err)
			}

			for _, diff := range diffs {
				// Build the removed set for fast lookup.
				removedSet := make(map[string]bool, len(diff.RemovedIDs))
				for _, rid := range diff.RemovedIDs {
					removedSet[rid] = true
				}

				sd := StepDiff{
					StepID:     step.ID,
					TargetName: diff.TargetName,
					Current:    toTrackRows(diff.CurrentIDs, diff.TrackLookup),
					Added:      toTrackRows(diff.AddedIDs, diff.TrackLookup),
					Removed:    toTrackRows(diff.RemovedIDs, diff.TrackLookup),
				}

				if kind == "remove" {
					// For remove step, we only care about removals.
					sd.Added = []TrackRow{}
				} else if kind == "add" {
					// For add step, we only care about additions.
					sd.Removed = []TrackRow{}
				}

				// Unchanged = current members that are NOT being removed.
				var unchangedIDs []string
				for _, cid := range diff.CurrentIDs {
					if !removedSet[cid] {
						unchangedIDs = append(unchangedIDs, cid)
					}
				}
				sd.Unchanged = toTrackRows(unchangedIDs, diff.TrackLookup)
				out = append(out, sd)
			}
		}
	}
	return out, nil
}

// connectionProviderLoc resolves an Endpoint to "<provider>/<resource>".
func connectionProviderLoc(ep config.Endpoint) (string, error) {
	if ep.ConnectionID == "m3u" || ep.ConnectionID == "m3u8" {
		return ep.ConnectionID + "://" + ep.Resource, nil
	}
	conn, err := config.FindConnectionByID(ep.ConnectionID)
	if err != nil {
		return "", err
	}
	return conn.Provider + "/" + ep.Resource, nil
}

// toTrackRows converts a slice of track IDs to TrackRow summaries using the
// lookup map populated by GetSyncDiff.
func toTrackRows(ids []string, lookup map[string]models.Track) []TrackRow {
	rows := make([]TrackRow, 0, len(ids))
	for _, id := range ids {
		t := lookup[id]
		bpm := ""
		if t.BPM > 0 {
			bpm = fmt.Sprintf("%.0f", t.BPM)
		}
		rows = append(rows, TrackRow{
			ID:       id,
			Title:    t.Title,
			Artist:   t.Artist,
			BPM:      bpm,
			Location: t.Location,
		})
	}
	return rows
}

// GroupRow is a lightweight group (playlist/folder) summary used in QueryResult.
type GroupRow struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Kind   string `json:"kind"`
	Parent string `json:"parent"`
	Items  int    `json:"items"`
}

// QueryResult wraps the output of PreviewQuery, covering both track and group resources.
type QueryResult struct {
	Kind   string     `json:"kind"` // "tracks" or "groups"
	Tracks []TrackRow `json:"tracks"`
	Groups []GroupRow `json:"groups"`
	Count  int        `json:"count"`
}

// PreviewQuery is the GUI equivalent of: djlt ls <provider>/<resource> <query>
// It runs a read-only List against the given Connection and returns a QueryResult
// that covers both track resources and group resources (playlists/folders).
func (a *App) PreviewQuery(connectionID, resource, query string) (QueryResult, error) {
	runOpts := orchestrator.RunOptions{}
	res, err := a.orch.List(a.ctx, connectionID+"/"+resource, query, runOpts, "")
	if err != nil {
		return QueryResult{}, err
	}
	if len(res.Groups) > 0 || resource == "playlists" || resource == "folders" {
		rows := make([]GroupRow, 0, len(res.Groups))
		for _, g := range res.Groups {
			rows = append(rows, GroupRow{
				ID:     g.ID,
				Name:   g.Name,
				Kind:   string(g.Kind),
				Parent: g.ParentFolder,
				Items:  g.Items,
			})
		}
		return QueryResult{Kind: "groups", Groups: rows, Count: len(rows)}, nil
	}
	rows := make([]TrackRow, 0, len(res.Tracks))
	for _, t := range res.Tracks {
		rows = append(rows, TrackRow{
			ID:       t.ID,
			Title:    t.Title,
			Artist:   t.Artist,
			BPM:      fmt.Sprintf("%.2f", t.BPM),
			Location: t.Location,
		})
	}
	return QueryResult{Kind: "tracks", Tracks: rows, Count: len(rows)}, nil
}

// ── File picker ───────────────────────────────────────────────────────────────

// OpenFileDialog opens a native file-picker and returns the selected path, or
// an empty string if the user cancelled. defaultDir sets the initial directory;
// pass "" to let Wails use the system default.
func (a *App) OpenFileDialog(defaultDir string) (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Select Library File",
		DefaultDirectory: defaultDir,
		Filters: []runtime.FileFilter{
			{DisplayName: "Rekordbox XML", Pattern: "*.xml"},
			{DisplayName: "M3U Playlist", Pattern: "*.m3u;*.m3u8"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
}

// SaveFileDialog opens a native save-file-dialog and returns the selected path,
// or an empty string if the user cancelled.
func (a *App) SaveFileDialog(defaultDir, defaultFile string) (string, error) {
	return runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            "Save Playlist",
		DefaultDirectory: defaultDir,
		DefaultFilename:  defaultFile,
		Filters: []runtime.FileFilter{
			{DisplayName: "M3U Playlist", Pattern: "*.m3u"},
			{DisplayName: "M3U8 Playlist", Pattern: "*.m3u8"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
}

// ── Plex PIN auth ─────────────────────────────────────────────────────────────

// PlexAuthChallenge carries the auth URL and PIN ID returned by InitPlexAuth.
type PlexAuthChallenge struct {
	URL   string `json:"url"`
	PinID int    `json:"pin_id"`
}

// InitPlexAuth requests a new Plex PIN and returns the browser auth URL and the
// PIN ID the frontend should pass to CheckPlexAuth.
func (a *App) InitPlexAuth() (PlexAuthChallenge, error) {
	client := plex.NewClient("")
	pin, err := client.RequestPin(a.ctx)
	if err != nil {
		return PlexAuthChallenge{}, fmt.Errorf("failed to request Plex PIN: %w", err)
	}
	url := fmt.Sprintf(
		"https://app.plex.tv/auth/#!?code=%s&context%%5Bdevice%%5D%%5Bproduct%%5D=%s&clientID=%s",
		pin.Code, plex.ClientName, plex.AppID,
	)
	return PlexAuthChallenge{URL: url, PinID: pin.ID}, nil
}

// CheckPlexAuth polls whether the PIN has been authorised. Returns the token
// string when authenticated, or "" if not yet authorised (not an error).
func (a *App) CheckPlexAuth(pinID int) (string, error) {
	client := plex.NewClient("")
	status, err := client.CheckPin(a.ctx, pinID)
	if err != nil {
		return "", fmt.Errorf("failed to check Plex PIN: %w", err)
	}
	return status.AuthToken, nil
}
