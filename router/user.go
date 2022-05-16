package router

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/middleware"
)

type User struct {
}

func (s User) Init(Router *gin.RouterGroup) {
	publicRouter := Router.Group("/user")

	privateRouter := publicRouter.Group("")
	privateRouter.Use(middleware.AuthorizedRoleRequired(dao.RoleAdmin))

	{
		publicRouter.GET("", api.Api.User.GetUserInfo)
		publicRouter.POST("", api.Api.User.UpdateSelfInfo)
		publicRouter.POST("/pwd", api.Api.User.UpdatePassword)
	}
	{
		privateRouter.GET("/list/:page/:size", api.Api.User.GetUserList)
		privateRouter.GET("/:id", api.Api.User.GetOtherUserInfo)
		privateRouter.POST("/:id", api.Api.User.UpdateUserInfo)
		privateRouter.DELETE("/:id", api.Api.User.DeleteUser)
		privateRouter.POST("/new", api.Api.User.NewUser)
		publicRouter.POST("/:id/pwd", api.Api.User.UpdatePassword)
	}
}
