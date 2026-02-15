package server

import (
	"github.com/gin-gonic/gin"
	assignment "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/module/assignment"
	auditlog "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/module/audit_log"
	auth "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/module/auth"
	availability "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/module/availability"
	changerequest "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/module/change_request"
	membership "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/module/membership"
	period "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/module/period"
)

type routeRegistrar interface {
	RegisterRoutes(apiV1 *gin.RouterGroup)
}

var moduleFactories = []func() routeRegistrar{
	func() routeRegistrar { return auth.NewHandler() },
	func() routeRegistrar { return membership.NewHandler() },
	func() routeRegistrar { return period.NewHandler() },
	func() routeRegistrar { return availability.NewHandler() },
	func() routeRegistrar { return assignment.NewHandler() },
	func() routeRegistrar { return changerequest.NewHandler() },
	func() routeRegistrar { return auditlog.NewHandler() },
}

func newRouteRegistrars() []routeRegistrar {
	registrars := make([]routeRegistrar, 0, len(moduleFactories))
	for _, factory := range moduleFactories {
		registrars = append(registrars, factory())
	}
	return registrars
}
