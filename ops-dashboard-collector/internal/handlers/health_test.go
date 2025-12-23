package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHealthHandler_ServeHTTP(t *testing.T) {
	// 1. Setup : On crée un logger "poubelle" pour le test
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	handler := &HealthHandler{Logger: logger}

	// 2. Execution : On simule une requête HTTP
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// ResponseRecorder est un "faux" navigateur qui enregistre la réponse
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// 3. Assertions (Vérifications)

	// Vérifier le Code HTTP
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Mauvais status code : obtenu %v, attendu %v", status, http.StatusOK)
	}

	// Vérifier le corps JSON
	var response HealthResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Réponse invalide (pas du JSON) : %v", err)
	}

	if response.Status != "UP" {
		t.Errorf("Mauvais status dans le JSON : obtenu %v, attendu UP", response.Status)
	}
}
