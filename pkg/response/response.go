package response

import "github.com/gofiber/fiber/v2"

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Send(c *fiber.Ctx, status int, data any, errorMessage string) error {
	if errorMessage == "" {
		return c.Status(status).JSON(fiber.Map{
			"success": true,
			"data":    data,
			"error":   ErrorDetail{},
		})
	}

	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"data":    data,
		"error": ErrorDetail{
			Message: errorMessage,
		},
	})
}
