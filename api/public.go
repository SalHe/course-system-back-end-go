package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api/token"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/model/req"
	"github.com/se2022-qiaqia/course-system/model/resp"
	"gorm.io/gorm"
)

type Public struct{}

type LoginCredit struct {
	Username string `json:"username" binding:"required,username" description:"用户名"`
	Password string `json:"password" binding:"required,password" description:"密码"`
}

// Login
// @Summary					登录。
// @Description				登录。
// @Tags					公共
// @Accept					json
// @Produce					json
// @Param					params			body		LoginCredit		true	"登录凭据"
// @Success 				200 			{object}	string
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/login		 	[post]
func (api Public) Login(c *gin.Context) {
	var loginCredit LoginCredit
	if !req.BindAndValidate(c, &loginCredit) {
		return
	}

	var user *dao.User
	if err := dao.DB.Model(&dao.User{}).Where("id = ? OR username = ?", loginCredit.Username, loginCredit.Username).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(resp.ErrCodeNotFound, "找不到对应用户", c)
		return
	}

	if user.ComparePassword(loginCredit.Password) {
		resp.Ok(token.NewToken(user), c)
		return
	}
	resp.Fail(resp.ErrCodeUnauthorized, "用户不存在或密码错误", c)
	return
}

type RegisterInfo struct {
	Username string `json:"username" binding:"required,username" description:"用户名"`
	Password string `json:"password" binding:"required,password" description:"密码"`
	Id       uint   `json:"id"`
}

// Register
// @Summary					注册。
// @Description				注册。
// @Tags					公共
// @Accept					json
// @Produce					json
// @Param					params			body		RegisterInfo		true	"注册信息"
// @Success 				200 			{object}	boolean
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/register		[post]
func (api Public) Register(c *gin.Context) {
	var b RegisterInfo
	if !req.BindAndValidate(c, &b) {
		return
	}

	var user dao.User
	if err := dao.DB.Unscoped().Model(&dao.User{}).Where("id = ? OR username = ?", b.Id, b.Username).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		user = dao.User{
			Model: dao.Model{
				ID: b.Id,
			},
			Username: b.Username,
		}
		user.SetPassword(b.Password)
		if err = dao.DB.Create(&user).Error; err != nil {
			resp.FailJust("注册失败！", c)
			return
		}
		resp.Ok(true, c)
		return
	} else {
		resp.Fail(resp.ErrCodeConflict, "用户已存在", c)
		return
	}
}
