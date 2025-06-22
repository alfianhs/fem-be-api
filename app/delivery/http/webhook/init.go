package webhook_http

import (
	"app/app/delivery/http/middleware"
	"app/domain"

	"github.com/gin-gonic/gin"
)

type routeWebhook struct {
	Usecase    domain.WebhookAppUsecase
	Route      *gin.RouterGroup
	Middleware middleware.AppMiddleware
}

func NewWebhookRouteHandler(usecase domain.WebhookAppUsecase, ginEngine *gin.Engine, middleware middleware.AppMiddleware) {
	handler := &routeWebhook{
		Usecase:    usecase,
		Route:      ginEngine.Group("/webhook"),
		Middleware: middleware,
	}

	handler.handleXenditRoute("/xendit")
}
