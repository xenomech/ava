package auth

import (
	"time"

	"ava/api/config"
	"ava/api/internal/dto"
	"ava/api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

const (
	AccessCookie  = middleware.AccessCookie
	RefreshCookie = "ava_refresh"

	accessCookiePath  = "/api/v1"
	refreshCookiePath = "/api/v1/auth"
)

func setSessionCookies(ctx *fiber.Ctx, tokens *dto.TokenResponse) {
	if tokens == nil {
		return
	}

	cfg := config.GetConfig()
	secure := cfg.ServerEnv != "local"

	ctx.Cookie(&fiber.Cookie{
		Name:     AccessCookie,
		Value:    tokens.AccessToken,
		Path:     accessCookiePath,
		Domain:   cfg.CookieDomain,
		MaxAge:   int(cfg.JwtAccessExpiry.Seconds()),
		Expires:  time.Now().Add(cfg.JwtAccessExpiry),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     RefreshCookie,
		Value:    tokens.RefreshToken,
		Path:     refreshCookiePath,
		Domain:   cfg.CookieDomain,
		MaxAge:   int(cfg.JwtRefreshExpiry.Seconds()),
		Expires:  time.Now().Add(cfg.JwtRefreshExpiry),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}

func clearSessionCookies(ctx *fiber.Ctx) {
	cfg := config.GetConfig()
	secure := cfg.ServerEnv != "local"

	for name, path := range map[string]string{
		AccessCookie:  accessCookiePath,
		RefreshCookie: refreshCookiePath,
	} {
		ctx.Cookie(&fiber.Cookie{
			Name:     name,
			Value:    "",
			Path:     path,
			Domain:   cfg.CookieDomain,
			MaxAge:   -1,
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
			Secure:   secure,
			SameSite: fiber.CookieSameSiteLaxMode,
		})
	}
}

func redactTokens(tokens *dto.TokenResponse) *dto.TokenResponse {
	if tokens == nil {
		return nil
	}

	return &dto.TokenResponse{ExpiresIn: tokens.ExpiresIn}
}
