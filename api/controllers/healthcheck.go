package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"go-printos-backend-quickstart/api/models"

	healthchecks "github.azc.ext.hp.com/3DSoftware/go-gravity-healthchecks"
	"github.azc.ext.hp.com/3DSoftware/log/v2"
)

type healthCheckInterface interface {
	Liveness(ctx context.Context) healthchecks.Health
	Readiness(ctx context.Context) healthchecks.Health
}

// Liveness responds to the liveness health check.
func Liveness(healthChecks healthCheckInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestLogger := log.GetLogger(r.Context())
		response := healthChecks.Liveness(r.Context())
		if !response.Healthy {
			requestLogger.WithField("dependencies", response.Dependencies).Error("Service is not healthy")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp, reason := json.Marshal(serializeHealthCheckResponse(response))
		if reason != nil {
			w.WriteHeader(http.StatusInternalServerError)
			requestLogger.WithError("error", reason).Error("Unable to marshal response")
			return
		}

		w.Header().Add("Content-Type", "application/json")
		_, err := w.Write(resp)
		if err != nil {
			requestLogger.Warn("Unable to write response on Liveness")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

// Readiness responds to the readiness health check including dependencies.
func Readiness(healthChecks healthCheckInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestLogger := log.GetLogger(r.Context())
		//If initalized if false it means we could not initialized the server with success not furder check required

		response := healthChecks.Readiness(r.Context())
		if !response.Healthy {
			requestLogger.WithField("dependencies", response.Dependencies).Error("Service is not healthy")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		resp, reason := json.Marshal(serializeHealthCheckResponse(response))
		if reason != nil {
			w.WriteHeader(http.StatusInternalServerError)
			requestLogger.WithError("error", reason).Error("Unable to marshal response")
			return
		}

		w.Header().Add("Content-Type", "application/json")
		_, err := w.Write(resp)
		if err != nil {
			requestLogger.Warn("Unable to write response on Liveness")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func serializeHealthCheckResponse(health healthchecks.Health) models.Health {
	return models.Health{
		Healthy:      health.Healthy,
		Dependencies: health.Dependencies,
		Versions:     health.Versions,
	}
}
