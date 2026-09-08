package model

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TokenPrefix marks a personal access token so the middleware can tell one from a session JWT.
const TokenPrefix = "ava_pat_"

// APIToken is a machine client's credential: issued once, scoped, and revoked rather than rotated.
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
	return token.Scopes.Has(wanted)
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
	// Hex, not base64url, whose alphabet contains the "_" that separates the two halves.
	lookup, err := randomHex(8)
	if err != nil {
		return nil, "", err
	}

	secret, err := randomSecret(24)
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

// randomHex is used where the value must not contain the "_" delimiter.
func randomHex(bytes int) (string, error) {
	buf, err := randomBytes(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}

// randomSecret packs more entropy per character; it may contain "_" safely, being the last half.
func randomSecret(bytes int) (string, error) {
	buf, err := randomBytes(bytes)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func randomBytes(count int) ([]byte, error) {
	buf := make([]byte, count)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}

	return buf, nil
}
