package controller

import (
	"net/http"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

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

    tr := otel.Tracer("profile-service")
    reqCtx, span := tr.Start(ctx.Request.Context(), "ProfileController.GetProfileByUsername")
    span.SetAttributes(attribute.String("app.username", username))
    defer span.End()

    if username == "" {
        log.Printf("missing username path parameter")
        span.SetStatus(codes.Error, "missing username")
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
        return
    }

    usersDelayMs := ctx.Query("usersDelayMs")
    usersFail := ctx.Query("usersFail")

    p, ok, err := c.svc.GetProfileByUsername(reqCtx, username, usersDelayMs, usersFail)
    if err != nil {
        log.Printf("GetProfileByUsername failed (username=%s): %v", username, err)
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        ctx.JSON(http.StatusBadGateway, gin.H{"error": "profile-service failed"})
        return
    }

    if !ok {
        log.Printf("profile not found (username=%s)", username)
        span.SetStatus(codes.Error, "profile not found")
        ctx.JSON(http.StatusNotFound, gin.H{
            "error":    "no profile found for user",
            "username": username,
        })
        return
    }

    ctx.JSON(http.StatusOK, p)
}