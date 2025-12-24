package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// On définit le compteur
// "Vec" signifie Vector: on peut trier par étiquettes (method,, status, path)
var httpRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Nombre total de requêtes HTTP",
	},
	[]string{"method", "status", "path"}, // Etiquettes
)

var httpRequestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Durée des requêtes en secondes",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "status", "path"},
)

// Astuce pour lire les Status Code
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Le Middleware
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// On emballe le writter pour pouvoir espionner le code de retour
		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// On laisse passer la requête vers le vrai handler
		next.ServeHTTP(recorder, r)

		//Une fois fini, on calcule la durée
		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(recorder.statusCode)

		//On enregistre dans Prometheus
		//Attention: pour "path", on évite de mettre des ID -> memoire prometheus

		httpRequestsTotal.WithLabelValues(r.Method, statusCode, r.URL.Path).Inc()
		httpRequestDuration.WithLabelValues(r.Method, statusCode, r.URL.Path).Observe(duration)
	})
}
