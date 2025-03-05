package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/sHelllWalker/heimdallr/internal/api/dto"
	"github.com/sHelllWalker/heimdallr/internal/domain/models/broadcast"
	"github.com/sHelllWalker/heimdallr/internal/enums"
	"github.com/sHelllWalker/heimdallr/internal/types"
	"github.com/sHelllWalker/heimdallr/internal/usecase"
)

type SendEvent struct {
	logger  *slog.Logger
	usecase *usecase.BroadcastEvent
}

func NewSendEvent(logger *slog.Logger, uc *usecase.BroadcastEvent) *SendEvent {
	return &SendEvent{
		logger:  logger,
		usecase: uc,
	}
}

func (r *SendEvent) Handle(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		r.writeError(writer)

		return
	}

	r.logger.Info("event received")

	queryParams := request.URL.Query()
	options := broadcast.MessangerOptions{
		Channel: queryParams.Get("channel"),
		ChatID:  queryParams.Get("chatId"),
	}

	event, err := encodeEvent(request)
	if err != nil {
		r.logger.Error("bad request", slog.Any("error", err))
		r.writeError(writer)

		return
	}

	mOptions := broadcast.MessageOptions{
		Color: matchEventLevelToColor(event.GetEventLevel()),
	}

	go r.usecase.Do(event, &options, &mOptions)

	r.writeSuccess(writer)
}

func encodeEvent(request *http.Request) (types.RenderableEvent, error) {
	resource := request.Header.Get("Sentry-Hook-Resource")

	rawBody, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	switch resource {
	case enums.Installation:
		return unmarshalEvent[*dto.Installation](rawBody, resource)
	case enums.IssueAlert:
		return unmarshalEvent[*dto.IssueAlert](rawBody, resource)
	case enums.MetricAlert:
		return unmarshalEvent[*dto.MetricAlert](rawBody, resource)
	case enums.Issues:
		return unmarshalEvent[*dto.Issues](rawBody, resource)
	case enums.Comments:
		return unmarshalEvent[*dto.Comment](rawBody, resource)
	case enums.Errors:
		fallthrough
	default:
		return unmarshalEvent[*dto.Error](rawBody, resource)
	}
}

func unmarshalEvent[T types.RenderableEvent](body []byte, resource string) (T, error) {
	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return result, err
	}

	result.SetResource(enums.Event(resource))

	return result, nil
}

func (r *SendEvent) writeError(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)

	err := json.NewEncoder(writer).Encode(map[string]any{
		"success": false,
	})
	if err != nil {
		r.logger.Error("json encode of error response failed", "error", err)
	}
}

func (r *SendEvent) writeSuccess(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	err := json.NewEncoder(writer).Encode(map[string]any{
		"success": true,
	})
	if err != nil {
		r.logger.Error("json encode of error response failed", "error", err)
	}
}

func matchEventLevelToColor(lvl enums.EventLevel) string {
	switch lvl {
	case enums.ProblemResolved:
		return "#008000"
	case enums.Problem:
		return "#FF0000"
	case enums.Notification:
		return "#808080"
	}

	return ""
}
