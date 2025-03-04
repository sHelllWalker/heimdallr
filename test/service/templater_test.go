package service

import (
	"io"
	"log/slog"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/sHelllWalker/heimdallr/internal/config"
	"github.com/sHelllWalker/heimdallr/internal/enums"
	"github.com/sHelllWalker/heimdallr/internal/service"
	"github.com/stretchr/testify/require"
)

func TestTemplater_GetTemplate_Default(t *testing.T) {
	templater := service.NewTemplater(slog.New(slog.NewTextHandler(io.Discard, nil)), config.TemplateConfig{})

	testAllTemplatesDefault(t, *templater)
}
func TestTemplater_GetTemplate_BadTemplate(t *testing.T) {
	mmErrPath, err := filepath.Abs("../resource/templates/broaken_template.txt")
	require.NoError(t, err)
	templater := service.NewTemplater(slog.New(slog.NewTextHandler(io.Discard, nil)), config.TemplateConfig{
		MMErrorsTemplatePath: mmErrPath,
	})

	testAllTemplatesDefault(t, *templater)
}

func TestTemplater_GetTemplate_BadPath(t *testing.T) {
	templater := service.NewTemplater(slog.New(slog.NewTextHandler(io.Discard, nil)), config.TemplateConfig{
		MMErrorsTemplatePath: "/",
	})

	testAllTemplatesDefault(t, *templater)
}

func TestTemplater_GetTemplate_CustomTemplate(t *testing.T) {
	mmErrPath, err := filepath.Abs("../resource/templates/issue_alert.txt")
	require.NoError(t, err)
	tgErrPath, err := filepath.Abs("../resource/templates/issue_alert.txt")
	require.NoError(t, err)

	templater := service.NewTemplater(slog.New(slog.NewTextHandler(io.Discard, nil)), config.TemplateConfig{
		MMErrorsTemplatePath: mmErrPath,
		TgErrorsTemplatePath: tgErrPath,
	})

	mmErrorTemplate := templater.GetTemplate(enums.MatterMost, enums.Errors)
	require.NotNil(t, mmErrorTemplate)
	tgErrorTemplate := templater.GetTemplate(enums.Telegram, enums.Errors)
	require.NotNil(t, tgErrorTemplate)

	require.NotSame(t, mmErrorTemplate, tgErrorTemplate)

	userErrorTemplates := map[enums.Messenger]*template.Template{
		enums.MatterMost: mmErrorTemplate,
		enums.Telegram:   tgErrorTemplate,
	}

	for _, m := range enums.Messengers() {
		mTemplates := make([]*template.Template, 0)
		for _, e := range enums.Events() {
			templ := templater.GetTemplate(m, e)
			require.NotNil(t, templ)
			if e == enums.Errors {
				require.Same(t, userErrorTemplates[m], templ)

				continue
			}

			mTemplates = append(mTemplates, templ)
		}

		firstTemplate := mTemplates[0]
		for _, templ := range mTemplates {
			require.Same(t, firstTemplate, templ)
			require.NotSame(t, userErrorTemplates[m], templ)
		}
	}
}

func testAllTemplatesDefault(t *testing.T, templater service.Templater) {
	t.Helper()

	seenTemplates := make(map[*template.Template]bool)
	for _, m := range enums.Messengers() {
		mTemplates := make([]*template.Template, 0)
		for _, e := range enums.Events() {
			templ := templater.GetTemplate(m, e)
			require.NotNil(t, templ)

			mTemplates = append(mTemplates, templ)
		}

		firstTemplate := mTemplates[0]
		for _, templ := range mTemplates {
			require.Same(t, firstTemplate, templ)
		}

		_, found := seenTemplates[firstTemplate]
		require.False(t, found)

		seenTemplates[firstTemplate] = true
	}
}
