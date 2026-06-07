package healthchecks

//REMOVE IF NOT USING MONGO: Remove file
import (
	"context"

	"github.azc.ext.hp.com/3DSoftware/log/v2"

	mongo "github.azc.ext.hp.com/3DSoftware/go-gravity-mongo-driver/v6"
)

type MongoLiveness struct {
	Mongo mongo.MongoInterface
}

func NewMongoLiveness(Mongo mongo.MongoInterface) *MongoLiveness {
	return &MongoLiveness{Mongo: Mongo}
}

func (m *MongoLiveness) IsHealthy(ctx context.Context) bool {
	logger := log.GetLogger(ctx)
	err := m.Mongo.Ping(ctx)
	if err != nil {
		logger.WithError("error", err).Errorf("MongoDB - Liveness %v", err)
	}
	return (err == nil)
}
