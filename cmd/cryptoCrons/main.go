package main

import (
	"context"
	"errors"
	"github.com/Gonnekone/cryptoCrons/internal/config"
	"github.com/Gonnekone/cryptoCrons/internal/http-server/handlers/add"
	"github.com/Gonnekone/cryptoCrons/internal/http-server/handlers/price"
	"github.com/Gonnekone/cryptoCrons/internal/http-server/handlers/remove"
	mwLogger "github.com/Gonnekone/cryptoCrons/internal/http-server/middleware/logger"
	"github.com/Gonnekone/cryptoCrons/internal/lib/logger/handlers/slogpretty"
	"github.com/Gonnekone/cryptoCrons/internal/lib/logger/sl"
	"github.com/Gonnekone/cryptoCrons/internal/parser"
	"github.com/Gonnekone/cryptoCrons/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)
	logger.Info("starting up the application", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	storage, err := postgres.New(cfg.Storage.DSN())
	if err != nil {
		logger.Error("failed to create storage", sl.Err(err))
		os.Exit(1)
	}

	parser := parser.New(logger, cfg.Interval, storage)

	parser.Start(context.Background())
	defer parser.Stop()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)

	router.Route("/currency", func(r chi.Router) {
		r.Get("/price", price.New(logger, storage))

		r.Post("/add", add.New(logger, parser))

		r.Delete("/remove", remove.New(logger, parser))
	})

	logger.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	logger.Info("server started")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to stop server", sl.Err(err))

		return
	}

	logger.Debug("closing storage")

	storage.Close()
	logger.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
