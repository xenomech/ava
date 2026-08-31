package model_test

import (
	"strings"
	"testing"
	"time"

	"ava/api/internal/model"

	"github.com/google/uuid"
)

func mint(t *testing.T, scopes []model.Scope, expiresAt *time.Time) (*model.APIToken, string) {
	t.Helper()

	token, plaintext, err := model.NewAPIToken(uuid.New(), uuid.New(), "Shortcuts", scopes, expiresAt)
	if err != nil {
		t.Fatalf("minting failed: %v", err)
	}

	return token, plaintext
}

func TestAPlaintextTokenMatchesTheStoredHashAndNothingElse(t *testing.T) {
	token, plaintext := mint(t, []model.Scope{model.ScopeDevicesWrite}, nil)

	_, secret, ok := model.SplitToken(plaintext)
	if !ok {
		t.Fatalf("could not split %q", plaintext)
	}

	if !token.MatchesSecret(secret) {
		t.Fatal("the secret it was created with did not match")
	}

	if token.MatchesSecret(secret + "x") {
		t.Fatal("a different secret matched")
	}
}

func TestTheSecretIsNeverStoredOnTheToken(t *testing.T) {
	token, plaintext := mint(t, []model.Scope{model.ScopeDevicesRead}, nil)

	_, secret, _ := model.SplitToken(plaintext)

	if strings.Contains(token.Hash, secret) {
		t.Fatal("the hash contains the secret verbatim")
	}

	if token.Lookup == secret {
		t.Fatal("the lookup half is the secret half")
	}
}

func TestTwoTokensNeverShareALookup(t *testing.T) {
	seen := make(map[string]bool, 200)

	for range 200 {
		token, _ := mint(t, []model.Scope{model.ScopeDevicesRead}, nil)

		if seen[token.Lookup] {
			t.Fatalf("lookup %q was issued twice", token.Lookup)
		}

		seen[token.Lookup] = true
	}
}

func TestSplitTokenRejectsAnythingNotShapedLikeAToken(t *testing.T) {
	for _, presented := range []string{
		"",
		"nonsense",
		"eyJhbGciOiJIUzI1NiJ9.e30.signature", // a session JWT
		model.TokenPrefix,
		model.TokenPrefix + "onlylookup",
		model.TokenPrefix + "_secretbutnolookup",
	} {
		if _, _, ok := model.SplitToken(presented); ok {
			t.Fatalf("%q was accepted as a token", presented)
		}
	}
}

func TestARevokedTokenIsNotActive(t *testing.T) {
	token, _ := mint(t, []model.Scope{model.ScopeDevicesRead}, nil)

	if !token.Active() {
		t.Fatal("a fresh token should be active")
	}

	now := time.Now()
	token.RevokedAt = &now

	if token.Active() {
		t.Fatal("a revoked token is still active")
	}
}

func TestAnExpiredTokenIsNotActiveButAFutureOneIs(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	expired, _ := mint(t, []model.Scope{model.ScopeDevicesRead}, &past)

	if expired.Active() {
		t.Fatal("a token past its expiry is still active")
	}

	future := time.Now().Add(time.Hour)
	live, _ := mint(t, []model.Scope{model.ScopeDevicesRead}, &future)

	if !live.Active() {
		t.Fatal("a token inside its expiry is not active")
	}
}

func TestATokenOnlyHoldsTheScopesItWasGiven(t *testing.T) {
	token, _ := mint(t, []model.Scope{model.ScopeDevicesRead, model.ScopeScenesRead}, nil)

	if !token.HasScope(model.ScopeDevicesRead) || !token.HasScope(model.ScopeScenesRead) {
		t.Fatal("a granted scope is missing")
	}

	if token.HasScope(model.ScopeDevicesWrite) {
		t.Fatal("a scope that was never granted is held")
	}
}

func TestValidScopeAcceptsOnlyKnownScopes(t *testing.T) {
	for _, scope := range model.AllScopes {
		if !model.ValidScope(string(scope)) {
			t.Fatalf("%q is in AllScopes but not valid", scope)
		}
	}

	for _, candidate := range []string{"", "devices", "devices:*", "admin", "DEVICES:READ"} {
		if model.ValidScope(candidate) {
			t.Fatalf("%q was accepted as a scope", candidate)
		}
	}
}
