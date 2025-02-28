package usecase

import (
	"bytes"
	"log/slog"

	"github.com/sHelllWalker/heimdallr/internal/enums"
	"github.com/sHelllWalker/heimdallr/internal/service"
	"github.com/sHelllWalker/heimdallr/internal/types"
)

type BroadcastEvent struct {
	logger    *slog.Logger
	mailer    *service.Mailer
	templater *service.Templater
}

func NewBroadcastEvent(logger *slog.Logger, mailer *service.Mailer, templater *service.Templater) *BroadcastEvent {
	return &BroadcastEvent{
		logger:    logger,
		mailer:    mailer,
		templater: templater,
	}
}

func (b BroadcastEvent) Do(
	event types.RenderableEvent,
	opts types.ProviderOptions,
	mOpts types.MessageOptions,
) {
	for _, messenger := range enums.Messengers() {
		t := b.templater.GetTemplate(messenger, event.GetResource())

		var message bytes.Buffer

		err := t.Execute(&message, event)
		if err != nil {
			b.logger.Error("template executing fail", slog.Any("messenger", messenger), slog.Any("error", err))
			continue
		}

		b.mailer.SendMessage(message.String(), messenger, opts, mOpts)
	}
}
