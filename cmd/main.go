package main

import (
	"context"
	"errors"
	"fmt"
	exporter "github.com/akyriako/typesense-prometheus-exporter"
	"github.com/caarlos0/env/v11"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	config exporter.Config
	logger *slog.Logger
)

const (
	exitCodeConfigurationError int = 78
)

func init() {
	err := env.Parse(&config)
	if err != nil {
		slog.Error(fmt.Sprintf("parsing env variables failed: %s", err.Error()))
		os.Exit(exitCodeConfigurationError)
	}

	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(config.LogLevel),
	}))

	slog.SetDefault(logger)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	collector := exporter.NewTypesenseCollector(ctx, logger, config)
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		landingPage := exporter.LandingPageTemplate
		w.Write([]byte(landingPage))
	})

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		<-sigChan

		logger.Warn("termination signal received, shutting down gracefully...")
		cancel()
	}()

	server := &http.Server{Addr: fmt.Sprintf(":%d", config.MetricsPort)}

	go func() {
		logger.Info("starting server...", "port", config.MetricsPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(fmt.Sprintf("error starting server: %v", err))
			os.Exit(-1)
		}
	}()

	<-ctx.Done()

	logger.Info("shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Error(fmt.Sprintf("error during server shutdown: %v", err))
		os.Exit(-1)
	}
	logger.Info("exporter shut down successfully")
}
