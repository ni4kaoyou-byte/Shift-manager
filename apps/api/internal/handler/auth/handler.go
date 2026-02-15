package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	usecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/auth"
)

type Handler struct {
	usecase *usecase.Usecase
}

func NewHandler(usecase *usecase.Usecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) RegisterRoutes(apiV1 *gin.RouterGroup) {
	group := apiV1.Group("/auth")
	group.GET("/ping", h.ping)
}

func (h *Handler) ping(c *gin.Context) {
	if err := h.usecase.Ping(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"module": "auth", "status": "ok"})
}
