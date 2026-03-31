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

	// pass-through injection params for experiments
	usersDelayMs := ctx.Query("usersDelayMs")
	usersFail := ctx.Query("usersFail")
	profileDelayMs := ctx.Query("profileDelayMs")
	profileFail := ctx.Query("profileFail")

	status, contentType, body, err := c.orch.FetchProfileByUsername(username, usersDelayMs, usersFail, profileDelayMs, profileFail)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.Data(status, contentType, body)
}