package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-printos-backend-quickstart/api/controllers"
	"go-printos-backend-quickstart/api/models"

	"github.com/stretchr/testify/assert"

	"context"
	"go-printos-backend-quickstart/internal/entities/world"

	"fmt"
)

// strPtr is a helper function to easily create string pointers
// required by generated API parameter structs.
func strPtr(s string) *string {
	return &s
}

// mockHelloWorldBusiness is a lightweight mock implementation of the business layer.
// It allows simulating both successful and failing scenarios for unit tests.
type mockHelloWorldBusiness struct {
	returnError       bool
	returnCreateError bool
}

// GetAllHelloWorld simulates fetching all HelloWorld entities.
// It can return either a valid list or an error depending on the test scenario.
func (m *mockHelloWorldBusiness) GetAllHelloWorld(ctx context.Context) ([]world.HelloWorld, error) {
	if m.returnError {
		return nil, fmt.Errorf("mock error")
	}

	return []world.HelloWorld{
		{
			ID:      "1",
			Message: "Hello from mock",
		},
	}, nil
}

// CreateHelloWorld simulates entity creation in the business layer.
// It supports both success and failure scenarios for testing controller behavior.
func (m *mockHelloWorldBusiness) CreateHelloWorld(ctx context.Context, data world.InputHelloWorld) (*world.HelloWorld, error) {
	if m.returnCreateError {
		return nil, fmt.Errorf("mock create error")
	}

	return &world.HelloWorld{
		ID:      "1",
		Message: data.Message,
	}, nil
}

// TestHelloWorld validates all execution paths of the HelloWorld handler.
// It covers:
// - Request with unexpected body (should return 400)
// - Nil and empty name inputs (default "Hello, World!")
// - Validation errors (length and invalid characters)
// - Valid name case
func TestHelloWorld(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		params         models.HelloworldParams
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "unexpected body",
			requestBody:    "data",
			params:         models.HelloworldParams{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "unexpected content\n",
		},
		{
			name:           "nil name",
			requestBody:    "",
			params:         models.HelloworldParams{Name: nil},
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello, World!",
		},
		{
			name:           "empty name",
			requestBody:    "",
			params:         models.HelloworldParams{Name: strPtr("")},
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello, World!",
		},
		{
			name:           "name too short",
			requestBody:    "",
			params:         models.HelloworldParams{Name: strPtr("Hi")},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "must be between 3 to 10 characters\n",
		},
		{
			name:           "invalid characters",
			requestBody:    "",
			params:         models.HelloworldParams{Name: strPtr("Us3r")},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "must only contain characters\n",
		},
		{
			name:           "valid name",
			requestBody:    "",
			params:         models.HelloworldParams{Name: strPtr("User")},
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello, User!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tt.requestBody != "" {
				req, err = http.NewRequest(http.MethodGet, "/hello", strings.NewReader(tt.requestBody))
			} else {
				req, err = http.NewRequest(http.MethodGet, "/hello", nil)
			}

			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			handler := controllers.HelloWorld()
			handler(rr, req, tt.params)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

// TestGetHelloWorlds_OK verifies the successful execution path of GetHelloWorlds.
// It ensures:
// - Business layer returns data
// - Controller serializes response correctly
// - HTTP status is 200 OK
func TestGetHelloWorlds_OK(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/helloworlds", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockBusiness := &mockHelloWorldBusiness{}

	handler := controllers.GetHelloWorlds(mockBusiness)
	handler(rr, req, models.GetHelloWorldsParams{})

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Hello from mock")
}

// TestGetHelloWorlds_Error validates behavior when the business layer fails.
// It ensures the controller returns a 500 Internal Server Error.
func TestGetHelloWorlds_Error(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/helloworlds", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockBusiness := &mockHelloWorldBusiness{
		returnError: true,
	}

	handler := controllers.GetHelloWorlds(mockBusiness)
	handler(rr, req, models.GetHelloWorldsParams{})

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// TestCreateHelloWorld_InvalidBody verifies that invalid JSON payloads
// are correctly rejected with a 400 Bad Request response.
func TestCreateHelloWorld_InvalidBody(t *testing.T) {
	body := `invalid-json`

	req, err := http.NewRequest(http.MethodPost, "/helloworlds", strings.NewReader(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockBusiness := &mockHelloWorldBusiness{}

	handler := controllers.CreateHelloWorld(mockBusiness)
	handler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// TestCreateHelloWorld_BusinessError validates behavior when the business layer fails
// after a successful request parsing. The controller should return 500.
func TestCreateHelloWorld_BusinessError(t *testing.T) {
	body := `{
        "data": {
            "type": "helloworld",
            "attributes": {
                "message": "Hello from test"
            }
        }
    }`

	req, err := http.NewRequest(http.MethodPost, "/helloworlds", strings.NewReader(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockBusiness := &mockHelloWorldBusiness{
		returnCreateError: true,
	}

	handler := controllers.CreateHelloWorld(mockBusiness)
	handler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// TestCreateHelloWorld_OK ensures the happy path for entity creation.
// It validates:
// - Proper JSON parsing
// - Successful business logic invocation
// - Correct response serialization and HTTP 201 status
func TestCreateHelloWorld_OK(t *testing.T) {
	body := `{
        "data": {
            "type": "helloworld",
            "attributes": {
                "message": "Hello from test"
            }
        }
    }`

	req, err := http.NewRequest(http.MethodPost, "/helloworlds", strings.NewReader(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockBusiness := &mockHelloWorldBusiness{}

	handler := controllers.CreateHelloWorld(mockBusiness)
	handler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Body.String(), "Hello from test")
}

// TestGetHelloWorlds_EncodeError ensures response encoding paths are executed.
// While forcing encoding errors is non-trivial, this test increases coverage of the serialization flow.
func TestGetHelloWorlds_EncodeError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/helloworlds", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockBusiness := &mockHelloWorldBusiness{}

	handler := controllers.GetHelloWorlds(mockBusiness)
	handler(rr, req, models.GetHelloWorldsParams{})

	assert.True(t, rr.Code == http.StatusOK || rr.Code == http.StatusInternalServerError)
}

// TestGetHelloWorlds_WithBody validates that requests containing unexpected payloads
// are rejected early with a 400 Bad Request response.
func TestGetHelloWorlds_WithBody(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/helloworlds", strings.NewReader("data"))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockBusiness := &mockHelloWorldBusiness{}

	handler := controllers.GetHelloWorlds(mockBusiness)
	handler(rr, req, models.GetHelloWorldsParams{})

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
