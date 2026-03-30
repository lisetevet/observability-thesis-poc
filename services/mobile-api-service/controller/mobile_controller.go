package controller

import (
	"net/http"

	"mobile-api-service/service"

	"github.com/gin-gonic/gin"
)

type MobileController struct {
	orch *service.Orchestrator
}

func NewMobileController(orch *service.Orchestrator) *MobileController {
	return &MobileController{orch: orch}
}

func (c *MobileController) Health(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
}

func (c *MobileController) GetProfile(ctx *gin.Context) {
	username := ctx.Param("username")

	status, contentType, body, err := c.orch.FetchProfileByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.Data(status, contentType, body)
}