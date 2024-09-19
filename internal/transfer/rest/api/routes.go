package api

import (
	"net/http"

	swagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/PatrykPasterny/scooter-rental/docs"
)

const (
	api          = "/api"
	version      = "/v1"
	scootersPath = "/scooters"
	rentPath     = "/rent"
	freePath     = "/free"
	swaggerDocs  = "/api-docs"
)

// registerRoutes sets service routes.
func (s *Server) registerRoutes() {
	versionRoute := s.router.PathPrefix(api + version).Subrouter()

	versionRoute.PathPrefix(swaggerDocs).Handler(swagger.WrapHandler)

	versionRoute.Path(scootersPath).Methods(http.MethodGet).
		HandlerFunc(AuthenticateUser(s.getScooters, s.logger, s.eligibleUsers))

	versionRoute.Path(rentPath).Methods(http.MethodPost).
		HandlerFunc(AuthenticateUser(s.rentScooter, s.logger, s.eligibleUsers))
	versionRoute.Path(freePath).Methods(http.MethodPost).
		HandlerFunc(AuthenticateUser(s.freeScooter, s.logger, s.eligibleUsers))
}
