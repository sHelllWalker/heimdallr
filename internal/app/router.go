package app

import (
	"net/http"

	"github.com/sHelllWalker/heimdallr/internal/api/handler"
	"github.com/sHelllWalker/heimdallr/internal/usecase"
)

func (a *App) initRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"/broadcast",
		handler.NewSendEvent(a.logger, usecase.NewBroadcastEvent(a.logger, a.mailer, a.templater)).Handle,
	)

	return mux
}
