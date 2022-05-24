package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/model/req"
	"github.com/se2022-qiaqia/course-system/model/resp"
	S "github.com/se2022-qiaqia/course-system/services"
	"github.com/se2022-qiaqia/course-system/token"
	"gorm.io/gorm"
)

type Public struct{}

// Login
// @Summary					登录。
// @Description				登录。
// @Tags					公共
// @Accept					json
// @Produce					json
// @Param					params			body		req.LoginCredit		true	"登录凭据"
// @Success 				200 			{object}	resp.OkResponse{data=string}
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/login		 	[post]
func (api *Public) Login(c *gin.Context) {
	var credit req.LoginCredit
	if !req.BindAndValidate(c, &credit) {
		return
	}

	user, err := S.Services.Public.Login(credit)
	if err == nil {
		resp.Ok(token.NewToken(user), c)
		return
	} else if errors.Is(err, S.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(resp.ErrCodeNotFound, "找不到对应用户", c)
		return
	} else if errors.Is(err, S.ErrWrongPassword) {
		resp.Fail(resp.ErrCodeUnauthorized, "密码错误", c)
		return
	}
	resp.FailJust("登录失败", c)
	return
}

// Register
// @Summary					注册。
// @Description				注册。
// @Tags					公共
// @Accept					json
// @Produce					json
// @Param					params			body		req.RegisterInfo		true	"注册信息"
// @Success 				200 			{object}	resp.OkResponse{data=boolean}
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/register		[post]
func (api *Public) Register(c *gin.Context) {
	var b req.RegisterInfo
	if !req.BindAndValidate(c, &b) {
		return
	}

	if err := dao.DB.Model(&dao.College{}).Where("id = ?", b.CollegeId).First(&dao.College{}).Error; err != nil {
		resp.Fail(resp.ErrCodeNotFound, "找不到对应学院", c)
		return
	}

	ok, err := S.Services.Public.Register(b)
	if err == nil {
		resp.Ok(ok, c)
		return
	} else if errors.Is(err, S.ErrConflict) {
		resp.Fail(resp.ErrCodeNotFound, "用户已存在", c)
		return
	} else {
		resp.FailJust("注册失败", c)
		return
	}
}
