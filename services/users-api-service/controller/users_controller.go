package controller

import (
	"net/http"

	"users-api-service/service"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	svc *service.UserService
}

func NewUsersController(svc *service.UserService) *UsersController {
	return &UsersController{svc: svc}
}

func (c *UsersController) Health(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
}

func (c *UsersController) GetUserUUID(ctx *gin.Context) {
	username := ctx.Param("username")

	uuid, ok, err := c.svc.GetUUID(ctx.Request.Context(), username)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "repository error"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":    "user not found",
			"username": username,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"username": username,
		"uuid":     uuid,
	})
}

func (c *UsersController) GetUserProfileSeed(ctx *gin.Context) {
	username := ctx.Param("username")

	u, ok, err := c.svc.GetUser(ctx.Request.Context(), username)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "repository error"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":    "user not found",
			"username": username,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"uuid":          u.UUID,
		"name":          u.Name,
		"surname":       u.Surname,
		"email":         u.Email,
		"personal_code": u.PersonalCode,
	})
}