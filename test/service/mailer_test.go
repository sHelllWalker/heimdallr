package service

import (
	"io"
	"log/slog"
	"testing"

	"github.com/sHelllWalker/heimdallr/internal/domain/models/broadcast"
	"github.com/sHelllWalker/heimdallr/internal/enums"
	"github.com/sHelllWalker/heimdallr/internal/service"
	"github.com/sHelllWalker/heimdallr/internal/types"
	"github.com/stretchr/testify/require"
)

type mockProvider struct {
	sendFunc func(text string, opts types.ProviderOptions, mOpts types.MessageOptions) (bool, error)
}

func (m *mockProvider) Send(text string, opts types.ProviderOptions, mOpts types.MessageOptions) (bool, error) {
	return m.sendFunc(text, opts, mOpts)
}

func TestMailer_SendMessage(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name        string
		messenger   enums.Messenger
		provider    *mockProvider
		expectedErr error
	}{
		{
			name:      "success send",
			messenger: enums.Telegram,
			provider: &mockProvider{
				sendFunc: func(_ string, _ types.ProviderOptions, _ types.MessageOptions) (bool, error) {
					return true, nil
				},
			},
		},
		{
			name:        "provider not found",
			messenger:   "unknown",
			provider:    nil,
			expectedErr: service.ErrUndefinedMessenger,
		},
		{
			name:      "failure send",
			messenger: enums.Telegram,
			provider: &mockProvider{
				sendFunc: func(_ string, _ types.ProviderOptions, _ types.MessageOptions) (bool, error) {
					return false, nil
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			providerMap := make(map[enums.Messenger]types.Provider)
			if tt.provider != nil {
				providerMap[enums.Telegram] = tt.provider
			}

			mailer := service.NewMailer(logger, providerMap)

			opts := &broadcast.MessangerOptions{
				Channel: "test",
				ChatID:  "123",
			}
			mOpts := &broadcast.MessageOptions{
				Color: "red",
			}

			err := mailer.SendMessage("test message", tt.messenger, opts, mOpts)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			}

			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
