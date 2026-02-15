package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/middleware"
	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	Error ErrorObject `json:"error"`
}

type ErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func AbortWithError(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, ErrorResponse{
		Error: ErrorObject{
			Code:    code,
			Message: message,
		},
	})
}

func AbortInternalServerError(c *gin.Context, err error) {
	log.Error().
		Err(err).
		Str("request_id", middleware.GetRequestID(c)).
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Msg("request failed")

	AbortWithError(c, http.StatusInternalServerError, "internal_error", "internal server error")
}
