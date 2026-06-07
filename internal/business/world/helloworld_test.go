package world

import (
	"context"
	"fmt"
	"go-printos-backend-quickstart/internal/entities/world"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockRepo is a simple in-memory implementation of the helloWorldRepository interface.
// It allows to simulate different scenarios for unit testing the HelloWorldService, enabling
// to test both successful operations and error handling without relying on external dependencies.
type mockRepo struct {
	getFunc    func(ctx context.Context, id world.HelloWorldID) (*world.HelloWorld, error)
	getAllFunc func(ctx context.Context) ([]world.HelloWorld, error)
	createFunc func(ctx context.Context, payload *world.HelloWorld) (*world.HelloWorld, error)
	updateFunc func(ctx context.Context, id world.HelloWorldID, model world.HelloWorld) (*world.HelloWorld, error)
}

// Get delegates the call to the injected getFunc.
// It allows to simulate fetching a HelloWorld entity by ID with customizable behavior for testing.
func (m *mockRepo) Get(ctx context.Context, id world.HelloWorldID) (*world.HelloWorld, error) {
	return m.getFunc(ctx, id)
}

// GetAll delegates the call to the injected getAllFunc, enabling simulation of fetching all HelloWorld
// entities with customizable behavior for testing.
func (m *mockRepo) GetAll(ctx context.Context) ([]world.HelloWorld, error) {
	return m.getAllFunc(ctx)
}

// Create delegates the call to the injected createFunc, allowing to simulate creating a HelloWorld entity
// with customizable behavior for testing.
func (m *mockRepo) Create(ctx context.Context, payload *world.HelloWorld) (*world.HelloWorld, error) {
	return m.createFunc(ctx, payload)
}

// Update delegates the call to the injected updateFunc, allowing to simulate updating a HelloWorld entity
// with customizable behavior for testing.
func (m *mockRepo) Update(ctx context.Context, id world.HelloWorldID, model world.HelloWorld) (*world.HelloWorld, error) {
	return m.updateFunc(ctx, id, model)
}

// TestGetHelloWorldById_OK tests the GetHelloWorldById method of HelloWorldService for a successful retrieval scenario.
// It uses a mock repository that returns a predefined HelloWorld entity when queried by ID, and asserts that the returned
// message matches the expected value.
func TestGetHelloWorldById_OK(t *testing.T) {
	repo := &mockRepo{
		getFunc: func(ctx context.Context, id world.HelloWorldID) (*world.HelloWorld, error) {
			return &world.HelloWorld{ID: id, Message: "hello"}, nil
		},
	}

	service := New(Config{HelloWorldRepository: repo})

	res, err := service.GetHelloWorldById(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, "hello", res.Message)
	assert.NotNil(t, res)
}

// TestGetHelloWorldById_Error tests the GetHelloWorldById method of HelloWorldService for an error scenario.
// It uses a mock repository that simulates an error when queried by ID, and asserts that an error is returned
// and the result is nil.
func TestGetHelloWorldById_Error(t *testing.T) {
	repo := &mockRepo{
		getFunc: func(ctx context.Context, id world.HelloWorldID) (*world.HelloWorld, error) {
			return nil, fmt.Errorf("error")
		},
	}

	service := New(Config{HelloWorldRepository: repo})

	res, err := service.GetHelloWorldById(context.Background(), "1")

	assert.Error(t, err)
	assert.Nil(t, res)
}

// TestCreateHelloWorld_Empty tests the CreateHelloWorld method for validation failure when the input message is empty.
// It ensures the service enforces business rules and does not attempt to create an entity with invalid data.
func TestCreateHelloWorld_Empty(t *testing.T) {
	service := New(Config{HelloWorldRepository: &mockRepo{}})

	res, err := service.CreateHelloWorld(context.Background(), world.InputHelloWorld{})

	assert.Error(t, err)
	assert.Nil(t, res)
}

// TestCreateHelloWorld_OK validates the successful creation flow.
// It checks that a valid input is passed to the repository and returned correctly
func TestCreateHelloWorld_OK(t *testing.T) {
	repo := &mockRepo{
		createFunc: func(ctx context.Context, payload *world.HelloWorld) (*world.HelloWorld, error) {
			return payload, nil
		},
	}

	service := New(Config{HelloWorldRepository: repo})

	res, err := service.CreateHelloWorld(context.Background(), world.InputHelloWorld{
		Message: "hello",
	})

	assert.NoError(t, err)
	assert.Equal(t, "hello", res.Message)
	assert.NotNil(t, res)
}

// TestCreateHelloWorld_Error validates the error scenario during creation.
// It checks that an error from the repository is correctly propagated.
func TestCreateHelloWorld_Error(t *testing.T) {
	repo := &mockRepo{
		createFunc: func(ctx context.Context, payload *world.HelloWorld) (*world.HelloWorld, error) {
			return nil, fmt.Errorf("error")
		},
	}

	service := New(Config{HelloWorldRepository: repo})

	res, err := service.CreateHelloWorld(context.Background(), world.InputHelloWorld{
		Message: "hello",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

// TestGetAllHelloWorld_OK tests the GetAllHelloWorld method of HelloWorldService for a successful retrieval scenario.
// It uses a mock repository that returns a predefined list of HelloWorld entities, and asserts that the result matches
// the expected values.
func TestGetAllHelloWorld_OK(t *testing.T) {
	repo := &mockRepo{
		getAllFunc: func(ctx context.Context) ([]world.HelloWorld, error) {
			return []world.HelloWorld{{Message: "hello"}}, nil
		},
	}

	service := New(Config{HelloWorldRepository: repo})

	res, err := service.GetAllHelloWorld(context.Background())

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "hello", res[0].Message)
}

// TestGetAllHelloWorld_Error tests the GetAllHelloWorld method of HelloWorldService for an error scenario.
// It uses a mock repository that simulates an error when fetching all entities, and asserts that an error is returned
// and the result is nil.
func TestGetAllHelloWorld_Error(t *testing.T) {
	repo := &mockRepo{
		getAllFunc: func(ctx context.Context) ([]world.HelloWorld, error) {
			return nil, fmt.Errorf("error")
		},
	}

	service := New(Config{HelloWorldRepository: repo})

	res, err := service.GetAllHelloWorld(context.Background())

	assert.Error(t, err)
	assert.Nil(t, res)
}

// TestUpdateHelloWorld_OK tests the UpdateHelloWorld method of HelloWorldService for a successful update scenario.
// It uses a mock repository that simulates fetching an existing entity and successfully updating it, and asserts that the
// updated message matches the expected value.
func TestUpdateHelloWorld_OK(t *testing.T) {
	repo := &mockRepo{
		getFunc: func(ctx context.Context, id world.HelloWorldID) (*world.HelloWorld, error) {
			return &world.HelloWorld{ID: id, Message: "old"}, nil
		},
		updateFunc: func(ctx context.Context, id world.HelloWorldID, model world.HelloWorld) (*world.HelloWorld, error) {
			return &model, nil
		},
	}

	service := New(Config{HelloWorldRepository: repo})

	res, err := service.UpdateHelloWorld(context.Background(), "1", world.InputHelloWorld{
		Message: "new",
	})

	assert.NoError(t, err)
	assert.Equal(t, "new", res.Message)
}

// TestUpdateHelloWorld_GetError tests the UpdateHelloWorld method of HelloWorldService for a scenario where fetching the existing entity fails.
// It uses a mock repository that simulates an error when fetching by ID, and asserts that an error is returned and the result is nil.
func TestUpdateHelloWorld_GetError(t *testing.T) {
	repo := &mockRepo{
		getFunc: func(ctx context.Context, id world.HelloWorldID) (*world.HelloWorld, error) {
			return nil, fmt.Errorf("error")
		},
	}

	service := New(Config{HelloWorldRepository: repo})

	res, err := service.UpdateHelloWorld(context.Background(), "1", world.InputHelloWorld{
		Message: "new",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

// TestUpdateHelloWorld_UpdateError tests the UpdateHelloWorld method of HelloWorldService for a scenario where updating the existing entity fails.
// It uses a mock repository that simulates an error when updating, and asserts that an error is returned and the result is nil.
func TestUpdateHelloWorld_UpdateError(t *testing.T) {
	repo := &mockRepo{
		getFunc: func(ctx context.Context, id world.HelloWorldID) (*world.HelloWorld, error) {
			return &world.HelloWorld{ID: id}, nil
		},
		updateFunc: func(ctx context.Context, id world.HelloWorldID, model world.HelloWorld) (*world.HelloWorld, error) {
			return nil, fmt.Errorf("error")
		},
	}

	service := New(Config{HelloWorldRepository: repo})

	res, err := service.UpdateHelloWorld(context.Background(), "1", world.InputHelloWorld{
		Message: "new",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}
