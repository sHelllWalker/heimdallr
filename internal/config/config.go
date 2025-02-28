package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	AppPort int `env:"APP_PORT" envDefault:"8666" validate:"min=0,max=65535"`

	LoggerConfig

	TemplateConfig

	TelegramConfig
	MatterMostConfig
}

type LoggerConfig struct {
	LogLevel     string `env:"LOG_LEVEL" envDefault:"DEBUG"`
	LogFilePath  string `env:"LOG_FILE_PATH"` // stdout default
	AddLogSource bool   `env:"ADD_LOG_SOURCE"`
}

type TelegramConfig struct {
	Token      string `env:"TOKEN"`
	ParseMode  string `env:"PARSE_MODE" envDefault:"MarkdownV2" validate:"oneof=MarkdownV2 Markdown HTML"`
	SilentMode bool   `env:"SILENT_MODE"` // send message without notification
}

type MatterMostConfig struct {
	WebhookURL      string `env:"WEBHOOK_URL"`
	IconURL         string `env:"ICON_URL"`
	Username        string `env:"USERNAME"`
	WithAttachments bool   `env:"WITH_ATTACHMENTS"`
}

type TemplateConfig struct {
	MMInstallationTemplatePath string `env:"MM_INSTALLATION_TEMPLATE_PATH"`
	MMIssueAlertTemplatePath   string `env:"MM_ISSUE_ALERT_TEMPLATE_PATH"`
	MMMetricAlertTemplatePath  string `env:"MM_METRIC_ALERT_TEMPLATE_PATH"`
	MMIssuesTemplatePath       string `env:"MM_ISSUES_TEMPLATE_PATH"`
	MMCommentsTemplatePath     string `env:"MM_COMMENTS_TEMPLATE_PATH"`
	MMErrorsTemplatePath       string `env:"MM_ERRORS_TEMPLATE_PATH"`

	TgInstallationTemplatePath string `env:"TG_INSTALLATION_TEMPLATE_PATH"`
	TgIssueAlertTemplatePath   string `env:"TG_ISSUE_ALERT_TEMPLATE_PATH"`
	TgMetricAlertTemplatePath  string `env:"TG_METRIC_ALERT_TEMPLATE_PATH"`
	TgIssuesTemplatePath       string `env:"TG_ISSUES_TEMPLATE_PATH"`
	TgCommentsTemplatePath     string `env:"TG_COMMENTS_TEMPLATE_PATH"`
	TgErrorsTemplatePath       string `env:"TG_ERRORS_TEMPLATE_PATH"`
}

func InitConfig() (*Config, error) {
	conf, err := env.ParseAs[Config]()
	if err != nil {
		return &conf, err
	}

	return &conf, validator.New().Struct(conf)
}
