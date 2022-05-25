package router

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/middleware"
)

type College struct{}

func (c College) Init(Router *gin.RouterGroup, RouterNoAuth *gin.RouterGroup) {
	publicRouter := RouterNoAuth.Group("/college")
	authRouter := Router.Group("/college")

	privateRouter := authRouter.Group("")
	privateRouter.Use(middleware.AuthorizedRoleRequired(dao.RoleAdmin))

	{
		publicRouter.POST("/list", api.Api.College.ListColleges)
	}
	{
		privateRouter.POST("/new", api.Api.College.NewCollege)
	}
}
