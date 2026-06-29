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
	"example.com/demo/internal/storage"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config", "err", err)
		os.Exit(1)
	}


	storageCtx, storageCancel := context.WithTimeout(context.Background(), 10*time.Second)
	store, err := storage.New(storageCtx, cfg)
	storageCancel()
	if err != nil {
		logger.Error("storage", "err", err)
		os.Exit(1)
	}
	defer func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer closeCancel()
		if err := store.Close(closeCtx); err != nil {
			logger.Error("storage close", "err", err)
		}
	}()
	_ = store // pass into your handlers / services

	logger.Info("starting", "service", "demo", "version", buildinfo.Version, "addr", cfg.HTTPAddr)

	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      server.New(logger),
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
