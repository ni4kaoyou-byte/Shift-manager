package auditlog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	commonhandler "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/handler"
	usecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/audit_log"
)

type Handler struct {
	service *usecase.Service
}

func NewHandler(service *usecase.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(apiV1 *gin.RouterGroup) {
	group := apiV1.Group("/audit-logs")
	group.GET("/ping", h.ping)
}

func (h *Handler) ping(c *gin.Context) {
	if err := h.service.Ping(c.Request.Context()); err != nil {
		commonhandler.AbortInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"module": "audit_log", "status": "ok"})
}
