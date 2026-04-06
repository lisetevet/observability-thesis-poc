package controller

import (
	"net/http"

	"profile-service/service"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	svc *service.ProfileService
}

func NewProfileController(svc *service.ProfileService) *ProfileController {
	return &ProfileController{svc: svc}
}

func (c *ProfileController) Health(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
}

func (c *ProfileController) GetProfile(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	p, ok, err := c.svc.GetProfile(ctx.Request.Context(), uuid)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "repository error"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "no profile found for user",
			"uuid":  uuid,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"uuid":          p.UUID,
		"name":          p.Name,
		"surname":       p.Surname,
		"email":         p.Email,
		"personal_code": p.PersonalCode,
	})
}

func (c *ProfileController) GetProfileByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	// allow experiment pass-through to users-service from profile-service
	usersDelayMs := ctx.Query("usersDelayMs")
	usersFail := ctx.Query("usersFail")

	p, ok, err := c.svc.GetProfileByUsernameDBFirst(ctx.Request.Context(), username, usersDelayMs, usersFail)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "profile-service failed"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":    "no profile found for user",
			"username": username,
		})
		return
	}

	ctx.JSON(http.StatusOK, p)
}