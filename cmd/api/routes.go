// Filename cmd/api/routes

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Create new http router instance
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.MethodNotAllowedReponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/forums", app.requiredPermission("forums:read", app.listForumsHandler)) // remove permissions
	router.HandlerFunc(http.MethodPost, "/v1/forums", app.requiredPermission("forums:write", app.createForumHandler))
	router.HandlerFunc(http.MethodGet, "/v1/forums/:id", app.showForumHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/forums/:id", app.requiredPermission("forums:write", app.updateForumHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/forums/:id", app.requiredPermission("forums:write", app.deleteForumHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activate", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}
