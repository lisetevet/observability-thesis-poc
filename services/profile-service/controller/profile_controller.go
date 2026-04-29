package controller

import (
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"profile-service/model"
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

	tr := otel.Tracer("profile-service")
	reqCtx, span := tr.Start(ctx.Request.Context(), "ProfileController.GetProfile")
	span.SetAttributes(attribute.String("app.uuid", uuid))
	defer span.End()

	if uuid == "" {
		log.Printf("missing uuid path parameter")
		span.SetStatus(codes.Error, "missing uuid")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}

	p, ok, err := c.svc.GetProfile(reqCtx, uuid)
	if err != nil {
		log.Printf("GetProfile failed (uuid=%s): %v", uuid, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "repository error"})
		return
	}

	if !ok {
		log.Printf("profile not found (uuid=%s)", uuid)
		span.SetStatus(codes.Error, "profile not found")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	ctx.JSON(http.StatusOK, p)
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

	query := model.UsersLookupQuery{
		DelayMs: ctx.Query("usersDelayMs"),
		Fail:    ctx.Query("usersFail"),
	}
	query.SetDefaults()

	p, ok, err := c.svc.GetProfileByUsername(reqCtx, username, query)
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
