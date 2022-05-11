package router

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/middleware"
)

type College struct{}

func (c College) Init(Router *gin.RouterGroup) {
	publicRouter := Router.Group("/college")

	privateRouter := publicRouter.Group("")
	privateRouter.Use(middleware.AuthorizedRoleRequired(dao.RoleAdmin))

	{
		publicRouter.GET("/list", api.Api.College.ListColleges)
	}
	{
		privateRouter.POST("/new", api.Api.College.NewCollege)
	}
}
