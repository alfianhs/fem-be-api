package superadmin_http

import (
	"app/app/delivery/http/middleware"
	"app/domain"

	"github.com/gin-gonic/gin"
)

type routeSuperadmin struct {
	Usecase    domain.SuperadminAppUsecase
	Route      *gin.RouterGroup
	Middleware middleware.AppMiddleware
}

func NewSuperadminRouteHandler(usecase domain.SuperadminAppUsecase, ginEngine *gin.Engine, middleware middleware.AppMiddleware) {
	handler := &routeSuperadmin{
		Usecase:    usecase,
		Route:      ginEngine.Group("/superadmin"),
		Middleware: middleware,
	}

	handler.handleAuthRoute("/auth")
	handler.handleSeasonRoute("/seasons")
	handler.handleVenueRoute("/venues")
	handler.handleTeamRoute("/teams")
	handler.handlePlayerRoute("/players")
	handler.handleSeasonTeamRoute("/season-teams")
	handler.handleSeasonTeamPlayerRoute("/season-team-players")
	handler.handleSeriesRoute("/series")
}
