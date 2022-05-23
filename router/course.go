package router

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/middleware"
)

type Course struct{}

func (c Course) Init(Router *gin.RouterGroup) {
	r := Router.Group("/course")

	sr := r.Group("")
	sr.Use(middleware.AuthorizedRoleRequired(dao.RoleStudent))

	ar := r.Group("")
	ar.Use(middleware.AuthorizedRoleRequired(dao.RoleAdmin))

	asr := r.Group("")
	asr.Use(middleware.AuthorizedRoleRequired(dao.RoleAdmin, dao.RoleStudent))

	{
		r.POST("/list", api.Api.Course.GetCourseList)
		r.POST("/schedules", api.Api.Course.GetCourseSchedules)
	}
	{
		// 选撤课
		asr.POST("/select", api.Api.Course.SelectCourse)
		asr.DELETE("/select", api.Api.Course.UnSelectCourse)
	}
	{
		// 管理课程相关
		ar.POST("", api.Api.Course.NewCourse)
		ar.PUT("/:id", api.Api.Course.UpdateCourseCommon)
		ar.PUT("/spec/:id", api.Api.Course.UpdateCourseSpecific)
		ar.POST("/open", api.Api.Course.OpenCourse)
	}
	{
		atr := r.Group("")
		atr.Use(middleware.AuthorizedRoleRequired(dao.RoleAdmin, dao.RoleTeacher))
		atr.GET("/spec/:id", api.Api.Course.GetCourseSpecificDetails)
	}
}
