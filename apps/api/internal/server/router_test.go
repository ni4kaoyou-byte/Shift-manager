package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthzReturnsOK(t *testing.T) {
	t.Parallel()

	router := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("X-Request-Id") == "" {
		t.Fatal("expected X-Request-Id header to be set")
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status=ok, got %q", body["status"])
	}
}

func TestAPIV1BaseRoute(t *testing.T) {
	t.Parallel()

	router := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestModuleRoutes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		path   string
		module string
	}{
		{path: "/api/v1/auth/ping", module: "auth"},
		{path: "/api/v1/membership/ping", module: "membership"},
		{path: "/api/v1/period/ping", module: "period"},
		{path: "/api/v1/availability/ping", module: "availability"},
		{path: "/api/v1/assignment/ping", module: "assignment"},
		{path: "/api/v1/change-requests/ping", module: "change_request"},
		{path: "/api/v1/audit-logs/ping", module: "audit_log"},
	}

	router := NewRouter()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
			}

			var body map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if body["module"] != tc.module {
				t.Fatalf("expected module %q, got %q", tc.module, body["module"])
			}
		})
	}
}
