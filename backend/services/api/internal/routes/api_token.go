package routes

import (
	apitokenctrl "ava/api/internal/controller/apitoken"
	"ava/api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// Session-only: a token that could mint another token would make revocation meaningless.
func apiTokenRoutes(
	router fiber.Router,
	controller *apitokenctrl.Controller,
	mw *middleware.Middleware,
) {
	// Authenticated, then SessionOnly: a token gets a 403 that says why, not a confusing 401.
	router.Get("/scopes", mw.Authenticated, controller.Scopes)
	router.Get("/", mw.Authenticated, middleware.SessionOnly(), controller.List)
	router.Post("/", mw.Authenticated, middleware.SessionOnly(), controller.Create)
	router.Post("/:tokenID/revoke", mw.Authenticated, middleware.SessionOnly(), controller.Revoke)
	router.Delete("/:tokenID", mw.Authenticated, middleware.SessionOnly(), controller.Delete)
}
