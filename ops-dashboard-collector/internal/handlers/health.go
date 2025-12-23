package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Response structure
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Hostanme  string    `json:"hostname"`
	Version   string    `json:"version"`
}

type HealthHandler struct {
	Logger *slog.Logger
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Health check request received",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)

	hostname, _ := os.Hostname()

	resp := HealthResponse{
		Status:    "UP",
		Timestamp: time.Now(),
		Hostanme:  hostname,
		Version:   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.Logger.Error("Error encoding JSON response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
