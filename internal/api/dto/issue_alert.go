package dto

import (
	"fmt"
	"strings"

	"github.com/sHelllWalker/heimdallr/internal/enums"
)

type IssueAlert struct {
	Event
}

func (ia *IssueAlert) GetTitle() string {
	alert, _ := ia.Data["issue_alert"].(map[string]any)
	title, ok := alert["title"].(string)
	if ok {
		return fmt.Sprintf("[%s] %s", strings.ToUpper(ia.Action), title)
	}

	return ""
}

func (ia *IssueAlert) GetDescription() string {
	event, _ := ia.Data["event"].(map[string]any)
	metadata, ok := event["metadata"].(map[string]any)
	if !ok || len(metadata) == 0 {
		return ""
	}
	value, _ := metadata["value"].(string)
	eventType, _ := metadata["type"].(string)
	filename, _ := metadata["filename"].(string)

	description := ""
	if filename != "" {
		description = fmt.Sprintf("file `%s`: ", filename)
	}

	return description + eventType + " " + value
}

func (ia *IssueAlert) GetLevel() string {
	err, _ := ia.Data["event"].(map[string]any)
	level, ok := err["level"].(string)
	if !ok {
		return ""
	}

	return strings.ToUpper(level)
}

func (ia *IssueAlert) GetLink() string {
	event, _ := ia.Data["event"].(map[string]any)
	webURL, ok := event["web_url"].(string)
	if !ok {
		return ""
	}

	return webURL
}

func (ia *IssueAlert) GetReadableResource() string {
	return "Issue Alert"
}

func (ia *IssueAlert) GetEventLevel() enums.EventLevel {
	return enums.Problem
}
