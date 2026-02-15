package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRecoverResponseFormat(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID(), Recover())
	router.GET("/panic", func(c *gin.Context) {
		panic("unexpected")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	expected := "{\"error\":{\"code\":\"internal_error\",\"message\":\"internal server error\"}}"
	if rec.Body.String() != expected {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
