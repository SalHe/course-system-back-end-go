package router

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api/admin"
	"github.com/se2022-qiaqia/course-system/api/user"
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
		publicRouter.GET("/info", user.GetUserInfo)
	}
	{
		privateRouter.GET("/list/:page/:size", admin.GetUserList)
		privateRouter.GET("/:id", admin.GetUser)
		privateRouter.POST("/:id", admin.UpdateUser)
		privateRouter.DELETE("/:id", admin.DeleteUser)
		privateRouter.POST("/new", admin.NewUser)
	}
}
