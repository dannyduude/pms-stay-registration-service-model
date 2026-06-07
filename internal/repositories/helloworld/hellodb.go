package hellowordrepository

import (
	"context"

	"go-printos-backend-quickstart/internal/entities/world"

	mongoPersistence "github.azc.ext.hp.com/3DSoftware/go-gravity-mongo-interface/v2/persistence"

	"github.azc.ext.hp.com/3DSoftware/log/v2"
)

type mongoInterface interface {
}
type helloWorldDB struct {
	world.HelloWorld `bson:",inline"`
	mongoPersistence.Context
}

// HelloWorldRepository is the data access layer for HelloWorld entity
type HelloWorldRepository struct {
	mongo mongoInterface
}

// Config holds the configuration for the HelloWorldRepository
type Config struct {
	Mongo mongoInterface
}

// New creates a new instance of HelloWorldRepository
func New(config Config) *HelloWorldRepository {
	return &HelloWorldRepository{
		mongo: config.Mongo,
	}
}

// Create inserts a new HelloWorld entity into the database
func (p *HelloWorldRepository) Create(ctx context.Context, helloWorld *world.HelloWorld) (*world.HelloWorld, error) {
	ctxLog := log.GetLogger(ctx)
	ctxLog.Trace("[Begin]hellowordrepository.Create")

	dbModel := helloWorldDB{
		HelloWorld: *helloWorld,
	}

	return &dbModel.HelloWorld, nil
}

// Get retrieves a HelloWorld entity by its ID from the database
func (p *HelloWorldRepository) Get(ctx context.Context, helloWorldID world.HelloWorldID) (*world.HelloWorld, error) {
	ctxLog := log.GetLogger(ctx)
	ctxLog.Trace("[Begin]hellowordrepository.Get HellowId[%s]", helloWorldID)

	if helloWorldID != "11111111-1111-1111-1111-111111111111" {
		ctxLog.Debugf("[End]hellowordrepository.Get HellowId[%s] not found", helloWorldID)
		return nil, world.ErrNotFound
	}

	dbModel := helloWorldDB{
		HelloWorld: world.HelloWorld{
			ID:      helloWorldID,
			Message: "",
		},
	}

	return &dbModel.HelloWorld, nil
}

// GetAll retrieves all HelloWorld entities from the database
func (p *HelloWorldRepository) GetAll(ctx context.Context) ([]world.HelloWorld, error) {
	ctxLog := log.GetLogger(ctx)
	ctxLog.Trace("[Begin]hellowordrepository.GetAll")

	dbModel := helloWorldDB{
		HelloWorld: world.HelloWorld{
			ID:      "",
			Message: "",
		},
	}

	ctxLog.Trace("[End]hellowordrepository.GetAll")
	return []world.HelloWorld{dbModel.HelloWorld}, nil
}

// Update modifies an existing HelloWorld entity in the database
func (p *HelloWorldRepository) Update(ctx context.Context, helloWorldID world.HelloWorldID, helloWorld world.HelloWorld) (*world.HelloWorld, error) {
	ctxLog := log.GetLogger(ctx)
	ctxLog.Trace("[Begin]hellowordrepository.Update")

	if helloWorldID != "11111111-1111-1111-1111-111111111111" {
		return nil, world.ErrNotFound
	}

	ctxLog.Trace("[End]hellowordrepository.Update")
	return &helloWorld, nil
}
