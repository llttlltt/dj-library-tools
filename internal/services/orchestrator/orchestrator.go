package orchestrator

import (
	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/llttlltt/dj-library-tools/internal/services/resolver"
)

type Orchestrator struct {
	Feedback             provider.Feedback
	RekordboxPrimaryPath string
}

type Options struct {
	RekordboxPrimaryPath string
}

func New(fb provider.Feedback, opts Options) *Orchestrator {
	if fb == nil {
		fb = provider.NoopFeedback{}
	}
	return &Orchestrator{
		Feedback:             fb,
		RekordboxPrimaryPath: opts.RekordboxPrimaryPath,
	}
}

type RunOptions struct {
	FilePath string
	Apply    bool
	Verbose  bool
}

func (o *Orchestrator) buildResolveOptions(opts RunOptions) resolver.ResolveOptions {
	return resolver.ResolveOptions{
		FilePath:             opts.FilePath,
		RekordboxPrimaryPath: o.RekordboxPrimaryPath,
		Apply:                opts.Apply,
		Verbose:              opts.Verbose,
		Feedback:             o.Feedback,
	}
}

func (o *Orchestrator) buildExecContext(opts RunOptions) provider.ExecutionContext {
	return provider.ExecutionContext{
		Apply:    opts.Apply,
		Verbose:  opts.Verbose,
		Feedback: o.Feedback,
	}
}

type ListResult struct {
	Tracks []models.Track
	Groups []models.ResourceGroup
	Provider provider.Provider
}

func (o *Orchestrator) List(locStr string, queryOverride string, opts RunOptions) (*ListResult, error) {
	sel, err := resolver.ResolveSelection(locStr, queryOverride, o.buildResolveOptions(opts))
	if err != nil {
		return nil, err
	}

	return &ListResult{
		Tracks:   sel.Tracks,
		Groups:   sel.Groups,
		Provider: sel.Provider,
	}, nil
}
