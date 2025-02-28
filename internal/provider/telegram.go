package provider

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/sHelllWalker/heimdallr/internal/config"
	"github.com/sHelllWalker/heimdallr/internal/types"
)

const tgAPIPathPattern = "https://api.telegram.org/bot%s/%s"

type Telegram struct {
	client *resty.Client
	conf   config.TelegramConfig
}

func NewTelegram(client *resty.Client, conf config.TelegramConfig) *Telegram {
	if conf.Token == "" {
		return nil
	}

	t := Telegram{
		client: client,
		conf:   conf,
	}
	if !t.getMe() {
		return nil
	}

	return &t
}

func (t *Telegram) Send(text string, opts types.ProviderOptions, _ types.MessageOptions) (isSend bool, err error) {
	if opts.GetChatID() == "" {
		return false, nil
	}

	body := struct {
		Text                string `json:"text"`
		ChatID              string `json:"chat_id"`
		ParseMode           string `json:"parse_mode,omitempty"`
		DisableNotification bool   `json:"disable_notification,omitempty"`
	}{
		Text:                text,
		ChatID:              opts.GetChatID(),
		ParseMode:           t.conf.ParseMode,
		DisableNotification: t.conf.SilentMode,
	}

	result := struct {
		Ok          bool   `json:"ok"`
		Description string `json:"description"`
		Result      struct {
			MessageID int `json:"message_id"`
		} `json:"result"`
	}{}

	resp, err := t.client.R().
		SetBody(body).
		SetContentLength(true).
		SetResult(&result).
		Post(fmt.Sprintf(tgAPIPathPattern, t.conf.Token, "sendMessage"))
	if err != nil {
		return false, err
	}
	if resp.StatusCode() != http.StatusOK || !result.Ok {
		return false, fmt.Errorf("fail request, description: `%s`", result.Description)
	}

	return true, nil
}

func (t *Telegram) getMe() bool {
	var result struct {
		Ok bool `json:"ok"`
	}

	resp, err := t.client.R().
		SetResult(&result).
		Post(fmt.Sprintf(tgAPIPathPattern, t.conf.Token, "getMe"))

	if err != nil || resp.StatusCode() != http.StatusOK {
		return false
	}

	return result.Ok
}
