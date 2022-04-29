package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api/resp"
	"github.com/se2022-qiaqia/course-system/dao"
	"gorm.io/gorm"
)

func isInitialized() bool {
	var count int64
	err := dao.DB.Model(&dao.User{}).Count(&count).Error
	return err == nil && count >= 1
}

func IsInitialized(c *gin.Context) {
	if isInitialized() {
		c.JSON(http.StatusOK, resp.Response{Msg: "系统已初始化"})
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, resp.Response{Msg: "系统未初始化"})
	}
}

type InitRequest struct {
	Id               uint   `json:"id"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	RealName         string `json:"realName"`
	AdminCollegeName string `json:"adminCollegeName"`
}

func InitSystem(c *gin.Context) {
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
		if err := dao.DB.Create(&dao.College{
			Model: gorm.Model{ID: 0},
			Name:  adminCollegeName,
		}).Error; err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		user := &dao.User{
			Username: b.Username,
			RealName: b.RealName,
			Role:     dao.RoleAdmin,
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
