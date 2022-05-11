package router

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api/admin"
)

type Start struct{}

func (i Start) Init(r *gin.RouterGroup) {
	r.GET("/init", admin.IsInitialized)
	r.POST("/init", admin.InitSystem)
}
