package jwt

import (
	"time"

	"ava/api/internal/model"
	"ava/api/pkg/serrors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
	DeviceToken  TokenType = "device"
)

var (
	ErrInvalidJWT = serrors.New("invalid token")
	ErrExpiredJWT = serrors.New("token has expired")
)

type Claims struct {
	jwt.RegisteredClaims
	UserID    uuid.UUID        `json:"user_id"`
	DeviceID  uuid.UUID        `json:"device_id,omitempty"`
	TenantID  uuid.UUID        `json:"tenant_id"`
	Role      model.TenantRole `json:"role"`
	SessionID uuid.UUID        `json:"session_id"`
	Email     string           `json:"email"`
	TokenType TokenType        `json:"type"`
}

func (tm *jwtTokenManager) GenerateToken(user *model.User, tenantID uuid.UUID, role model.TenantRole, sessionID uuid.UUID, tokenType TokenType, rid string) (string, error) {
	var expiry time.Duration

	switch tokenType {
	case AccessToken:
		expiry = tm.accessExpiry
	case RefreshToken:
		expiry = tm.refreshExpiry
	case DeviceToken:
		expiry = tm.deviceExpiry
	default:
		return "", serrors.New("invalid token type")
	}

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ava",
			Subject:   user.ID.String(),
		},
		UserID:    user.ID,
		TenantID:  tenantID,
		Role:      role,
		SessionID: sessionID,
		Email:     user.Email,
		TokenType: tokenType,
	}

	if rid != "" {
		claims.ID = rid
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tm.secretKey))
}

func (tm *jwtTokenManager) GenerateDeviceToken(device *model.Device) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.deviceExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ava",
			Subject:   device.ID.String(),
		},
		DeviceID:  device.ID,
		TenantID:  device.TenantID,
		TokenType: DeviceToken,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tm.secretKey))
}

func (tm *jwtTokenManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidJWT
		}

		return []byte(tm.secretKey), nil
	})
	if err != nil {
		if serrors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredJWT
		}

		return nil, ErrInvalidJWT
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidJWT
	}

	return claims, nil
}

func (tm *jwtTokenManager) GetAccessExpiry() time.Duration {
	return tm.accessExpiry
}

func (tm *jwtTokenManager) GetRefreshExpiry() time.Duration {
	return tm.refreshExpiry
}

func (tm *jwtTokenManager) GetDeviceExpiry() time.Duration {
	return tm.deviceExpiry
}
