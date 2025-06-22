package webhook_http

import (
	"app/domain/request"
	"app/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *routeWebhook) handleXenditRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.POST("/snap", h.Middleware.AuthXendit(), h.XenditSnapWebhook)
}

// XenditSnapWebhook godoc
//
// @Summary      Xendit Snap Webhook
// @Description  Handle Xendit Snap Webhook
// @Tags         Webhook
// @Accept       json
// @Produce      json
//
//	@Param			payload	body	request.SnapWebhookRequest	true	"Xendit Snap Webhook Payload"
//
// @Success      200 {object} helpers.Response
// @Router       /webhook/xendit/snap [post]
func (h *routeWebhook) XenditSnapWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	var payload request.SnapWebhookRequest
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, err.Error(), nil, nil))
		return
	}

	response := h.Usecase.HandleXenditWebhook(ctx, payload)
	c.JSON(response.Status, response)
}
