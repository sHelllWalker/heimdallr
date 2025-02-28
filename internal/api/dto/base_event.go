package dto

import "github.com/sHelllWalker/heimdallr/internal/enums"

type Event struct {
	Resource enums.Event

	Action       string         `json:"action"`
	Data         map[string]any `json:"data"`
	Actor        map[string]any `json:"actor"`
	Installation struct {
		UUID string `json:"uuid"`
	} `json:"installation"`
}

func (e *Event) GetResource() enums.Event {
	return e.Resource
}

func (e *Event) SetResource(eventType enums.Event) {
	e.Resource = eventType
}
