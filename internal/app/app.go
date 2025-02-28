package app

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sHelllWalker/heimdallr/internal/config"
	"github.com/sHelllWalker/heimdallr/internal/enums"
	"github.com/sHelllWalker/heimdallr/internal/provider"
	"github.com/sHelllWalker/heimdallr/internal/service"
	"github.com/sHelllWalker/heimdallr/internal/types"
)

type App struct {
	conf *config.Config

	logger  *slog.Logger
	logFile *os.File

	mailer    *service.Mailer
	templater *service.Templater
}

func NewApp() (*App, error) {
	app := &App{}
	conf, err := config.InitConfig()
	if err != nil {
		log.Fatalf("can`t init app config, error: %s", err)
	}
	app.conf = conf

	app.initLogger(conf.LoggerConfig)

	httpClient := resty.New()
	httpClient.SetRetryCount(3)

	app.mailer = service.NewMailer(app.logger, map[enums.Messenger]types.Provider{
		enums.MatterMost: provider.NewMatterMost(httpClient, app.conf.MatterMostConfig),
		enums.Telegram:   provider.NewTelegram(httpClient, app.conf.TelegramConfig),
	})
	app.templater = service.NewTemplater(app.logger, app.conf.TemplateConfig)

	return app, nil
}

func (a *App) Listen() error {
	addr := ":" + strconv.Itoa(a.conf.AppPort)

	mux := a.initRouter()

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	var serveErr error
	go func() {
		a.logger.Info("start heimdallr serve", slog.Int("port", a.conf.AppPort))
		serveErr = server.ListenAndServe()
		quit <- syscall.SIGTERM
	}()

	<-quit

	if !errors.Is(serveErr, http.ErrServerClosed) {
		a.logger.Error("serve error", slog.Any("error", serveErr))
	}

	a.logger.Info("shutting down heimdallr")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		a.logger.Error("shutting down failed", "error", err)

		return err
	}

	a.logger.Info("gracefully stopped")

	return nil
}

func (a *App) Close() error {
	if a.logFile != nil {
		err := a.logFile.Close()

		return err
	}

	return nil
}

func (a *App) initLogger(conf config.LoggerConfig) {
	opts := &slog.HandlerOptions{
		Level:     getLogLevel(conf.LogLevel),
		AddSource: conf.AddLogSource,
	}

	logFile, err := a.getLoggerFile(conf.LogFilePath)
	logger := slog.New(slog.NewJSONHandler(logFile, opts))
	if err != nil {
		logger.Error("log init error", slog.Any("error", err))
	}

	a.logger = logger
}

func (a *App) getLoggerFile(filePath string) (*os.File, error) {
	if filePath != "" {
		file, err := os.OpenFile(
			filePath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0666,
		)
		if err != nil {
			return os.Stdout, err
		}
		a.logFile = file

		return file, err
	}

	return os.Stdout, nil
}

func getLogLevel(logLevel string) slog.Level {
	switch logLevel {
	case "ERROR":
		return slog.LevelError
	case "WARN":
		return slog.LevelWarn
	case "INFO":
		return slog.LevelInfo
	case "DEBUG":
		fallthrough
	default:
		return slog.LevelDebug
	}
}
