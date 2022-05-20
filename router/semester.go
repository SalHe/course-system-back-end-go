package router

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/middleware"
)

type Semester struct{}

func (s *Semester) Init(Router *gin.RouterGroup) {
	r := Router.Group("/semester")
	ar := r.Group("")
	ar.Use(middleware.AuthorizedRoleRequired(dao.RoleAdmin))
	{
		r.GET("/", api.Api.Semester.QuerySemester)
		r.GET("/curr", api.Api.Semester.GetCurrentSemester)
	}
	{
		ar.POST("/", api.Api.Semester.CreateSemester)
		ar.POST("/curr", api.Api.Semester.SetCurrentSemester)
	}
}
