package model

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql/driver"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TokenPrefix marks a personal access token so the middleware can tell one from a session JWT.
const TokenPrefix = "ava_pat_"

// Scope is one permission a token may carry; a session always behaves as if it held every scope.
type Scope string

const (
	ScopeDevicesRead  Scope = "devices:read"
	ScopeDevicesWrite Scope = "devices:write"
	ScopeRoomsRead    Scope = "rooms:read"
	ScopeRoomsWrite   Scope = "rooms:write"
	ScopeScenesRead   Scope = "scenes:read"
	ScopeScenesWrite  Scope = "scenes:write"
	ScopeHubsRead     Scope = "hubs:read"
	ScopeHubsWrite    Scope = "hubs:write"
)

// AllScopes is every scope a token may be granted, in the order they are shown.
var AllScopes = []Scope{
	ScopeDevicesRead,
	ScopeDevicesWrite,
	ScopeRoomsRead,
	ScopeRoomsWrite,
	ScopeScenesRead,
	ScopeScenesWrite,
	ScopeHubsRead,
	ScopeHubsWrite,
}

// Scopes is the granted set, stored as jsonb so it needs no array driver.
type Scopes []string

func (scopes *Scopes) Scan(value any) error {
	if value == nil {
		*scopes = nil

		return nil
	}

	switch raw := value.(type) {
	case []byte:
		return json.Unmarshal(raw, scopes)
	case string:
		return json.Unmarshal([]byte(raw), scopes)
	default:
		return errors.New("scopes: unsupported database value")
	}
}

func (scopes Scopes) Value() (driver.Value, error) {
	if scopes == nil {
		return "[]", nil
	}

	encoded, err := json.Marshal([]string(scopes))

	return string(encoded), err
}

// ValidScope reports whether the string names a scope this server knows.
func ValidScope(candidate string) bool {
	for _, scope := range AllScopes {
		if string(scope) == candidate {
			return true
		}
	}

	return false
}

/*
APIToken is a long-lived credential for a machine client: a Shortcut, a script, a home panel.

Deliberately not a session. Sessions rotate their refresh token on every use, which is right for a
browser and useless for a client that cannot store the new one. A token here is created once,
carries only the scopes it was given, and is revoked rather than rotated.
*/
type APIToken struct {
	BaseModel
	TenantID uuid.UUID `gorm:"type:uuid;not null;index:idx_api_token_tenant_user" json:"tenant_id"`
	Tenant   *Tenant   `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index:idx_api_token_tenant_user" json:"user_id"`
	User     *User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`

	Name string `gorm:"type:varchar(80);not null" json:"name"`
	// Lookup is the token's public half: indexed, safe to log, and what finds the row.
	Lookup string `gorm:"type:varchar(32);not null;uniqueIndex" json:"-"`
	// Hash is the SHA-256 of the secret half. The secret itself is shown once and never stored.
	Hash   string `gorm:"type:varchar(64);not null" json:"-"`
	Scopes Scopes `gorm:"type:jsonb;not null;default:'[]'" json:"scopes"`

	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
}

func (token *APIToken) BeforeCreate(tx *gorm.DB) error {
	return token.BaseModel.BeforeCreate(tx)
}

// Active reports whether the token may still authenticate a request.
func (token *APIToken) Active() bool {
	if token.RevokedAt != nil {
		return false
	}

	return token.ExpiresAt == nil || time.Now().Before(*token.ExpiresAt)
}

// HasScope reports whether the token was granted this permission.
func (token *APIToken) HasScope(wanted Scope) bool {
	for _, held := range token.Scopes {
		if held == string(wanted) {
			return true
		}
	}

	return false
}

// MatchesSecret compares in constant time, so a wrong guess cannot be timed against a right one.
func (token *APIToken) MatchesSecret(secret string) bool {
	sum := sha256.Sum256([]byte(secret))
	want := hex.EncodeToString(sum[:])

	return subtle.ConstantTimeCompare([]byte(token.Hash), []byte(want)) == 1
}

// NewAPIToken mints a token, returning it alongside the one and only copy of its plaintext.
func NewAPIToken(
	tenantID, userID uuid.UUID,
	name string,
	scopes []Scope,
	expiresAt *time.Time,
) (*APIToken, string, error) {
	lookup, err := randomString(9)
	if err != nil {
		return nil, "", err
	}

	secret, err := randomString(24)
	if err != nil {
		return nil, "", err
	}

	granted := make(Scopes, 0, len(scopes))
	for _, scope := range scopes {
		granted = append(granted, string(scope))
	}

	sum := sha256.Sum256([]byte(secret))

	token := &APIToken{
		TenantID:  tenantID,
		UserID:    userID,
		Name:      name,
		Lookup:    lookup,
		Hash:      hex.EncodeToString(sum[:]),
		Scopes:    granted,
		ExpiresAt: expiresAt,
	}

	return token, TokenPrefix + lookup + "_" + secret, nil
}

// SplitToken pulls the lookup and secret halves out of a presented token.
func SplitToken(presented string) (lookup, secret string, ok bool) {
	if !strings.HasPrefix(presented, TokenPrefix) {
		return "", "", false
	}

	lookup, secret, found := strings.Cut(strings.TrimPrefix(presented, TokenPrefix), "_")
	if !found || lookup == "" || secret == "" {
		return "", "", false
	}

	return lookup, secret, true
}

func randomString(bytes int) (string, error) {
	buf := make([]byte, bytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}
