package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"slices"
)

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

// ValidScope reports whether the string names a scope this server knows.
func ValidScope(candidate string) bool {
	return slices.Contains(AllScopes, Scope(candidate))
}

// Scopes is the granted set, stored as jsonb so it needs no array driver.
type Scopes []string

// Has reports whether this set carries the permission.
func (scopes Scopes) Has(wanted Scope) bool {
	return slices.Contains(scopes, string(wanted))
}

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
