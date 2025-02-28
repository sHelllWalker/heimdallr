package enums

type EventLevel int

const (
	ProblemResolved EventLevel = iota
	Notification
	Problem
)
