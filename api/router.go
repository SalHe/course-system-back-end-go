package api

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
	{
		r := r.Group("/")
		r.POST("/login", Login)
		r.POST("/register", Register)
	}
	r = r.Group("/")
	r.Use(middleware.AuthorizedRequired)
	{
		user := r.Group("/user")
		user.GET("/info", GetUserInfo)
	}
	return engine
}
