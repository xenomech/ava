package config_test

import (
	"errors"
	"strings"
	"testing"

	"ava/api/config"
)

func TestAnEmptySecretIsRefusedInEveryEnvironment(t *testing.T) {
	for _, env := range []string{"local", "staging", "production"} {
		t.Run(env, func(t *testing.T) {
			cfg := &config.Config{ServerEnv: env}

			if err := cfg.Validate(); !errors.Is(err, config.ErrMissingJWTSecret) {
				t.Fatalf("Validate() = %v, want ErrMissingJWTSecret", err)
			}

			// Whitespace is not a secret, and viper hands back whatever the environment held.
			cfg.JwtSecretKey = "   \t "
			if err := cfg.Validate(); !errors.Is(err, config.ErrMissingJWTSecret) {
				t.Fatalf("Validate() with blank secret = %v, want ErrMissingJWTSecret", err)
			}
		})
	}
}

func TestAShortSecretPassesLocallyAndFailsEverywhereElse(t *testing.T) {
	const short = "dev"

	if err := (&config.Config{ServerEnv: "local", JwtSecretKey: short}).Validate(); err != nil {
		t.Fatalf("local rejected the throwaway secret the docs hand out: %v", err)
	}

	err := (&config.Config{ServerEnv: "production", JwtSecretKey: short}).Validate()
	if err == nil {
		t.Fatal("production accepted a 3-character signing key")
	}

	if !strings.Contains(err.Error(), "openssl") {
		t.Errorf("error does not say how to fix it: %v", err)
	}
}

func TestASecretOfTheHashLengthIsAccepted(t *testing.T) {
	// 32 bytes, the HS256 output size and what `openssl rand -hex 32` produces as 64 hex characters.
	cfg := &config.Config{ServerEnv: "production", JwtSecretKey: strings.Repeat("a", 32)}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
}
