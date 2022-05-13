package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api/token"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/model/req"
	"github.com/se2022-qiaqia/course-system/model/resp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type User struct{}

func (api User) GetUserInfo(c *gin.Context) {
	cla, _ := c.Get("claims")
	claims := cla.(*token.JwtClaims)
	resp.Ok(map[string]interface{}{
		"id":       claims.Id,
		"username": claims.Username,
		"role":     claims.Role,
	}, c)
}

func (api User) GetOtherUserInfo(c *gin.Context) {
	id := c.Param("id")
	var user dao.User

	if err := dao.DB.Preload(clause.Associations).Model(&dao.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		resp.Fail(resp.ErrCodeNotFound, "未找到对应用户", c)
		return
	}
	resp.Ok(resp.NewUser(&user), c)
	return
}

type NewUserRequest struct {
	Id           uint     `json:"id" binding:"required"`
	Username     string   `json:"username" binding:"required,username" description:"用户名"`
	Password     string   `json:"password" binding:"required,password" description:"密码"`
	RealName     string   `json:"realName" binding:"required,min=1,max=10" description:"真实姓名"`
	CollegeId    uint     `json:"collegeId" binding:"required" description:"学院id"`
	Role         dao.Role `json:"role" description:"角色"`
	EntranceYear uint     `json:"entranceYear" binding:"required,min=1980" description:"入学/入职年份"`
}

func (api User) NewUser(c *gin.Context) {
	var b NewUserRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	var user *dao.User
	if err := dao.DB.Model(&dao.User{}).Where("id = ? OR username = ? OR username = ?", b.Id, b.Username, b.Username, b.Id).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		user := &dao.User{
			Model:        dao.Model{ID: b.Id},
			Username:     b.Username,
			RealName:     b.RealName,
			Role:         b.Role,
			CollegeId:    b.CollegeId,
			EntranceYear: b.EntranceYear,
		}
		user.SetPassword(b.Password)
		result := dao.DB.Create(user)
		if err := result.Error; err != nil {
			resp.Fail(resp.ErrCodeInternal, "添加用户失败!", c)
			return
		}
		resp.Ok(true, c)
		return
	} else {
		resp.Fail(resp.ErrCodeConflict, "用户已存在!", c)
		return
	}
}

func (api User) GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.Param("page"))
	size, _ := strconv.Atoi(c.Param("size"))
	var users []dao.User
	if err := dao.DB.Preload(clause.Associations).Offset((page - 1) * size).Limit(size).Find(&users).Error; err != nil {
		resp.Fail(resp.ErrCodeNotFound, "未找到用户列表", c)
		return
	}
	usersResp := make([]*resp.User, len(users))
	for i, user := range users {
		usersResp[i] = resp.NewUser(&user)
	}
	resp.Ok(usersResp, c)
	return
}

func (api User) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user dao.User
	var count int64
	if dao.DB.Model(&dao.User{}).Find(&user, "id = ?", id).Count(&count); count <= 0 {
		resp.Fail(resp.ErrCodeNotFound, "未找到对应用户或已被删除", c)
		return
	}
	if err := dao.DB.Where("id = ?", id).Delete(&user).Error; err != nil {
		resp.Fail(resp.ErrCodeInternal, "删除用户失败", c)
		return
	}
	resp.Ok(true, c)
	return
}

type UpdateUserRequest struct {
	Username     string   `json:"username" binding:"required,username" description:"用户名"`
	RealName     string   `json:"realName" binding:"required,min=1,max=10" description:"真实姓名"`
	CollegeId    uint     `json:"collegeId" binding:"required,min=1" description:"学院id"`
	Role         dao.Role `json:"role" binding:"required" description:"角色"`
	EntranceYear uint     `json:"entranceYear" binding:"required,min=1980" description:"入学/入职年份"`
}

func (api User) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user dao.User
	if err := dao.DB.Find(&user, id).Error; err != nil {
		resp.Fail(resp.ErrCodeNotFound, "未找到对应用户", c)
		return
	}
	var b UpdateUserRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	user.Username = b.Username
	user.RealName = b.RealName
	user.Role = b.Role
	user.CollegeId = b.CollegeId
	user.EntranceYear = b.EntranceYear
	if err := dao.DB.Save(&user).Error; err != nil {
		resp.Fail(resp.ErrCodeInternal, "更新用户失败", c)
		return
	}
	resp.Ok(true, c)
	return
}
