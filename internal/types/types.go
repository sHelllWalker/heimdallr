package types

import "github.com/sHelllWalker/heimdallr/internal/enums"

type Provider interface {
	Send(text string, opts ProviderOptions, mOpts MessageOptions) (isSend bool, err error)
}

type ProviderOptions interface {
	GetChannel() string
	GetChatID() string
}

type MessageOptions interface {
	GetColor() string
}

type RenderableEvent interface {
	GetResource() enums.Event
	SetResource(event enums.Event)
	GetTitle() string
	GetDescription() string
	GetLevel() string
	GetLink() string

	GetEventLevel() enums.EventLevel
}
