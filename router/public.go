package router

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/middleware"
)

type Public struct{}

func (p Public) Init(r *gin.RouterGroup) {
	r.POST("/login", api.Api.Public.Login)
	r.POST("/register", api.Api.Public.Register)

	r.GET("/register/enable", api.Api.Public.CanRegister)

	ar := r.Group("")
	ar.Use(middleware.AuthorizedRoleRequired(dao.RoleAdmin))
	ar.POST("/register/enable", api.Api.Public.EnableRegister)
}
