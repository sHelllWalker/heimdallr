package dto

import (
	"fmt"
	"strings"

	"github.com/sHelllWalker/heimdallr/internal/enums"
)

type Issues struct {
	Event
}

func (i *Issues) GetTitle() string {
	issue, _ := i.Data["issue"].(map[string]any)
	title, _ := issue["title"].(string)

	return fmt.Sprintf("[%s] %s", i.Action, title)
}

func (i *Issues) GetDescription() string {
	issue, _ := i.Data["issue"].(map[string]any)
	metadata, ok := issue["metadata"].(map[string]any)
	if !ok || len(metadata) == 0 {
		return ""
	}
	value, _ := metadata["value"].(string)
	eventType, _ := metadata["type"].(string)

	description := eventType + " " + value
	if filename, _ := metadata["filename"].(string); filename != "" {
		description = fmt.Sprintf("\nfile `%s` ", filename)
	}

	if status, _ := issue["status"].(string); status != "" {
		description = fmt.Sprintf("\nstatus `%s` \n", status)
	}

	return description
}

func (i *Issues) GetLevel() string {
	issue, _ := i.Data["issue"].(map[string]any)
	level, ok := issue["level"].(string)
	if !ok {
		return ""
	}

	return strings.ToUpper(level)
}

func (i *Issues) GetLink() string {
	issue, _ := i.Data["issue"].(map[string]any)
	webURL, ok := issue["web_url"].(string)
	if !ok {
		return ""
	}

	return webURL
}

func (i *Issues) GetReadableResource() string {
	return "Issue"
}

func (i *Issues) GetEventLevel() enums.EventLevel {
	switch i.Action {
	case "resolved", "archived":
		return enums.ProblemResolved
	case "unresolved", "created":
		return enums.Problem
	case "assigned":
		return enums.Notification
	}

	return enums.Notification
}
