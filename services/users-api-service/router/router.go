package router

import (
	"users-api-service/controller"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
}

func New() *Router {
	r := gin.New()
	r.Use(gin.Recovery())
	return &Router{engine: r}
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}

func (r *Router) Setup(ctrl *controller.UsersController, basePath string) {
	r.engine.GET("/health", ctrl.Health)
	r.engine.GET(basePath+"/:username", ctrl.GetUserUUID)
}
