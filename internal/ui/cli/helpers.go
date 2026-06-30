package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/config"
	djerrors "github.com/llttlltt/dj-library-tools/internal/core/errors"
	"github.com/llttlltt/dj-library-tools/internal/services/orchestrator"
)

func getOrchestrator() *orchestrator.Orchestrator {
	cfg, _ := config.LoadAppConfig()
	opts := orchestrator.Options{
		RekordboxPrimaryPath: cfg.Rekordbox.PrimaryFilePath,
	}
	return orchestrator.New(&TerminalFeedback{}, opts)
}

func getRunOptions() orchestrator.RunOptions {
	return orchestrator.RunOptions{
		FilePath: filePath,
		Apply:    apply,
		Verbose:  verbose,
	}
}

// HandleError provides user-friendly messages for sentinel provider errors.
func HandleError(err error) error {
	if err == nil {
		return nil
	}

	kind := djerrors.KindOf(err)
	switch kind {
	case djerrors.KindReadOnly:
		return fmt.Errorf("operation failed: this provider is read-only")
	case djerrors.KindUnsupportedResource:
		return fmt.Errorf("operation failed: this resource type is not supported by the provider")
	case djerrors.KindInvalidParent:
		return fmt.Errorf("operation failed: cannot create the resource in that location (structural constraint)")
	case djerrors.KindNotFound:
		return fmt.Errorf("operation failed: resource not found")
	}

	return err
}

func stringsTitle(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
