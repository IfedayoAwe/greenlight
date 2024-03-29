package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.listMoviesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write", app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:read", app.showMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/change-password", app.requireActivatedUser(app.changePasswordHandler))
	router.HandlerFunc(http.MethodPost, "/v1/tokens/password", app.createPasswordResetTokenHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/password", app.resetUserPasswordHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.createActivationTokenHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/users/update-details", app.requireActivatedUser(app.updateUserDetailsHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/users/logout", app.requireActivatedUser(app.userLogoutHandler))
	router.HandlerFunc(http.MethodPut, "/v1/users/profile", app.requireActivatedUser(app.userProfileHandler))
	router.HandlerFunc(http.MethodGet, "/v1/user/profile", app.requireActivatedUser(app.getUserProfileHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/users/delete", app.requireActivatedUser(app.deleteUserAccountHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users/movie-permission", app.requireAdmin(app.addMovieWritePermissionForUser))

	handler := http.StripPrefix("/profile", http.HandlerFunc(app.showProfilePictureHandler))
	router.HandlerFunc(http.MethodGet, "/profile/:filepath", app.requireActivatedUser(handler.ServeHTTP))
	// router.HandlerFunc(http.MethodGet, "/profile/:filepath", app.requireActivatedUser(app.enableGzip(handler.ServeHTTP)))

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())
	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
