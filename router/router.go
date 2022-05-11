package router

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/middleware"
)

func NewRouter() *gin.Engine {
	engine := gin.Default()
	gin.SetMode(gin.DebugMode)
	if !config.Config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine.Use(middleware.Authorize)

	r := engine.Group("/api/v1")

	// 不需要认证的API
	publicRouter := r.Group("")
	Router.Public.Init(publicRouter)
	Router.Start.Init(publicRouter)

	// 需要认证的API
	authenticatedRouter := r.Group("")
	authenticatedRouter.Use(middleware.AuthorizedRequired)
	Router.User.Init(authenticatedRouter)
	Router.Course.Init(authenticatedRouter)
	Router.College.Init(authenticatedRouter)

	return engine
}
