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