package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AppEnv            string
	Port              string
	DatabaseURL       string
	SupabaseURL       string
	SupabaseJWTSecret string
}

func Load() (Config, error) {
	return loadFromLookup(os.LookupEnv)
}

func loadFromLookup(lookup func(string) (string, bool)) (Config, error) {
	cfg := Config{
		AppEnv:            getOrDefault(lookup, "APP_ENV", "development"),
		Port:              getOrDefault(lookup, "PORT", "8080"),
		DatabaseURL:       strings.TrimSpace(getRequired(lookup, "DATABASE_URL")),
		SupabaseURL:       strings.TrimSpace(getRequired(lookup, "SUPABASE_URL")),
		SupabaseJWTSecret: strings.TrimSpace(getRequired(lookup, "SUPABASE_JWT_SECRET")),
	}

	missing := make([]string, 0)
	if cfg.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if cfg.SupabaseURL == "" {
		missing = append(missing, "SUPABASE_URL")
	}
	if cfg.SupabaseJWTSecret == "" {
		missing = append(missing, "SUPABASE_JWT_SECRET")
	}

	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

func getRequired(lookup func(string) (string, bool), key string) string {
	value, _ := lookup(key)
	return value
}

func getOrDefault(lookup func(string) (string, bool), key, fallback string) string {
	value, ok := lookup(key)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
