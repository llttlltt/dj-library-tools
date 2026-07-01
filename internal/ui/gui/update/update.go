package update

import (
	"fmt"
	"io/ioutil"
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
	v, err := semver.ParseTolerant(currentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse current version %q: %w", currentVersion, err)
	}

	latest, found, err := selfupdate.DetectLatest(Repo)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}

	if !found || latest.Version.LTE(v) {
		return &UpdateInfo{Available: false, Current: currentVersion}, nil
	}

	return &UpdateInfo{
		Available:    true,
		Version:      latest.Version.String(),
		Current:      currentVersion,
		ReleaseNotes: latest.ReleaseNotes,
		URL:          latest.AssetURL,
	}, nil
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
