package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/demo/internal/buildinfo"
	"example.com/demo/internal/config"
	"example.com/demo/internal/server"
	"example.com/demo/internal/telemetry"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config", "err", err)
		os.Exit(1)
	}


	otelInitCtx, otelInitCancel := context.WithTimeout(context.Background(), 10*time.Second)
	shutdownTracer, err := telemetry.Init(otelInitCtx, "demo", buildinfo.Version)
	otelInitCancel()
	if err != nil {
		logger.Error("init telemetry", "err", err)
		os.Exit(1)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracer(shutdownCtx); err != nil {
			logger.Error("shutdown telemetry", "err", err)
		}
	}()



	e := server.New(logger)
	logger.Info("starting", "service", "demo", "version", buildinfo.Version, "addr", cfg.HTTPAddr)

	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      e,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server", "err", err)
			stop()
		}
	}()


	<-ctx.Done()
	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown", "err", err)
		os.Exit(1)
	}
}
