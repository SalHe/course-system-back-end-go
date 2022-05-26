package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/middleware"
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
// @Success 				200 			{object}	resp.OkResponse{data=resp.User}
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/user		[get]
func (api *User) GetUserInfo(c *gin.Context) {
	cla, _ := c.Get(middleware.ClaimsKey)
	claims := cla.(*token.JwtClaims)
	// resp.Ok(resp.NewUser(claims.User), c)

	user, err := S.Services.User.GetUser(claims.User.ID)
	if err == nil {
		resp.Ok(resp.NewUser(user), c)
		return
	} else {
		resp.FailJust("查询失败", c)
		return
	}
}

// GetOtherUserInfo
// @Summary					获取任意用户信息。
// @Description
// @Tags					用户
// @Accept					json
// @Produce					json
// @Security				ApiKeyAuth
// @Param					id				path		int			true		"用户id"
// @Success 				200 			{object}	resp.OkResponse{data=resp.User}
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/user/{id}		[get]
func (api *User) GetOtherUserInfo(c *gin.Context) {
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
// @Success 				200 			{object}	resp.OkResponse{data=boolean}
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/user/new		[post]
func (api *User) NewUser(c *gin.Context) {
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
// @Param					page			body			req.Page			false	"分页"
// @Success 				200 			{object}		resp.OkResponse{data=resp.Page{contents=[]resp.User}}
// @Failure 				400 			{object} 		resp.ErrorResponse
// @Router					/user/list		[post]
func (api *User) GetUserList(c *gin.Context) {
	var b req.Page
	if !req.BindAndValidate(c, &b) {
		return
	}

	count := S.Services.User.GetUserCount()

	users, err := S.Services.User.GetUserList(b)
	if err == nil {
		usersResp := make([]*resp.User, len(users))
		for i, user := range users {
			usersResp[i] = resp.NewUser(&user)
		}
		resp.OkPage(usersResp, b.ActualPage(), b.ActualSize(), count, c)
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
// @Success 				200 			{object}	resp.OkResponse{data=boolean}
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/user/{id}		[delete]
func (api *User) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	cla, _ := c.Get(middleware.ClaimsKey)
	claims := cla.(*token.JwtClaims)

	if claims.User.ID == uint(id) {
		resp.FailJust("不能删除自己哦~", c)
		return
	}

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

// UpdateSelfInfo
// @Summary					更新用户信息。
// @Description
// @Tags					用户
// @Accept					json
// @Produce					json
// @Security				ApiKeyAuth
// @Param					info 			body		req.UpdateUserRequest			true		"新用户信息"
// @Success 				200 			{object}	resp.OkResponse{data=resp.User} "更新后的用户信息"
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/user [post]
func (api *User) UpdateSelfInfo(c *gin.Context) {
	cla, _ := c.Get(middleware.ClaimsKey)
	claims := cla.(*token.JwtClaims)
	id := claims.ID
	var b req.UpdateUserRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	updated, err := S.Services.User.UpdateUser(id, b, claims.IsAdmin())
	if err == nil {
		resp.Ok(resp.NewUser(updated), c)
		return
	} else {
		resp.Fail(resp.ErrCodeInternal, "更新个人信息失败", c)
		return
	}
}

// UpdateUserInfo
// @Summary					更新任意用户信息。
// @Description
// @Tags					用户
// @Accept					json
// @Produce					json
// @Security				ApiKeyAuth
// @Param					id				path		int							true		"用户id"
// @Param					info 			body		req.UpdateUserRequest			true		"新用户信息"
// @Success 				200 			{object}	resp.OkResponse{data=resp.User} "更新后的用户信息"
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/user/{id}		[post]
func (api *User) UpdateUserInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var b req.UpdateUserRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	updated, err := S.Services.User.UpdateUser(uint(id), b, true)
	if err == nil {
		resp.Resp(resp.NewUser(updated),
			"更新成功，非管理员仅允许更新用户名，同时不允许对用户提权为管理员，也不允许将管理员降权为其他用户。",
			c)
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, S.ErrNotFound) {
		resp.Fail(resp.ErrCodeNotFound, "未找到对应用户", c)
		return
	} else {
		resp.Fail(resp.ErrCodeInternal, "更新用户失败", c)
		return
	}
}

// UpdatePassword
// @Summary					更新用户密码。
// @Description
// @Tags					用户
// @Accept					json
// @Produce					json
// @Security				ApiKeyAuth
// @Param					id				path		int								false		"用户id"
// @Param					info 			body		req.UpdateUserPassword			true		"新用户信息"
// @Success 				200 			{object}	resp.OkResponse{data=boolean}
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/user/pwd [post]
// @Router					/user/{id}/pwd [post]
func (api *User) UpdatePassword(c *gin.Context) {
	cla, _ := c.Get(middleware.ClaimsKey)
	claims := cla.(*token.JwtClaims)

	idString, hasId := c.Params.Get("id")
	i, _ := strconv.Atoi(idString)

	// 确定一下被修改的用户ID
	id := uint(i)                             // 路径参数里指定的用户ID
	curToken, _ := c.Get(middleware.TokenKey) // 当前登录的用户token
	if hasId {
		if !claims.IsAdmin() {
			// 正常来说不会给普通用户授权，但是这里还是判断下，防止写错
			resp.Fail(resp.ErrCodeUnauthorized, "还是不要乱改别人密码哟~", c)
			return
		}
		curToken = "" // TODO 找到被修改用户的token，在后面将其从tokenStorage中删除（登出）
	} else {
		id = claims.User.ID
	}

	var b req.UpdateUserPassword
	if !req.BindAndValidate(c, &b) {
		return
	}

	err := S.Services.User.UpdatePassword(id, b)
	if err == nil {
		token.Storage.Delete(curToken.(string))
		resp.Ok(true, c)
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, S.ErrNotFound) {
		resp.FailJust("未找到对应用户", c)
		return
	} else {
		resp.FailJust("修改失败", c)
		return
	}
}
