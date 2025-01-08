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

	//logger = logger.With("namespace", config.Namespace).With("cluster", config.Cluster)
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
		landingPage := `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Typesense Prometheus Exporter</title>
				<style>
					* {
						margin: 0;
						padding: 0;
						box-sizing: border-box;
					}
					html, body {
						height: 100%;
						display: flex;
						align-items: center;
						justify-content: center;
						background-color: black;
						color: white;
						font-family: Arial, sans-serif;
					}
					.container {
						text-align: center;
					}
					img {
						margin-top: 20px;
						max-width: 200px;
						height: auto;
						margin-bottom: 20px;
					}
					a {
						text-decoration: none;
						color: #00bcd4;
						font-size: 18px;
					}
					a:hover {
						text-decoration: underline;
					}
				</style>
			</head>
			<body>
				<div class="container">
					<img src="https://prometheus.io/assets/prometheus_logo_grey.svg" alt="Prometheus Logo"/><br/>
					<img src="https://typesense.org/_nuxt/img/typesense_logo_white.0f9fb0a.svg" alt="Typesense Logo"/>
					<p><a href="/metrics">Go to Metrics</a></p>
				</div>
			</body>
			</html>
		`
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
