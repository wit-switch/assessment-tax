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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Execute ...
// @title                     Assessment Tax API
// @version                   1.0
// @description               This is a assessment tax api.
// @BasePath                  /
// @securityDefinitions.basic BasicAuth
// @name                      Authorization
// @description               BasicAuth protects our entity endpoints.
func Execute(cfg *config.Config) {
	validate, err := validator.New()
	if err != nil {
		slog.Error("[!] failed to create validate", slog.Any("err", err))
		os.Exit(1)
	}

	dbClient, err := infrastructure.NewPostgresClient(context.Background(), cfg.PostgreSQL)
	if err != nil {
		slog.Error("[!] failed to connect postgres", slog.Any("err", err))
		os.Exit(1)
	}

	hdl := getHandler(dbClient)

	mdw := middleware.NewMiddleware(middleware.Dependencies{
		Auth: cfg.Auth,
	})

	e := echo.New()

	e.HTTPErrorHandler = httphdl.HTTPErrorHandler
	e.Validator = httphdl.NewValidator(validate)
	e.IPExtractor = echo.ExtractIPDirect()

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
				hdl.Tax.Calculate, httphdl.WithBodyParser(), httphdl.WithBodyValidator(),
			),
		)
		taxGroup.POST("/calculations/upload-csv", httphdl.BindRoute(hdl.Tax.CalculateFromCSV))
	}

	adminGroup := e.Group("/admin")
	{
		adminGroup.Use(mdw.Auth())

		adminGroup.POST("/deductions/personal",
			httphdl.BindRoute(
				hdl.Admin.UpdatePersonalDeduct, httphdl.WithBodyParser(), httphdl.WithBodyValidator(),
			),
		)
		adminGroup.POST("/deductions/k-receipt",
			httphdl.BindRoute(
				hdl.Admin.UpdateKReceiptDeduct, httphdl.WithBodyParser(), httphdl.WithBodyValidator(),
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

func getHandler(dbClient *pgxpool.Pool) *handler.Handler {
	repositories := repository.New(repository.Dependencies{
		DB: dbClient,
	})

	services := service.New(service.Dependencies{
		Repositories: repositories,
	})

	hdl := handler.New(handler.Dependencies{
		Services: services,
	})

	return hdl
}
