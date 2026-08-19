package response

import (
	"ava/api/pkg/serrors"

	"github.com/gofiber/fiber/v2"
)

const (
	CodeBadRequest      = "bad_request"
	CodeUnauthorized    = "unauthorized"
	CodeForbidden       = "forbidden"
	CodeNotFound        = "not_found"
	CodeConflict        = "conflict"
	CodeValidationError = "validation_error"
	CodeRateLimited     = "rate_limited"
	CodeInternalError   = "internal_error"
)

type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func codeForStatus(status int) string {
	switch status {
	case fiber.StatusBadRequest:
		return CodeBadRequest
	case fiber.StatusUnauthorized:
		return CodeUnauthorized
	case fiber.StatusForbidden:
		return CodeForbidden
	case fiber.StatusNotFound:
		return CodeNotFound
	case fiber.StatusConflict:
		return CodeConflict
	case fiber.StatusUnprocessableEntity:
		return CodeValidationError
	case fiber.StatusTooManyRequests:
		return CodeRateLimited
	default:
		return CodeInternalError
	}
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
			Code:    codeForStatus(status),
			Message: errorMessage,
		},
	})
}

func SendWithCode(c *fiber.Ctx, status int, data any, code, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": code == "",
		"data":    data,
		"error": ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func SendError(c *fiber.Ctx, status int, err error) error {
	return SendErrorWithData(c, status, nil, err)
}

func SendErrorWithData(c *fiber.Ctx, status int, data any, err error) error {
	code := serrors.CodeOf(err)
	if code == "" {
		code = codeForStatus(status)
	}

	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"data":    data,
		"error": ErrorDetail{
			Code:    code,
			Message: err.Error(),
		},
	})
}

func SendValidation(c *fiber.Ctx, message string, fields map[string]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"error": ErrorDetail{
			Code:    CodeValidationError,
			Message: message,
			Details: fields,
		},
	})
}
