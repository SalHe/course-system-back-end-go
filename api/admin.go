package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api/resp"
	"github.com/se2022-qiaqia/course-system/dao"
)

func IsInitialized(c *gin.Context) {
	var count int64
	if err := dao.DB.Model(&dao.User{}).Count(&count).Error; err == nil && count >= 1 {
		c.JSON(http.StatusOK, resp.Response{Msg: "系统已初始化"})
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, resp.Response{Msg: "系统未初始化"})
	}
}

func InitSystem(c *gin.Context) {
	// TODO 实现系统初始化
	c.JSON(http.StatusOK, resp.Response{Msg: "初始化成功"})
}
