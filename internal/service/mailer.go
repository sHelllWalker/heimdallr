package service

import (
	"log/slog"

	"github.com/sHelllWalker/heimdallr/internal/enums"
	"github.com/sHelllWalker/heimdallr/internal/types"
)

const sendMessageError = "message sending failure"

type Mailer struct {
	logger      *slog.Logger
	providerMap map[enums.Messenger]types.Provider
}

func NewMailer(logger *slog.Logger, providerMap map[enums.Messenger]types.Provider) *Mailer {
	return &Mailer{
		logger:      logger,
		providerMap: providerMap,
	}
}

func (m *Mailer) SendMessage(
	message string,
	messenger enums.Messenger,
	messengerOpts types.ProviderOptions,
	messageOpts types.MessageOptions,
) {
	provider, found := m.providerMap[messenger]
	if !found {
		m.logger.Error(
			sendMessageError,
			slog.Any("messenger", messenger),
			slog.String("reason", "provider not found"),
		)

		return
	}

	isSend, err := provider.Send(message, messengerOpts, messageOpts)
	if err != nil {
		m.logger.Error(
			sendMessageError,
			slog.Any("messenger", messenger),
			slog.Any("reason", err),
		)

		return
	}

	if !isSend {
		m.logger.Debug("message not sended", slog.Any("messenger", messenger))

		return
	}

	m.logger.Debug("message sending success", slog.Any("messenger", messenger))
}
