package member_http

import (
	"app/app/delivery/http/middleware"
	"app/domain"

	"github.com/gin-gonic/gin"
)

type routeMember struct {
	Usecase    domain.MemberAppUsecase
	Route      *gin.RouterGroup
	Middleware middleware.AppMiddleware
}

func NewMemberRouteHandler(usecase domain.MemberAppUsecase, ginEngine *gin.Engine, middleware middleware.AppMiddleware) {
	handler := &routeMember{
		Usecase:    usecase,
		Route:      ginEngine.Group("/member"),
		Middleware: middleware,
	}

	handler.handleAuthRoute("/auth")
	handler.handleVotingRoute("/votings")
	handler.handleCandidateRoute("/candidates")
	handler.handlePurchaseRoute("/purchases")
	handler.handleSeasonRoute("/seasons")
	handler.handleTicketRoute("/tickets")
	handler.handleSeriesRoute("/series")
	handler.handleTicketPurchaseRoute("/ticket-purchases")
}
