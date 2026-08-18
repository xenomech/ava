package routes

import (
	flowctrl "ava/api/internal/controller/flow"
	"ava/api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func flowRoutes(router fiber.Router, controller *flowctrl.Controller, mw *middleware.Middleware) {
	router.Use(mw.ValidateAccessToken)

	router.Get("/:type", controller.Get)
	router.Put("/:type/steps/:stepId", controller.SubmitStep)
	router.Post("/:type/back", controller.GoBack)
	router.Post("/:type/skip", controller.SkipStep)
}
