package api

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/middleware"
	R "github.com/se2022-qiaqia/course-system/router"
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
	R.Router.Public.Init(publicRouter)
	R.Router.Start.Init(publicRouter)

	// 需要认证的API
	authenticatedRouter := r.Group("")
	authenticatedRouter.Use(middleware.AuthorizedRequired)
	R.Router.User.Init(authenticatedRouter)
	R.Router.Course.Init(authenticatedRouter)
	R.Router.College.Init(authenticatedRouter)

	return engine
}
