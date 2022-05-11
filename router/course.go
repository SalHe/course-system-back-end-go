package router

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api/courses"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/middleware"
)

type Course struct{}

func (c Course) Init(Router *gin.RouterGroup) {
	publicRouter := Router.Group("/course")

	privateRouter := publicRouter.Group("")
	privateRouter.Use(middleware.AuthorizedRoleRequired(dao.RoleAdmin))

	{
		publicRouter.POST("/list", courses.GetCourseList)
	}
	{
		privateRouter.POST("", courses.NewCourse)
		privateRouter.PUT("/:id", courses.UpdateCourse)
		privateRouter.POST("/open", courses.OpenCourse)
	}
}
