package router

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/middleware"
	"github.com/se2022-qiaqia/course-system/model/req"
)

func NewRouter(engine *gin.Engine) *gin.Engine {
	if config.Config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	req.InitValidation()

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
	Router.College.Init(authenticatedRouter, publicRouter)
	Router.Semester.Init(authenticatedRouter)

	return engine
}
