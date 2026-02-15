package server

import (
	"github.com/gin-gonic/gin"
	assignmenthandler "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/handler/assignment"
	auditloghandler "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/handler/audit_log"
	authhandler "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/handler/auth"
	availabilityhandler "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/handler/availability"
	changerequesthandler "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/handler/change_request"
	membershiphandler "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/handler/membership"
	periodhandler "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/handler/period"
	assignmentrepository "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/repository/assignment"
	auditlogrepository "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/repository/audit_log"
	authrepository "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/repository/auth"
	availabilityrepository "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/repository/availability"
	changerequestrepository "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/repository/change_request"
	membershiprepository "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/repository/membership"
	periodrepository "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/repository/period"
	assignmentusecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/assignment"
	auditlogusecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/audit_log"
	authusecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/auth"
	availabilityusecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/availability"
	changerequestusecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/change_request"
	membershipusecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/membership"
	periodusecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/period"
)

type moduleHandlers struct {
	auth          *authhandler.Handler
	membership    *membershiphandler.Handler
	period        *periodhandler.Handler
	availability  *availabilityhandler.Handler
	assignment    *assignmenthandler.Handler
	changeRequest *changerequesthandler.Handler
	auditLog      *auditloghandler.Handler
}

func newModuleHandlers() *moduleHandlers {
	authRepo := authrepository.NewInMemoryRepository()
	authUC := authusecase.NewUsecase(authRepo)

	membershipRepo := membershiprepository.NewInMemoryRepository()
	membershipUC := membershipusecase.NewUsecase(membershipRepo)

	periodRepo := periodrepository.NewInMemoryRepository()
	periodUC := periodusecase.NewUsecase(periodRepo)

	availabilityRepo := availabilityrepository.NewInMemoryRepository()
	availabilityUC := availabilityusecase.NewUsecase(availabilityRepo)

	assignmentRepo := assignmentrepository.NewInMemoryRepository()
	assignmentUC := assignmentusecase.NewUsecase(assignmentRepo)

	changeRequestRepo := changerequestrepository.NewInMemoryRepository()
	changeRequestUC := changerequestusecase.NewUsecase(changeRequestRepo)

	auditLogRepo := auditlogrepository.NewInMemoryRepository()
	auditLogUC := auditlogusecase.NewUsecase(auditLogRepo)

	return &moduleHandlers{
		auth:          authhandler.NewHandler(authUC),
		membership:    membershiphandler.NewHandler(membershipUC),
		period:        periodhandler.NewHandler(periodUC),
		availability:  availabilityhandler.NewHandler(availabilityUC),
		assignment:    assignmenthandler.NewHandler(assignmentUC),
		changeRequest: changerequesthandler.NewHandler(changeRequestUC),
		auditLog:      auditloghandler.NewHandler(auditLogUC),
	}
}

func (h *moduleHandlers) RegisterRoutes(apiV1 *gin.RouterGroup) {
	h.auth.RegisterRoutes(apiV1)
	h.membership.RegisterRoutes(apiV1)
	h.period.RegisterRoutes(apiV1)
	h.availability.RegisterRoutes(apiV1)
	h.assignment.RegisterRoutes(apiV1)
	h.changeRequest.RegisterRoutes(apiV1)
	h.auditLog.RegisterRoutes(apiV1)
}
