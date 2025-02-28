package provider

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/sHelllWalker/heimdallr/internal/config"
	"github.com/sHelllWalker/heimdallr/internal/types"
)

type MatterMost struct {
	client *resty.Client
	conf   config.MatterMostConfig
}

func NewMatterMost(client *resty.Client, config config.MatterMostConfig) *MatterMost {
	if config.WebhookURL == "" {
		return nil
	}

	return &MatterMost{
		client: client,
		conf:   config,
	}
}

func (m *MatterMost) Send(
	text string,
	opts types.ProviderOptions,
	mOpts types.MessageOptions,
) (isSend bool, err error) {
	if opts.GetChannel() == "" {
		return false, nil
	}

	body := struct {
		Channel     string            `json:"channel"`
		Text        string            `json:"text,omitempty"`
		Username    string            `json:"username,omitempty"`
		IconURL     string            `json:"icon_url,omitempty"`
		Attachments map[string]string `json:"attachments,omitempty"`
	}{
		Channel:  opts.GetChannel(),
		Username: m.conf.Username,
		IconURL:  m.conf.IconURL,
	}
	if m.conf.WithAttachments {
		body.Attachments["text"] = text
		body.Attachments["color"] = mOpts.GetColor()
	} else {
		body.Text = text
	}

	result := ""

	resp, err := m.client.R().
		SetBody(body).
		SetContentLength(true).
		SetResult(&result).
		Post(tgAPIPathPattern)
	if err != nil {
		return false, err
	}
	if resp.StatusCode() != http.StatusOK || result != "ok" {
		return false, fmt.Errorf("can`t send request, get message: `%s`", result)
	}

	return true, nil
}
