package dto

import (
	"fmt"
	"strings"

	"github.com/sHelllWalker/heimdallr/internal/enums"
)

type Error struct {
	Event
}

func (e *Error) GetTitle() string {
	err, _ := e.Data["error"].(map[string]any)
	title, ok := err["title"].(string)
	if !ok {
		return ""
	}

	return title
}

func (e *Error) GetDescription() string {
	err, _ := e.Data["error"].(map[string]any)
	metadata, ok := err["metadata"].(map[string]any)
	if !ok || len(metadata) == 0 {
		return ""
	}
	value, _ := metadata["value"].(string)
	errType, _ := metadata["type"].(string)
	filename, _ := metadata["filename"].(string)

	description := ""
	if filename != "" {
		description = fmt.Sprintf("file `%s`: ", filename)
	}

	return description + errType + " " + value
}

func (e *Error) GetLevel() string {
	err, _ := e.Data["error"].(map[string]any)
	level, ok := err["level"].(string)
	if !ok {
		return ""
	}

	return strings.ToUpper(level)
}

func (e *Error) GetLink() string {
	err, _ := e.Data["error"].(map[string]any)
	webURL, ok := err["web_url"].(string)
	if !ok {
		return ""
	}

	return webURL
}

func (e *Error) GetReadableResource() string {
	return "Error"
}

func (e *Error) GetEventLevel() enums.EventLevel {
	return enums.Problem
}
