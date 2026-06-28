package cli

import (
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type CLIProgressListener struct {
	p        *mpb.Progress
	totalBar *mpb.Bar
	trackBar *mpb.Bar
}

func NewCLIProgressListener() *CLIProgressListener {
	return &CLIProgressListener{
		p: mpb.New(mpb.WithWidth(64)),
	}
}

func (l *CLIProgressListener) OnStart(total int64) {
	l.totalBar = l.p.AddBar(total,
		mpb.PrependDecorators(
			decor.Name("Overall Sync", decor.WCSyncSpaceR),
			decor.CountersNoUnit("%d / %d", decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WCSyncSpace),
			decor.OnComplete(decor.EwmaETA(decor.ET_STYLE_GO, 60), "done!"),
		),
	)
}

func (l *CLIProgressListener) OnTrackStart(trackTitle string) {
	displayName := trackTitle
	if len(displayName) > 20 {
		displayName = displayName[:17] + "..."
	}
	l.trackBar = l.p.AddBar(1,
		mpb.BarRemoveOnComplete(),
		mpb.PrependDecorators(
			decor.Name(fmt.Sprintf("  -> %s", displayName), decor.WCSyncSpaceR),
		),
	)
}

func (l *CLIProgressListener) OnTrackEnd() {
	if l.trackBar != nil {
		l.trackBar.Increment()
	}
	if l.totalBar != nil {
		l.totalBar.Increment()
	}
}

func (l *CLIProgressListener) OnComplete() {
	l.p.Wait()
}
