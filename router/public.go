package router

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api"
)

type Public struct{}

func (p Public) Init(r *gin.RouterGroup) {
	r.POST("/login", api.Api.Public.Login)
	r.POST("/register", api.Api.Public.Register)
}
