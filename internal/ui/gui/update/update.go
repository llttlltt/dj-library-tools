package update

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

const Repo = "llttlltt/dj-library-tools"

// UpdateInfo carries information about a found update.
type UpdateInfo struct {
	Available    bool   `json:"available"`
	Version      string `json:"version"`
	Current      string `json:"current"`
	ReleaseNotes string `json:"release_notes"`
	URL          string `json:"url"`
}

// Check queries GitHub for the latest release and compares it with the current version.
func Check(currentVersion string) (*UpdateInfo, error) {
	// In development mode, we use a placeholder version that is always
	// lower than any real release.
	semverVersion := currentVersion
	if semverVersion == "v0.0.0-dev" || semverVersion == "v0.0.0" {
		semverVersion = "0.0.0"
	}

	v, err := semver.ParseTolerant(semverVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse current version %q: %w", currentVersion, err)
	}

	latest, found, err := selfupdate.DetectLatest(Repo)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}

	// Fallback: If not found, it might be a platform mismatch (common in dev).
	// Try to get the latest release version string directly via GitHub API.
	if !found {
		resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", Repo))
		if err == nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
			var githubRel struct {
				TagName string `json:"tag_name"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&githubRel); err == nil && githubRel.TagName != "" {
				latestVer, err := semver.ParseTolerant(githubRel.TagName)
				if err == nil {
					latest = &selfupdate.Release{
						Version: latestVer,
					}
					found = true
				}
			}
		}
	}

	if !found || latest == nil {
		return &UpdateInfo{
			Available: false,
			Current:   currentVersion,
			Version:   "",
		}, nil
	}

	isDev := currentVersion == "v0.0.0-dev" || currentVersion == "v0.0.0"

	info := &UpdateInfo{
		Available:    !isDev && latest.Version.GT(v),
		Version:      latest.Version.String(),
		Current:      currentVersion,
		ReleaseNotes: latest.ReleaseNotes,
		URL:          latest.AssetURL,
	}

	return info, nil
}

// Apply downloads and applies the latest update.
func Apply(currentVersion string) error {
	v, err := semver.ParseTolerant(currentVersion)
	if err != nil {
		return err
	}

	_, err = selfupdate.UpdateSelf(v, Repo)
	return err
}

// GetPermissionStatus checks if the application has write access to its own directory.
func GetPermissionStatus() string {
	execPath, err := os.Executable()
	if err != nil {
		return "Unknown"
	}

	// For macOS .app bundles, we check if we can write inside the bundle
	targetDir := filepath.Dir(execPath)
	if runtime.GOOS == "darwin" && filepath.Base(targetDir) == "MacOS" {
		// Inside /Contents/MacOS, go up to /Contents
		targetDir = filepath.Dir(targetDir)
	}

	tempFile := filepath.Join(targetDir, ".write_test")
	err = ioutil.WriteFile(tempFile, []byte("test"), 0644)
	if err != nil {
		return "Limited"
	}
	_ = os.Remove(tempFile)
	return "Healthy"
}

// FixPermissions attempts to fix write permissions by changing ownership of the app bundle.
// Currently only implemented for macOS.
func FixPermissions() error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("permission fixing is only supported on macOS")
	}

	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// Identify the .app bundle path.
	// Binary is at: .../AppName.app/Contents/MacOS/binary
	bundlePath := execPath
	if strings.Contains(bundlePath, ".app/Contents/MacOS/") {
		bundlePath = strings.Split(bundlePath, ".app/Contents/MacOS/")[0] + ".app"
	} else {
		return fmt.Errorf("could not identify .app bundle path from %s", execPath)
	}

	// We use osascript to run chown with admin privileges.
	// This will trigger the standard macOS password/TouchID prompt.
	user := os.Getenv("USER")
	if user == "" {
		// Fallback to id command if USER env is not set
		out, err := exec.Command("id", "-un").Output()
		if err != nil {
			return fmt.Errorf("could not determine current user: %w", err)
		}
		user = strings.TrimSpace(string(out))
	}

	script := fmt.Sprintf("do shell script \"chown -R %s \\\"%s\\\"\" with administrator privileges", user, bundlePath)
	cmd := exec.Command("osascript", "-e", script)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to fix permissions: %s: %w", string(out), err)
	}

	return nil
}
