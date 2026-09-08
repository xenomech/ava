package middleware

import (
	"ava/api/internal/model"
	apitokensvc "ava/api/internal/services/apitoken"
	"ava/api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// scopesLocal holds the scopes of the token behind this request; absent for a session.
const scopesLocal = "scopes"

// ValidateAPIToken authenticates a token and leaves the request looking like any other, plus its scopes.
func ValidateAPIToken(service apitokensvc.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authenticated, err := service.Authenticate(c.Context(), extractBearer(c))
		if err != nil {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Invalid or expired token")
		}

		c.Locals("userID", authenticated.Token.UserID)
		c.Locals("tenantID", authenticated.Token.TenantID)
		c.Locals("role", authenticated.Role)
		c.Locals("user", authenticated.User)
		c.Locals(scopesLocal, authenticated.Token.Scopes)

		return c.Next()
	}
}

// Authenticated picks session or token by the shape of the credential, never by trying each in turn.
func Authenticated(
	sessions fiber.Handler,
	tokens fiber.Handler,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if _, _, ok := model.SplitToken(extractBearer(c)); ok {
			return tokens(c)
		}

		return sessions(c)
	}
}

// RequireScope refuses a token lacking the permission; a session holds no scopes and passes untouched.
func RequireScope(required model.Scope) fiber.Handler {
	return func(c *fiber.Ctx) error {
		held, isToken := c.Locals(scopesLocal).(model.Scopes)
		if !isToken {
			return c.Next()
		}

		for _, scope := range held {
			if scope == string(required) {
				return c.Next()
			}
		}

		return response.Send(
			c,
			fiber.StatusForbidden,
			nil,
			"This token does not carry the "+string(required)+" scope",
		)
	}
}

// SessionOnly refuses a personal access token, for things a machine client must never do.
func SessionOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if _, isToken := c.Locals(scopesLocal).(model.Scopes); isToken {
			return response.Send(
				c,
				fiber.StatusForbidden,
				nil,
				"This action needs a signed-in session, not an API token",
			)
		}

		return c.Next()
	}
}
