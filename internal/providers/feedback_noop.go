package provider

type NoopFeedback struct{}

func (f NoopFeedback) OnPreview(message string)                  {}
func (f NoopFeedback) OnSuccess(message string)                  {}
func (f NoopFeedback) OnWarning(message string)                  {}
func (f NoopFeedback) OnStatus(message string)                   {}
func (f NoopFeedback) OnProgress(done, total int)                {}
func (f NoopFeedback) OnTable(headers []string, rows [][]string) {}
