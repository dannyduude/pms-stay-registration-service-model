package controllers

import (
	"context"
	"fmt"
	"go-printos-backend-quickstart/api/models"
	"go-printos-backend-quickstart/internal/entities/world"
	"io"
	"net/http"
	"regexp"

	"go-printos-backend-quickstart/api/models/logical"

	jsonApiErrors "github.azc.ext.hp.com/3DSoftware/go-gravity-jsonapi/errors"
	"github.azc.ext.hp.com/3DSoftware/log/v2"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/manyminds/api2go/jsonapi"
)

const (
	contentTypeJSONAPI = "application/vnd.api+json"
	contentTypeHeader  = "Content-type"
)

type HelloWorldGetController func(http.ResponseWriter, *http.Request, models.HelloworldParams)
type GetHelloWorldsGetController func(http.ResponseWriter, *http.Request, models.GetHelloWorldsParams)
type CreateHelloWorldsGetController func(http.ResponseWriter, *http.Request)

type helloWorldBusiness interface {
	GetAllHelloWorld(ctx context.Context) ([]world.HelloWorld, error)
	CreateHelloWorld(ctx context.Context, data world.InputHelloWorld) (*world.HelloWorld, error)
}

func HelloWorld() HelloWorldGetController {
	return func(w http.ResponseWriter, r *http.Request, params models.HelloworldParams) {
		requestLogger := log.GetLogger(r.Context())
		requestLogger.Info("Log for track hello world")

		if r.ContentLength > 0 {
			http.Error(w, "unexpected content", http.StatusBadRequest)
			return
		}

		if params.Name == nil || *params.Name == "" {
			fmt.Fprintf(w, "Hello, World!")
			return
		}

		err := validation.Validate(params.Name,
			validation.Length(3, 10).Error("must be between 3 to 10 characters"),
			validation.Match(regexp.MustCompile("^[a-zA-Z]*$")).Error("must only contain characters"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Hello, %s!", *params.Name)
	}
}

func GetHelloWorlds(hellowWorldBusiness helloWorldBusiness) GetHelloWorldsGetController {
	return func(w http.ResponseWriter, r *http.Request, params models.GetHelloWorldsParams) {
		requestLogger := log.GetLogger(r.Context())
		requestLogger.Info("Log for track hello world")

		if r.ContentLength > 0 {
			http.Error(w, "unexpected content", http.StatusBadRequest)
			return
		}

		helloWorlds, err := hellowWorldBusiness.GetAllHelloWorld(r.Context())
		if err != nil {
			requestLogger.Errorf("Unknown error %v", err)
			jsonApiErrors.JsonAPIError(w, r, "Internal Error", http.StatusInternalServerError)
			return
		}

		var logicalHelloWorlds = []logical.HelloWorld{}
		for _, helloWorld := range helloWorlds {
			logicalHelloWorlds = append(logicalHelloWorlds, serializeHelloWorld(helloWorld))
		}

		jsonResp, err := jsonapi.Marshal(logicalHelloWorlds)
		if err != nil {
			requestLogger.Error("error encoding response")
			jsonApiErrors.JsonAPIError(w, r, "error encoding response", http.StatusInternalServerError)
			return
		}

		w.Header().Add(contentTypeHeader, contentTypeJSONAPI)
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(jsonResp); err != nil {
			requestLogger.Warn("Unable to write response on get of HelloWorlds")
		}
	}
}

func CreateHelloWorld(hellowWorldBusiness helloWorldBusiness) CreateHelloWorldsGetController {
	return func(w http.ResponseWriter, r *http.Request) {
		requestLogger := log.GetLogger(r.Context())
		requestLogger.Info("Log for track hello world")
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			requestLogger.Error("error reading request body")
			jsonApiErrors.JsonAPIError(w, r, "error reading request body", http.StatusBadRequest)
			return
		}

		var helloWorldAttributes logical.POSTHelloWorldAttributes
		if err := jsonapi.Unmarshal(payload, &helloWorldAttributes); err != nil {
			requestLogger.Error("invalid body format for HelloWorld", err)
			jsonApiErrors.JsonAPIError(w, r, "error encoding response", http.StatusBadRequest)
			return
		}

		helloWorld, err := hellowWorldBusiness.CreateHelloWorld(r.Context(), deserializePOSTHelloWorld(&helloWorldAttributes))
		if err != nil {
			requestLogger.Errorf("Unknown error %v", err)
			jsonApiErrors.JsonAPIError(w, r, "Internal Error", http.StatusInternalServerError)
			return
		}

		jsonResp, err := jsonapi.Marshal(serializeHelloWorld(*helloWorld))
		if err != nil {
			requestLogger.Error("error encoding response")
			jsonApiErrors.JsonAPIError(w, r, "error encoding response", http.StatusInternalServerError)
			return
		}

		w.Header().Add(contentTypeHeader, contentTypeJSONAPI)
		w.WriteHeader(http.StatusCreated)
		if _, err = w.Write(jsonResp); err != nil {
			requestLogger.Warn("Unable to write response on get of HelloWorld")
		}
	}
}

func serializeHelloWorld(doc world.HelloWorld) logical.HelloWorld {
	helloWorld := logical.HelloWorld{
		Base: logical.Base{
			Id: string(doc.ID),
		},
		HelloWorldAttributes: models.HelloWorldAttributes{
			Message:   doc.Message,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
		},
	}
	helloWorld.SerializeRelationships(doc)
	return helloWorld
}

func deserializePOSTHelloWorld(helloWorld *logical.POSTHelloWorldAttributes) world.InputHelloWorld {
	return world.InputHelloWorld{
		Message: *helloWorld.Message,
	}
}
