package enums

type Event string

const (
	Installation = "installation"
	IssueAlert   = "event_alert"
	MetricAlert  = "metric_alert"
	Issues       = "issue"
	Comments     = "comment"
	Errors       = "error"
)
