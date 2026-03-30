package router

import (
	"profile-service/controller"

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

func (r *Router) Setup(ctrl *controller.ProfileController, basePath string) {
	r.engine.GET("/health", ctrl.Health)
	r.engine.GET(basePath+"/profiles/:uuid", ctrl.GetProfile)
}