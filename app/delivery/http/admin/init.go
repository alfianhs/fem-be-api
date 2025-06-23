package admin_http

import (
	"app/app/delivery/http/middleware"
	"app/domain"

	"github.com/gin-gonic/gin"
)

type routeAdmin struct {
	Usecase    domain.AdminAppUsecase
	Route      *gin.RouterGroup
	Middleware middleware.AppMiddleware
}

func NewAdminRouteHandler(usecase domain.AdminAppUsecase, ginEngine *gin.Engine, middleware middleware.AppMiddleware) {
	handler := &routeAdmin{
		Usecase:    usecase,
		Route:      ginEngine.Group("/admin"),
		Middleware: middleware,
	}

	handler.handleAuthRoute("/auth")
	handler.handleTicketPurchaseRoute("/ticket-purchases")
}
