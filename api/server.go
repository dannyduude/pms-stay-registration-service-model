package api

import (
	"net/http"

	"go-printos-backend-quickstart/api/controllers"
	"go-printos-backend-quickstart/api/models"
)

type server struct {
	b               *Bootstrap
	liveness        http.HandlerFunc
	readiness       http.HandlerFunc
	helloworld      controllers.HelloWorldGetController
	gethelloworlds  controllers.GetHelloWorldsGetController
	creteHelloWorld controllers.CreateHelloWorldsGetController
}

func newServer(b *Bootstrap) *server {
	return &server{
		b:               b,
		liveness:        controllers.Liveness(&b.HealthChecks),
		readiness:       controllers.Readiness(&b.HealthChecks),
		helloworld:      controllers.HelloWorld(),
		gethelloworlds:  controllers.GetHelloWorlds(b.HelloWorldBusiness),
		creteHelloWorld: controllers.CreateHelloWorld(b.HelloWorldBusiness),
	}
}

func (s *server) Liveness(w http.ResponseWriter, r *http.Request, params models.LivenessParams) {
	s.liveness(w, r)
}

func (s *server) Readiness(w http.ResponseWriter, r *http.Request, params models.ReadinessParams) {
	s.readiness(w, r)
}

func (s *server) Helloworld(w http.ResponseWriter, r *http.Request, params models.HelloworldParams) {
	s.helloworld(w, r, params)
}

func (s *server) GetHelloWorlds(w http.ResponseWriter, r *http.Request, params models.GetHelloWorldsParams) {
	s.gethelloworlds(w, r, params)
}

func (s *server) CreateHelloWorld(w http.ResponseWriter, r *http.Request) {
	s.creteHelloWorld(w, r)
}
