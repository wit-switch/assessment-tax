package rest

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wit-switch/assessment-tax/config"
	_ "github.com/wit-switch/assessment-tax/docs"
	"github.com/wit-switch/assessment-tax/infrastructure"
	"github.com/wit-switch/assessment-tax/internal/core/service"
	"github.com/wit-switch/assessment-tax/internal/handler"
	httphdl "github.com/wit-switch/assessment-tax/internal/handler/http"
	"github.com/wit-switch/assessment-tax/internal/handler/middleware"
	"github.com/wit-switch/assessment-tax/internal/repository"
	"github.com/wit-switch/assessment-tax/pkg/errorx"
	"github.com/wit-switch/assessment-tax/pkg/validator"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title       Assessment Tax API
// @version     1.0
// @description This is a assessment tax api.
// @BasePath    /
func Execute(cfg *config.Config) {
	dbClient, err := infrastructure.NewPostgresClient(context.Background(), cfg.PostgreSQL)
	if err != nil {
		slog.Error("[!] failed to connect postgres", slog.Any("err", err))
		os.Exit(1)
	}

	repositories := repository.New(repository.Dependencies{
		DB: dbClient,
	})

	services := service.New(service.Dependencies{
		Repositories: repositories,
	})

	hdl := handler.New(handler.Dependencies{
		Services: services,
	})

	mdw := middleware.NewMiddleware(middleware.Dependencies{})

	e := echo.New()

	e.HTTPErrorHandler = httphdl.HTTPErrorHandler
	e.Validator = httphdl.NewValidator(validator.New())
	// with no proxy
	e.IPExtractor = echo.ExtractIPDirect()
	// with proxies using X-Forwarded-For header
	// e.IPExtractor = echo.ExtractIPFromXFFHeader()

	if cfg.Server.Docs {
		e.GET("/docs/*", echoSwagger.WrapHandler)
	}

	e.GET("/healthcheck", func(c echo.Context) error {
		ctx := c.Request().Context()
		if dbErr := dbClient.Ping(ctx); dbErr != nil {
			return dbErr
		}

		return c.String(http.StatusOK, "OK")
	})

	e.Use(
		mdw.Logger(),
	)

	taxGroup := e.Group("/tax")
	{
		taxGroup.POST("/calculations",
			httphdl.BindRoute(
				hdl.Tax.Calculate,
				httphdl.WithBodyParser(),
				httphdl.WithBodyValidator(),
			),
		)
	}

	server := &http.Server{
		Addr:              cfg.Server.HTTPAddress(),
		Handler:           e,
		ReadHeaderTimeout: 30 * time.Second,
	}

	quit := make(chan os.Signal, 1)

	go func() {
		if svErr := e.StartServer(server); !errorx.Is(svErr, http.ErrServerClosed) {
			slog.Error("[!] failed to serve server", slog.Any("err", svErr))
			os.Exit(1)
		}
	}()

	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	{
		<-quit
		slog.Info("gracefully shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if svErr := e.Shutdown(ctx); svErr != nil {
			slog.Error("[!] failed to shutdown server", slog.Any("err", svErr))
		}

		slog.Info("close postgres connection")
		dbClient.Close()

		slog.Info("shutting down the server")
	}
}
