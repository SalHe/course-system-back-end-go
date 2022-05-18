package router

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/middleware"
)

type Course struct{}

func (c Course) Init(Router *gin.RouterGroup) {
	publicRouter := Router.Group("/course")

	privateRouter := publicRouter.Group("")
	privateRouter.Use(middleware.AuthorizedRoleRequired(dao.RoleAdmin))

	{
		publicRouter.POST("/list", api.Api.Course.GetCourseList)
	}
	{
		privateRouter.POST("", api.Api.Course.NewCourse)
		privateRouter.PUT("/:id", api.Api.Course.UpdateCourseCommon)
		privateRouter.PUT("/spec/:id", api.Api.Course.UpdateCourseSpecific)
		privateRouter.POST("/open", api.Api.Course.OpenCourse)
	}
}
