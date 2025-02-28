package service

import (
	"fmt"
	"log/slog"
	"strings"
	"text/template"

	"github.com/sHelllWalker/heimdallr/internal/config"
	"github.com/sHelllWalker/heimdallr/internal/enums"
	"github.com/sHelllWalker/heimdallr/web/templates"
)

type Templater struct {
	logger           *slog.Logger
	markdownReplacer *strings.Replacer

	templateMap map[templateKey]*template.Template
}

type templateKey struct {
	Messenger enums.Messenger
	Event     enums.Event
}

func NewTemplater(logger *slog.Logger, config config.TemplateConfig) *Templater {
	messengers := enums.Messengers()

	templater := &Templater{
		logger:           logger,
		markdownReplacer: markDownEscaper(),
		templateMap:      make(map[templateKey]*template.Template, len(messengers)),
	}

	for _, messenger := range messengers {
		templater.createTemplate(config, messenger)
	}

	return templater
}

func (t *Templater) GetTemplate(messenger enums.Messenger, event enums.Event) *template.Template {
	templateKey := templateKey{Messenger: messenger, Event: event}

	if t, found := t.templateMap[templateKey]; found {
		return t
	}

	return t.getDefault(messenger)
}

func (t *Templater) createTemplate(config config.TemplateConfig, m enums.Messenger) {
	t.initDefaulTemplate(m)

	userTemplatePaths := getUserPaths(config, m)
	for event, templatePath := range userTemplatePaths {
		if templatePath == "" {
			continue
		}

		var err error
		userTemplate := template.New(templatePath)
		t.registerFunctions(userTemplate)
		userTemplate, err = userTemplate.ParseFiles(templatePath)
		if err != nil {
			t.logger.Error(
				fmt.Sprintf("can`t create user template by path `%s`", templatePath),
				slog.Any("error", err),
			)

			continue
		}

		t.templateMap[templateKey{Messenger: m, Event: event}] = userTemplate
	}
}

func (t *Templater) initDefaulTemplate(m enums.Messenger) {
	tmplt := template.New(string(m))
	t.registerFunctions(tmplt)

	template.Must(tmplt.ParseFS(templates.Templates, string(m)))

	t.templateMap[templateKey{Messenger: m, Event: "default"}] = tmplt
}

func (t *Templater) getDefault(m enums.Messenger) *template.Template {
	return t.templateMap[templateKey{Messenger: m, Event: "default"}]
}

func (t *Templater) registerFunctions(templ *template.Template) {
	templ.Funcs(template.FuncMap{
		"escapeMarkdown": func(s string) string {
			return t.markdownReplacer.Replace(s)
		},
	})
}

func getUserPaths(config config.TemplateConfig, m enums.Messenger) map[enums.Event]string {
	switch m {
	case enums.MatterMost:
		return map[enums.Event]string{
			enums.Installation: config.MMInstallationTemplatePath,
			enums.IssueAlert:   config.MMIssueAlertTemplatePath,
			enums.MetricAlert:  config.MMMetricAlertTemplatePath,
			enums.Issues:       config.MMIssuesTemplatePath,
			enums.Comments:     config.MMCommentsTemplatePath,
			enums.Errors:       config.MMErrorsTemplatePath,
		}
	case enums.Telegram:
		return map[enums.Event]string{
			enums.Installation: config.TgInstallationTemplatePath,
			enums.IssueAlert:   config.TgIssueAlertTemplatePath,
			enums.MetricAlert:  config.TgMetricAlertTemplatePath,
			enums.Issues:       config.TgIssuesTemplatePath,
			enums.Comments:     config.TgCommentsTemplatePath,
			enums.Errors:       config.TgErrorsTemplatePath,
		}
	default:
		return nil
	}
}

func markDownEscaper() *strings.Replacer {
	return strings.NewReplacer(
		"\\", "\\\\",
		"`", "\\`",
		"*", "\\*",
		"_", "\\_",
		"{", "\\{",
		"}", "\\}",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		".", "\\.",
		"!", "\\!",
		"|", "\\|",
		"~", "\\~",
		":", "\\:",
		">", "\\>",
		"=", "\\=",
	)
}
