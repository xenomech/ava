package middleware

import (
	"ava/api/internal/model"
	apitokensvc "ava/api/internal/services/apitoken"
	"ava/api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// scopesLocal holds the scopes of the token behind this request; absent for a session.
const scopesLocal = "scopes"

/*
ValidateAPIToken authenticates a personal access token and leaves the request looking like any
other: same locals, same tenant, same role. Downstream handlers cannot tell the difference, which
is the point — a token is a way in, not a second kind of caller.

It sets `scopes` so RequireScope can narrow what this particular token may do. A session sets no
scopes at all and is therefore unrestricted.
*/
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

/*
Authenticated accepts either a session or a personal access token.

Which one is decided by the shape of the presented credential, not by trying each in turn: a token
carries a recognisable prefix, so a malformed session JWT is never retried as a token.
*/
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

/*
RequireScope refuses a token that was not granted this permission.

A session holds no scopes and passes untouched — the browser is already limited by the member's
role, and re-stating that here would mean two systems disagreeing about the same question.
*/
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
