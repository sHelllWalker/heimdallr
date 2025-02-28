package provider

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/sHelllWalker/heimdallr/internal/config"
	"github.com/sHelllWalker/heimdallr/internal/enums"
	"github.com/sHelllWalker/heimdallr/internal/types"
)

func CreateProvider(m enums.Messenger, conf config.Config, client *resty.Client) (types.Provider, error) {
	switch m {
	case enums.MatterMost:
		return NewMatterMost(client, conf.MatterMostConfig), nil
	case enums.Telegram:
		return NewTelegram(client, conf.TelegramConfig), nil

	default:
		return nil, fmt.Errorf("no provider for messenger `%s`", m)
	}
}
