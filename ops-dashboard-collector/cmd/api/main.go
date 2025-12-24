package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mtk14m/ops-dashboard-collector/internal/handlers"
	"github.com/mtk14m/ops-dashboard-collector/internal/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	//1. Configuration du Logger Structuré (JSON)
	//Indispensable pour Loki/Datadog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//2. Wiring (Injection des dépendances)
	HealthHandler := &handlers.HealthHandler{Logger: logger}

	mux := http.NewServeMux()
	//ROUTES METIERS
	mux.Handle("/health", HealthHandler)
	//ROUTES TECHNIQUES
	mux.Handle("/metrics", promhttp.Handler())

	handlerWithMetrics :=  middleware.MetricsMiddleware(mux)

	//3. Configuration du serveur avec Timeouts (Sécurité SRE)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handlerWithMetrics, // On a ajouté ici notre middleware
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	//4. Démarrage dans une Goroutine pour ne pas bloquer
	go func() {
		logger.Info("Starting server", "port", port, "env", "production")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	//5. Graceful shutdown

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit //On bloque juqu'a recevoir le signal
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Error forced shutting down server", "error", err)
		os.Exit(1)
	}

	logger.Info("Server gracefully stopped")

}
