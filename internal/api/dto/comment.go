package dto

import (
	"fmt"

	"github.com/sHelllWalker/heimdallr/internal/enums"
)

type Comment struct {
	Event
}

func (c *Comment) GetTitle() string {
	actorName, _ := c.Actor["name"].(string)

	return fmt.Sprintf("%s by %s", c.Action, actorName)
}

func (c *Comment) GetDescription() string {
	comment, _ := c.Data["comment"].(string)
	if projectSlug, _ := c.Data["project_slug"].(string); projectSlug != "" {
		return fmt.Sprintf("%s\nproject: %s", comment, projectSlug)
	}

	return comment
}

func (c *Comment) GetLevel() string {
	return ""
}

func (c *Comment) GetLink() string {
	return ""
}

func (c *Comment) GetReadableResource() string {
	return "Comment"
}

func (c *Comment) GetEventLevel() enums.EventLevel {
	return enums.Notification
}
