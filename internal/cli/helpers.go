package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/djerr"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/resolver"
)

// HandleError provides user-friendly messages for sentinel provider errors.
func HandleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, djerr.ErrReadOnly) {
		return fmt.Errorf("operation failed: this provider is read-only")
	}
	if errors.Is(err, djerr.ErrUnsupportedResource) {
		return fmt.Errorf("operation failed: this resource type is not supported by the provider")
	}
	if errors.Is(err, djerr.ErrInvalidParent) {
		return fmt.Errorf("operation failed: cannot create the resource in that location (structural constraint)")
	}

	return err
}

func getExecContext() provider.ExecutionContext {
	return provider.ExecutionContext{
		Apply:   apply,
		Verbose: verbose,
	}
}

func ResolveSelection(locStr string, queryOverride string) (*resolver.Selection, error) {
	// Standard resolution with global context
	opts := resolver.ResolveOptions{
		FilePath: filePath,
		Apply:    apply,
		Verbose:  verbose,
	}
	return resolver.ResolveSelection(locStr, queryOverride, opts)
}

func stringsTitle(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
