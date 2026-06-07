package healthchecks_test

import (
	"context"
	"fmt"
	"go-printos-backend-quickstart/internal/healthchecks"
	"testing"

	mongomocks "github.azc.ext.hp.com/3DSoftware/go-gravity-mongo-driver/v6/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMongoDbLiveness_Ping(t *testing.T) {
	tests := []struct {
		name        string
		returnError error
		result      bool
	}{
		{name: "error", returnError: fmt.Errorf("mongo unavailble"), result: false},
		{name: "no-error", returnError: nil, result: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedMongo := new(mongomocks.MongoInterfaceMock)
			mockedMongo.On("Ping", mock.Anything).Return(tt.returnError)
			mongoLive := healthchecks.NewMongoLiveness(mockedMongo)
			got := mongoLive.IsHealthy(context.Background())
			assert.Equal(t, got, tt.result)
		})

	}
}
