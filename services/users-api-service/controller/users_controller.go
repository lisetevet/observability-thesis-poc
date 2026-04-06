package controller

import (
	"net/http"

	"users-api-service/model"
	"users-api-service/service"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

	tr := otel.Tracer("users-api-service")
	reqCtx, span := tr.Start(ctx.Request.Context(), "UsersController.GetUserUUID")
	span.SetAttributes(attribute.String("app.username", username))
	defer span.End()

	uuid, ok, err := c.svc.GetUUID(reqCtx, username)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "service error")
		ctx.JSON(http.StatusBadGateway, model.ErrorResponse{Error: "repository error", Username: username})
		return
	}
	if !ok {
		ctx.JSON(http.StatusNotFound, model.ErrorResponse{Error: "user not found", Username: username})
		return
	}

	ctx.JSON(http.StatusOK, model.UserUUIDResponse{
		Username: username,
		UUID:     uuid,
	})
}

func (c *UsersController) GetUserProfileSeed(ctx *gin.Context) {
	username := ctx.Param("username")

	tr := otel.Tracer("users-api-service")
	reqCtx, span := tr.Start(ctx.Request.Context(), "UsersController.GetUserProfileSeed")
	span.SetAttributes(attribute.String("app.username", username))
	defer span.End()

	u, ok, err := c.svc.GetUser(reqCtx, username)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "service error")
		ctx.JSON(http.StatusBadGateway, model.ErrorResponse{Error: "repository error", Username: username})
		return
	}
	if !ok {
		ctx.JSON(http.StatusNotFound, model.ErrorResponse{Error: "user not found", Username: username})
		return
	}

	ctx.JSON(http.StatusOK, model.ProfileSeedResponse{
		UUID:         u.UUID,
		Name:         u.Name,
		Surname:      u.Surname,
		Email:        u.Email,
		PersonalCode: u.PersonalCode,
	})
}