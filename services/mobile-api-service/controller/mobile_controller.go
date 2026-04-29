package controller

import (
	"log"
	"net/http"

	"mobile-api-service/model"
	"mobile-api-service/service"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	reqCtx := ctx.Request.Context()

	tr := otel.Tracer("mobile-api-service")
	reqCtx, span := tr.Start(reqCtx, "MobileController.GetProfile")
	defer span.End()

	username := ctx.Param("username")
	span.SetAttributes(attribute.String("app.username", username))

	var query model.ProfileLookupQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Printf("invalid profile lookup query parameters (username=%s): %v", username, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid query parameters")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query.SetDefaults()

	status, contentType, body, err := c.orch.FetchProfileByUsername(
		reqCtx,
		username,
		query,
	)
	if err != nil {
		log.Printf("GetProfile failed (username=%s): %v", username, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.Data(status, contentType, body)
}
