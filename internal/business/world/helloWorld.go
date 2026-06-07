package world

import (
	"context"
	"fmt"
	"go-printos-backend-quickstart/internal/entities/world"

	"github.azc.ext.hp.com/3DSoftware/log/v2"
	"github.com/google/uuid"
)

type helloWorldRepository interface {
	Create(ctx context.Context, payload *world.HelloWorld) (*world.HelloWorld, error)
	Get(ctx context.Context, versionID world.HelloWorldID) (*world.HelloWorld, error)
	GetAll(ctx context.Context) ([]world.HelloWorld, error)
	Update(ctx context.Context, versionID world.HelloWorldID, VersionModel world.HelloWorld) (*world.HelloWorld, error)
}

// HelloWorldService is the business logic layer for HelloWorld entity
type HelloWorldService struct {
	hellowordrepository helloWorldRepository
}

// Config holds the configuration for the HelloWorldService
type Config struct {
	HelloWorldRepository helloWorldRepository
}

// New creates a new instance of HelloWorldService
func New(config Config) *HelloWorldService {
	return &HelloWorldService{
		hellowordrepository: config.HelloWorldRepository,
	}
}

// GetHelloWorldById retrieves a HelloWorld entity by its ID
func (b *HelloWorldService) GetHelloWorldById(ctx context.Context, id world.HelloWorldID) (*world.HelloWorld, error) {
	ctxLog := log.GetLogger(ctx)
	ret, err := b.hellowordrepository.Get(ctx, id)
	if err != nil {
		ctxLog.WithError("error", err).WithField("id", id).Error("Error retrieving HelloWorld by ID")
		return nil, err
	}

	return ret, nil
}

// GetAllHelloWorld retrieves all HelloWorld entities
func (b *HelloWorldService) GetAllHelloWorld(ctx context.Context) ([]world.HelloWorld, error) {
	ctxLog := log.GetLogger(ctx)
	ret, err := b.hellowordrepository.GetAll(ctx)
	if err != nil {
		ctxLog.WithError("error", err).Error("Error retrieving HelloWorlds")
		return nil, err
	}

	return ret, nil
}

// CreateHelloWorld creates a new HelloWorld entity
func (b *HelloWorldService) CreateHelloWorld(ctx context.Context, data world.InputHelloWorld) (*world.HelloWorld, error) {
	ctxLog := log.GetLogger(ctx)
	if data.Message == "" {
		return nil, fmt.Errorf("message %v", world.ErrCannotBeEmpty)
	}

	internalHelloWorld := world.HelloWorld{
		Message: data.Message,
		ID:      world.HelloWorldID(uuid.NewString()),
	}

	ret, err := b.hellowordrepository.Create(ctx, &internalHelloWorld)
	if err != nil {
		ctxLog.WithError("error", err).Error("Error creating HelloWorld")
		return nil, err
	}

	return ret, nil

}

// UpdateHelloWorld updates an existing HelloWorld entity
func (b *HelloWorldService) UpdateHelloWorld(ctx context.Context, id world.HelloWorldID, data world.InputHelloWorld) (*world.HelloWorld, error) {
	ctxLog := log.GetLogger(ctx)
	internalHelloWorld, err := b.GetHelloWorldById(ctx, id)
	if err != nil {
		return nil, err
	}

	internalHelloWorld.Message = data.Message
	ret, err := b.hellowordrepository.Update(ctx, id, *internalHelloWorld)
	if err != nil {
		ctxLog.WithError("error", err).WithField("id", id).Error("Error updating HelloWorld")
		return nil, err
	}

	return ret, nil

}
