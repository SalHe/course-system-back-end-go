package router

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api"
)

type Start struct{}

func (i Start) Init(r *gin.RouterGroup) {
	r.GET("/init", api.Api.Start.IsInitialized)
	r.POST("/init", api.Api.Start.InitSystem)
}
