package service

import (
	"errors"
	"log/slog"

	"github.com/sHelllWalker/heimdallr/internal/enums"
	"github.com/sHelllWalker/heimdallr/internal/types"
)

var ErrUndefinedMessenger = errors.New("undefined messenger")

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
) error {
	provider, found := m.providerMap[messenger]
	if !found {
		return ErrUndefinedMessenger
	}

	isSend, err := provider.Send(message, messengerOpts, messageOpts)
	if err != nil {
		return err
	}

	if !isSend {
		m.logger.Debug("message not sended", slog.Any("messenger", messenger))

		return nil
	}

	m.logger.Debug("message sending success", slog.Any("messenger", messenger))

	return nil
}
