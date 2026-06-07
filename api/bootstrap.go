//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config ./oapi-codegen/server.cfg.yml openapi.yml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config ./oapi-codegen/types.cfg.yml openapi.yml

package api

import (
	"context"
	"io"
	"net/http"

	"go-printos-backend-quickstart/api/models"
	helloWorldBusiness "go-printos-backend-quickstart/internal/business/world"
	internalHealthchecks "go-printos-backend-quickstart/internal/healthchecks"
	hellowordrepository "go-printos-backend-quickstart/internal/repositories/helloworld"

	healthchecks "github.azc.ext.hp.com/3DSoftware/go-gravity-healthchecks"

	mongo "github.azc.ext.hp.com/3DSoftware/go-gravity-mongo-driver/v6"
	featureflagsdk "github.azc.ext.hp.com/3DSoftware/gravity-feature-flag/v2"
	featureflagsdkMiddleware "github.azc.ext.hp.com/3DSoftware/gravity-feature-flag/v2/middleware"

	"github.azc.ext.hp.com/3DSoftware/go-gravity-jsonapi/errors"
	requestcontext "github.azc.ext.hp.com/3DSoftware/go-gravity-middleware-requestid"

	"github.azc.ext.hp.com/3DSoftware/go-cleanhttp/middleware"
	authMiddleware "github.azc.ext.hp.com/3DSoftware/gravity-middleware-auth/v2/dqarauth"
	authModels "github.azc.ext.hp.com/3DSoftware/gravity-middleware-auth/v2/models"
	"github.azc.ext.hp.com/3DSoftware/log/v2"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	readiness   = "*/healthcheck/readiness*"
	liveness    = "*/healthcheck/liveness*"
	metrics     = "*/metrics*"
	serviceName = "backend"
)

// BootstrapConfig defines the necessary bootstrap configuration.
type BootstrapConfig struct {
	ServiceVersion     string
	ExecuteAsDryRun    bool
	LogWritter         io.Writer
	FeatureFlagsConfig featureflagsdk.Config
	MongoConfig        mongo.Config
}

func (b BootstrapConfig) validateFeatureFlagsdk(input featureflagsdk.Config) error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.SDKkey, validation.Required),
	)
}

func (b BootstrapConfig) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.FeatureFlagsConfig, validation.Required, validation.By(func(value interface{}) error {
			return b.validateFeatureFlagsdk(value.(featureflagsdk.Config))
		})),
	)
}

// Bootstrap exports the necessary objects to run the backend.
type Bootstrap struct {
	ServiceVersion     string
	Logger             log.Logger
	HelloWorldBusiness *helloWorldBusiness.HelloWorldService
	FeatureFlags       featureflagsdk.FeatureFlagsInterface
	HealthChecks       healthchecks.HealthChecks
}

// NewBootstrap initializes and injects all the dependencies.
func NewBootstrap(config BootstrapConfig) *Bootstrap {
	healthChecks := healthchecks.NewHealthChecks(healthchecks.HealthChecksConfig{
		LocalMode:  config.ExecuteAsDryRun,
		SvcName:    serviceName,
		SvcVersion: config.ServiceVersion,
	})

	healthChecks.SetDependency("helloWorldRepo", healthchecks.NewSimpleHealthCheck(true))
	healthChecks.SetDependency("feature-flags", healthchecks.NewSimpleHealthCheck(true))
	healthChecks.SetDependency("mongodb", healthchecks.NewSimpleHealthCheck(true))
	healthChecks.SetDependency("business", healthchecks.NewSimpleHealthCheck(true))

	if config.LogWritter == nil {
		config.LogWritter = io.Discard
	}
	logger := log.New(config.LogWritter, log.InfoLevel, config.ServiceVersion)

	ctx := context.Background()

	err := config.Validate()
	if err != nil {
		logger.Errorf("missing parameters: %v", err)
	}

	//REMOVE IF NOT USING MONGO
	mongo, err := mongo.New(ctx, config.MongoConfig, logger)
	if err != nil {
		logger.Errorf("error creating mongodb connection: %v", err)
		healthChecks.SetDependency("mongodb", healthchecks.NewSimpleHealthCheck(false))
	} else {
		logger.Info("Connected to mongodb successfully")
		healthChecks.SetDependency("mongodb", internalHealthchecks.NewMongoLiveness(mongo))
	}

	helloWorldRepo := hellowordrepository.New(hellowordrepository.Config{
		Mongo: mongo,
	})

	featureflags, err := featureflagsdk.NewClient(config.FeatureFlagsConfig)
	if err != nil {
		logger.Errorf("error connecting  feature-flags: %v", err)
		healthChecks.SetDependency("feature-flags", healthchecks.NewSimpleHealthCheck(false))
	}

	helloWorldBusiness := helloWorldBusiness.New(helloWorldBusiness.Config{
		HelloWorldRepository: helloWorldRepo,
	})
	if err != nil {
		logger.Errorf("Failed to initialize business layer Error: %v", err)
		healthChecks.SetDependency("business", healthchecks.NewSimpleHealthCheck(false))
	}

	return &Bootstrap{
		Logger:             logger,
		HelloWorldBusiness: helloWorldBusiness,
		ServiceVersion:     config.ServiceVersion,
		HealthChecks:       healthChecks,
		FeatureFlags:       featureflags,
	}
}

// NewHandler returns an HTTP handler.
func NewHandler(bootstrap *Bootstrap) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(requestcontext.RequestID)
	r.Use(middleware.Logger(bootstrap.Logger))
	r.Use(middleware.OWASPSecureHeaders)
	r.Use(errors.JsonApiRecoverer)
	r.Use(authMiddleware.Authorize(jwt.NewParser(),
		authModels.UrlMatch{Method: http.MethodGet, Pattern: readiness},
		authModels.UrlMatch{Method: http.MethodGet, Pattern: liveness},
		authModels.UrlMatch{Method: http.MethodGet, Pattern: metrics},
	))
	r.Use(featureflagsdkMiddleware.PopulateContext(bootstrap.FeatureFlags, []featureflagsdkMiddleware.FeatureFlagDefinition{}))

	r.Get("/metrics", http.HandlerFunc(promhttp.Handler().ServeHTTP))
	return models.HandlerFromMux(newServer(bootstrap), r)
}
