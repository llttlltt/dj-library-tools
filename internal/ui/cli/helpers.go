package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/config"
	djerrors "github.com/llttlltt/dj-library-tools/internal/core/errors"
	"github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/llttlltt/dj-library-tools/internal/services/orchestrator"
	"github.com/llttlltt/dj-library-tools/internal/services/resolver"
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

func getExecContext() provider.ExecutionContext {
	return provider.ExecutionContext{
		Apply:    apply,
		Verbose:  verbose,
		Feedback: &TerminalFeedback{},
	}
}

func ResolveSelection(locStr string, queryOverride string) (*resolver.Selection, provider.Provider, error) {
	cfg, _ := config.LoadAppConfig()
	// Standard resolution with global context
	opts := resolver.ResolveOptions{
		FilePath:             filePath,
		RekordboxPrimaryPath: cfg.Rekordbox.PrimaryFilePath,
		Apply:                apply,
		Verbose:              verbose,
		Feedback:             &TerminalFeedback{},
	}
	return resolver.ResolveSelection(locStr, queryOverride, opts)
}

func stringsTitle(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
