package controllers_test

import (
	"context"
	"encoding/json"
	"go-printos-backend-quickstart/api/controllers"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-printos-backend-quickstart/api/models"

	healthchecks "github.azc.ext.hp.com/3DSoftware/go-gravity-healthchecks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHealthchecks struct {
	mock.Mock
}

func (m *mockHealthchecks) Liveness(ctx context.Context) healthchecks.Health {
	args := m.Called(ctx)
	return args.Get(0).(healthchecks.Health)
}

func (m *mockHealthchecks) Readiness(ctx context.Context) healthchecks.Health {
	args := m.Called(ctx)
	return args.Get(0).(healthchecks.Health)
}

func TestHealthCheck_Liveness(t *testing.T) {
	tests := []struct {
		name            string
		healthyResponse healthchecks.Health
		expectedStatus  int
		expectedHealth  bool
	}{
		{
			name:            "returns internal server error when healthy is false",
			healthyResponse: healthchecks.Health{Healthy: false},
			expectedStatus:  http.StatusInternalServerError,
			expectedHealth:  false,
		},
		{
			name:            "returns ok when healthy is true",
			healthyResponse: healthchecks.Health{Healthy: true},
			expectedStatus:  http.StatusOK,
			expectedHealth:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			health := mockHealthchecks{}
			health.On("Liveness", mock.Anything).Return(tt.healthyResponse)

			reqURI := "/healthcheck/liveness"
			req, err := http.NewRequest(http.MethodGet, reqURI, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			processHandler := controllers.Liveness(&health)
			handler := func(handlerFn http.HandlerFunc) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					handlerFn(w, r)
				}
			}(processHandler)
			handler.ServeHTTP(rr, req)
			response := rr.Result()
			assert.Equal(t, tt.expectedStatus, response.StatusCode)

			if tt.expectedStatus != http.StatusOK {
				return
			}

			data, err := io.ReadAll(response.Body)
			assert.Nil(t, err)

			currentResponse := models.Health{}
			err = json.Unmarshal(data, &currentResponse)
			assert.Nil(t, err)

			assert.Equal(t, tt.expectedHealth, currentResponse.Healthy)
		})
	}

}

func TestHealthCheck_Readiness(t *testing.T) {
	tests := []struct {
		name            string
		healthyResponse healthchecks.Health
		expectedStatus  int
		expectedHealth  bool
	}{
		{
			name:            "returns service unavailable when healthy is false",
			healthyResponse: healthchecks.Health{Healthy: false},
			expectedStatus:  http.StatusServiceUnavailable,
			expectedHealth:  false,
		},
		{
			name:            "returns ok when healthy is true",
			healthyResponse: healthchecks.Health{Healthy: true},
			expectedStatus:  http.StatusOK,
			expectedHealth:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			health := mockHealthchecks{}
			health.On("Readiness", mock.Anything).Return(tt.healthyResponse)

			reqURI := "/healthcheck/readiness"
			req, err := http.NewRequest(http.MethodGet, reqURI, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			processHandler := controllers.Readiness(&health)
			handler := func(handlerFn http.HandlerFunc) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					handlerFn(w, r)
				}
			}(processHandler)
			handler.ServeHTTP(rr, req)
			response := rr.Result()
			assert.Equal(t, tt.expectedStatus, response.StatusCode)

			if tt.expectedStatus != http.StatusOK {
				return
			}

			data, err := io.ReadAll(response.Body)
			assert.Nil(t, err)

			currentResponse := models.Health{}
			err = json.Unmarshal(data, &currentResponse)
			assert.Nil(t, err)

			assert.Equal(t, tt.expectedHealth, currentResponse.Healthy)
		})
	}
}
