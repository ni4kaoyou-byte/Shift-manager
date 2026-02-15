package config

import (
	"strings"
	"testing"
)

func TestLoadFromLookupSuccess(t *testing.T) {
	lookup := func(key string) (string, bool) {
		values := map[string]string{
			"APP_ENV":             "test",
			"PORT":                "18080",
			"DATABASE_URL":        "postgres://user:pass@localhost:5432/shift_manager",
			"SUPABASE_URL":        "https://example.supabase.co",
			"SUPABASE_JWT_SECRET": "secret",
		}
		value, ok := values[key]
		return value, ok
	}

	cfg, err := loadFromLookup(lookup)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.AppEnv != "test" {
		t.Fatalf("expected AppEnv to be test, got %s", cfg.AppEnv)
	}
	if cfg.Port != "18080" {
		t.Fatalf("expected Port to be 18080, got %s", cfg.Port)
	}
}

func TestLoadFromLookupMissingRequiredEnv(t *testing.T) {
	lookup := func(key string) (string, bool) {
		values := map[string]string{
			"DATABASE_URL": "postgres://user:pass@localhost:5432/shift_manager",
		}
		value, ok := values[key]
		return value, ok
	}

	_, err := loadFromLookup(lookup)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	errorMessage := err.Error()
	if !strings.Contains(errorMessage, "SUPABASE_URL") {
		t.Fatalf("expected missing SUPABASE_URL, got %s", errorMessage)
	}
	if !strings.Contains(errorMessage, "SUPABASE_JWT_SECRET") {
		t.Fatalf("expected missing SUPABASE_JWT_SECRET, got %s", errorMessage)
	}
}
