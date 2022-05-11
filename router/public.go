package router

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api/user"
)

type Public struct{}

func (p Public) Init(r *gin.RouterGroup) {
	r.POST("/login", user.Login)
	r.POST("/register", user.Register)
}
