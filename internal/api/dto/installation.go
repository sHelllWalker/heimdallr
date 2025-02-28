package dto

import (
	"fmt"

	"github.com/sHelllWalker/heimdallr/internal/enums"
)

type Installation struct {
	Event
}

func (i *Installation) GetTitle() string {
	actorName, _ := i.Actor["name"].(string)
	if actorName != "" {
		return "by " + actorName
	}

	return ""
}

func (i *Installation) GetDescription() string {
	installation, _ := i.Data["installation"].(map[string]any)
	status, _ := installation["status"].(string)
	app, _ := installation["app"].(map[string]any)
	slug, _ := app["slug"].(string)

	if slug == "" {
		return "unknown app"
	}

	return fmt.Sprintf("%s app `%s` status `%s`", i.Action, slug, status)
}

func (i *Installation) GetLevel() string {
	return ""
}

func (i *Installation) GetLink() string {
	return ""
}

func (i *Installation) GetReadableResource() string {
	return "Installation"
}

func (i *Installation) GetEventLevel() enums.EventLevel {
	return enums.Notification
}
