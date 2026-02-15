package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAbortWithError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	AbortWithError(c, http.StatusBadRequest, "invalid_request", "invalid request")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	expected := "{\"error\":{\"code\":\"invalid_request\",\"message\":\"invalid request\"}}"
	if rec.Body.String() != expected {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestAbortInternalServerError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	AbortInternalServerError(c, errors.New("boom"))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	expected := "{\"error\":{\"code\":\"internal_error\",\"message\":\"internal server error\"}}"
	if rec.Body.String() != expected {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
