package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requireActivatedUser(app.listMoviesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requireActivatedUser(app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requireActivatedUser(app.showMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requireActivatedUser(app.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requireActivatedUser(app.deleteMovieHandler))

	router.HandlerFunc(http.MethodPost, "/v1/moduleinfo", app.createModuleInfoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/moduleinfo/:id", app.showModuleInfoHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/moduleinfo/:id", app.updateModuleInfoHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/moduleinfo/:id", app.deleteModuleInfoHandler)

	router.HandlerFunc(http.MethodPost, "/v1/departamentinfo", app.createDepInfoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/departamentinfo/:id", app.getDepInfoHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodPost, "/v1/userinfo", app.createUserInfoHandler)
	router.HandlerFunc(http.MethodPut, "/v1/userinfo/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/userinfo", app.getAllUserInfoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/userinfo/:id", app.getUserInfoHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/userinfo/:id", app.editUserInfoHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/userinfo/:id", app.deleteUserInfoHandler)

	return app.recoverPanic(app.authenticate(router))
}
