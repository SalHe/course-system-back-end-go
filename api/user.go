package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/model/req"
	"github.com/se2022-qiaqia/course-system/model/resp"
	S "github.com/se2022-qiaqia/course-system/services"
	"github.com/se2022-qiaqia/course-system/token"
	"gorm.io/gorm"
	"strconv"
)

type User struct{}

// GetUserInfo
// @Summary					获取当前登录用户信息。
// @Description
// @Tags					用户
// @Accept					json
// @Produce					json
// @Security				ApiKeyAuth
// @Success 				200 			{object}	resp.User
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/user/info		[get]
func (api User) GetUserInfo(c *gin.Context) {
	cla, _ := c.Get("claims")
	claims := cla.(*token.JwtClaims)
	resp.Ok(resp.NewUser(claims.User), c)
}

// GetOtherUserInfo
// @Summary					获取任意用户信息。
// @Description
// @Tags					用户
// @Accept					json
// @Produce					json
// @Security				ApiKeyAuth
// @Param					id				path		int			true		"用户id"
// @Success 				200 			{object}	resp.User
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/user/{id}		[get]
func (api User) GetOtherUserInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := S.Services.User.GetUser(uint(id))
	if err == nil {
		resp.Ok(resp.NewUser(user), c)
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(resp.ErrCodeNotFound, "未找到对应用户", c)
		return
	} else {
		resp.FailJust("查询失败", c)
		return
	}
}

// NewUser
// @Summary					添加用户。
// @Description
// @Tags					用户
// @Accept					json
// @Produce					json
// @Security				ApiKeyAuth
// @Param					params			body 		req.NewUserRequest	true		"添加用户"
// @Success 				200 			{object}	boolean
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/user/new		[post]
func (api User) NewUser(c *gin.Context) {
	var b req.NewUserRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	err := S.Services.User.NewUser(b)
	if err == nil {
		resp.Ok(true, c)
		return
	} else if errors.Is(err, S.ErrConflict) {
		resp.Fail(resp.ErrCodeConflict, "用户已存在", c)
		return
	} else {
		resp.FailJust("添加失败", c)
		return
	}
}

// GetUserList
// @Summary					获取用户列表。
// @Description
// @Tags					用户
// @Accept					json
// @Produce					json
// @Security				ApiKeyAuth
// @Param					page 			path 		int			false		"页码"
// @Param					size 			path 		int			false		"每页数量"
// @Success 				200 			{array}		resp.User
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/user/list/{page}/{size}		[get]
func (api User) GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.Param("page"))
	size, _ := strconv.Atoi(c.Param("size"))

	users, err := S.Services.User.GetUserList(req.Page{Page: page, Size: size})
	if err == nil {
		usersResp := make([]*resp.User, len(users))
		for i, user := range users {
			usersResp[i] = resp.NewUser(&user)
		}
		resp.Ok(usersResp, c)
		return
	} else {
		resp.Fail(resp.ErrCodeNotFound, "未找到用户列表", c)
		return
	}
}

// DeleteUser
// @Summary					删除指定用户。
// @Description
// @Tags					用户
// @Accept					json
// @Produce					json
// @Security				ApiKeyAuth
// @Param					id				path		int			true		"用户id"
// @Success 				200 			{object}	resp.User
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/user/{id}		[delete]
func (api User) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := S.Services.User.DeleteUser(uint(id))
	if err == nil {
		resp.Ok(true, c)
		return
	} else if errors.Is(err, S.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(resp.ErrCodeNotFound, "未找到对应用户或已被删除", c)
		return
	} else {
		resp.Fail(resp.ErrCodeInternal, "删除用户失败", c)
		return
	}
}

// UpdateUser
// @Summary					更新任意用户信息。
// @Description
// @Tags					用户
// @Accept					json
// @Produce					json
// @Security				ApiKeyAuth
// @Param					id				path		int							true		"用户id"
// @Param					info 			body		req.UpdateUserRequest			true		"新用户信息"
// @Success 				200 			{object}	resp.User
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/user/{id}		[post]
func (api User) UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var b req.UpdateUserRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	oldUser, err := S.Services.User.UpdateUser(uint(id), b)
	if err == nil {
		resp.Ok(oldUser, c)
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, S.ErrNotFound) {
		resp.Fail(resp.ErrCodeNotFound, "未找到对应用户", c)
		return
	} else {
		resp.Fail(resp.ErrCodeInternal, "更新用户失败", c)
		return
	}
}
