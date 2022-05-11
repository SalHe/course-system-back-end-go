package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api/resp"
	"github.com/se2022-qiaqia/course-system/api/token"
	"github.com/se2022-qiaqia/course-system/dao"
	"gorm.io/gorm"
)

type LoginCredit struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var loginCredit LoginCredit
	err := c.ShouldBindJSON(&loginCredit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Response{Msg: "请指定用户名和密码"})
		return
	}

	var user *dao.User
	if err := dao.DB.Model(&dao.User{}).Where("id = ? OR username = ?", loginCredit.Username, loginCredit.Username).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, resp.Response{Msg: "找不到对应用户"})
		return
	}

	if user.ComparePassword(loginCredit.Password) {
		c.JSON(http.StatusOK, resp.Response{
			Msg:  "登录成功",
			Data: token.NewToken(user),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, resp.Response{Msg: "用户不存在或密码错误"})
}

func Register(c *gin.Context) {
	var b struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Id       uint   `json:"id"`
	}
	err := c.ShouldBindJSON(&b)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Response{Msg: "请正确指定注册信息"})
		return
	}

	var user dao.User
	if err = dao.DB.Model(&dao.User{}).Where("username = ?", b.Username).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		user = dao.User{
			Model: gorm.Model{
				ID: b.Id,
			},
			Username: b.Username,
		}
		user.SetPassword(b.Password)
		if err = dao.DB.Create(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp.Response{Msg: "注册失败！"})
			return
		}
		c.JSON(http.StatusOK, resp.Response{Msg: "注册成功！"})
		return
	} else {
		c.JSON(http.StatusConflict, resp.Response{Msg: "用户已存在"})
		return
	}
}

func GetUserInfo(c *gin.Context) {
	cla, _ := c.Get("claims")
	claims := cla.(*token.JwtClaims)
	c.JSON(http.StatusOK, resp.Response{
		Data: map[string]interface{}{
			"id":       claims.Id,
			"username": claims.Username,
			"role":     claims.Role,
		},
	})
}
