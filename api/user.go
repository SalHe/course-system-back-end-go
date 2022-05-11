package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api/token"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/model/resp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strconv"
)

type User struct{}

func (api User) GetUserInfo(c *gin.Context) {
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

func (api User) GetOtherUserInfo(c *gin.Context) {
	id := c.Param("id")
	var user dao.User

	if err := dao.DB.Preload(clause.Associations).Find(&user, id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, resp.Response{Msg: "未找到对应用户"})
		return
	}
	c.JSON(http.StatusOK, resp.Response{Msg: "获取用户成功", Data: resp.User{
		Id:       user.ID,
		Username: user.Username,
		RealName: user.RealName,
		Role:     user.Role,
		College: resp.College{
			Id:   user.CollegeId,
			Name: user.College.Name,
		},
		EntranceYear: user.EntranceYear,
	}})
	return
}

type NewUserRequest struct {
	Id           uint     `json:"id" binding:"required"`
	Username     string   `json:"username" binding:"required" validate:"min=5;max=20"`
	Password     string   `json:"password" binding:"required" validate:"regexp=^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)[a-zA-Z\\d]{8,}$"`
	RealName     string   `json:"realName" binding:"required" validate:"min=1;max=10"`
	CollegeId    uint     `json:"collegeId" binding:"required"`
	Role         dao.Role `json:"role"`
	EntranceYear uint     `json:"entranceYear" binding:"required"`
}

func (api User) NewUser(c *gin.Context) {
	var b NewUserRequest
	err := c.BindJSON(&b)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Response{Msg: "请输入参数"})
		return
	}
	var user *dao.User
	if err := dao.DB.Model(&dao.User{}).Where("id = ? OR username = ? OR username = ?", b.Id, b.Username, b.Username, b.Id).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		user := &dao.User{
			Model:        gorm.Model{ID: b.Id},
			Username:     b.Username,
			RealName:     b.RealName,
			Role:         b.Role,
			CollegeId:    b.CollegeId,
			EntranceYear: b.EntranceYear,
		}
		user.SetPassword(b.Password)
		result := dao.DB.Create(user)
		if err := result.Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp.Response{Msg: "添加用户失败!"})
			return
		}
		c.JSON(http.StatusOK, resp.Response{Msg: "创建用户成功"})
		return
	} else {
		c.AbortWithStatusJSON(http.StatusConflict, resp.Response{Msg: "用户已存在"})
		return
	}
}

func (api User) GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.Param("page"))
	size, _ := strconv.Atoi(c.Param("size"))
	var users []dao.User
	if err := dao.DB.Preload(clause.Associations).Offset((page - 1) * size).Limit(size).Find(&users).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, resp.Response{Msg: "未找到用户列表"})
		return
	}
	usersResp := make([]*resp.User, len(users))
	for i, user := range users {
		usersResp[i] = &resp.User{
			Id:       user.ID,
			Username: user.Username,
			RealName: user.RealName,
			Role:     user.Role,
			College: resp.College{
				Id:   user.CollegeId,
				Name: user.College.Name,
			},
			EntranceYear: user.EntranceYear,
		}
	}
	c.JSON(http.StatusOK, resp.Response{Msg: "获取用户列表成功", Data: usersResp})
	return
}

func (api User) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user dao.User
	var count int64
	if dao.DB.Model(&dao.User{}).Find(&user, "id = ?", id).Count(&count); count <= 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, resp.Response{Msg: "未找到对应用户或已被删除"})
		return
	}
	if err := dao.DB.Where("id = ?", id).Delete(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.Response{Msg: "删除用户失败"})
		return
	}
	c.JSON(http.StatusOK, resp.Response{Msg: "删除用户成功"})
	return
}

type UpdateUserRequest struct {
	Username     string   `json:"username" binding:"required" validate:"min=5;max=20"`
	RealName     string   `json:"realName" binding:"required" validate:"min=1;max=10"`
	CollegeId    uint     `json:"collegeId" binding:"required"`
	Role         dao.Role `json:"role"`
	EntranceYear uint     `json:"entranceYear" binding:"required"`
}

func (api User) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user dao.User
	if err := dao.DB.Find(&user, id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, resp.Response{Msg: "未找到对应用户"})
		return
	}
	var b UpdateUserRequest
	err := c.BindJSON(&b)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Response{Msg: "请输入参数"})
		return
	}
	user.Username = b.Username
	user.RealName = b.RealName
	user.Role = b.Role
	user.CollegeId = b.CollegeId
	user.EntranceYear = b.EntranceYear
	if err := dao.DB.Save(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.Response{Msg: "更新用户失败"})
		return
	}
	c.JSON(http.StatusOK, resp.Response{Msg: "更新用户成功"})
	return
}
