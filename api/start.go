package api

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/model/resp"
	"gorm.io/gorm"
	"net/http"
)

type Start struct{}

func isInitialized() bool {
	var count int64
	err := dao.DB.Model(&dao.User{}).Count(&count).Error
	return err == nil && count >= 1
}

func (api Start) IsInitialized(c *gin.Context) {
	if isInitialized() {
		c.JSON(http.StatusOK, resp.Response{Msg: "系统已初始化"})
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, resp.Response{Msg: "系统未初始化"})
	}
}

type InitRequest struct {
	Id               uint   `json:"id" binding:"required"`
	Username         string `json:"username" binding:"required"`
	Password         string `json:"password" binding:"required"`
	RealName         string `json:"realName" binding:"required"`
	AdminCollegeName string `json:"adminCollegeName" binding:"required"`
}

func (api Start) InitSystem(c *gin.Context) {
	if isInitialized() {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.Response{Msg: "系统已初始化，不可再次初始化"})
		return
	} else {
		var b InitRequest
		if err := c.ShouldBindJSON(&b); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, resp.Response{Msg: "请输入参数"})
			return
		}

		if b.Id == 0 || b.Username == "" || b.Password == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, resp.Response{Msg: "请完整指定有效的ID、用户名、密码"})
			return
		}

		var adminCollegeName = b.AdminCollegeName
		if len(adminCollegeName) == 0 {
			adminCollegeName = "管理员"
		}

		user := &dao.User{
			Model:    gorm.Model{ID: b.Id},
			Username: b.Username,
			RealName: b.RealName,
			Role:     dao.RoleAdmin,
			College: dao.College{
				Name: adminCollegeName,
			},
		}
		user.SetPassword(b.Password)
		result := dao.DB.Create(user)
		if err := result.Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp.Response{Msg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp.Response{Msg: "初始化成功"})
		return
	}
}
