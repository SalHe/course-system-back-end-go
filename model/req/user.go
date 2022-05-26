package req

import "github.com/se2022-qiaqia/course-system/dao"

type QueryUserRequest struct {
	Page             `json:"page"`
	Id               uint       `json:"id"`
	Username         string     `json:"username"`
	RealName         string     `json:"realName"`
	Roles            []dao.Role `json:"roles"`
	CollegesId       []uint     `json:"collegesId"`
	EntranceYearFrom uint       `json:"entranceYearFrom"`
	EntranceYearTo   uint       `json:"entranceYearTo"`
}

// NewUserRequest 新增用户信息
type NewUserRequest struct {
	Id           uint     `json:"id" binding:"required"`                                          // 用户id
	Username     string   `json:"username" binding:"required,username" description:"用户名"`         // 用户名
	Password     string   `json:"password" binding:"required,password" description:"密码"`          // 密码
	RealName     string   `json:"realName" binding:"required,min=1,max=10" description:"真实姓名"`    // 真实姓名
	CollegeId    uint     `json:"collegeId" binding:"required" description:"学院id"`                // 学院id
	Role         dao.Role `json:"role" description:"角色"`                                          // 角色
	EntranceYear uint     `json:"entranceYear" binding:"required,min=1980" description:"入学/入职年份"` // 入学/入职年份
}

// UpdateUserRequest 更新用户信息
type UpdateUserRequest struct {
	Username     string   `json:"username" binding:"required,username" description:"用户名"`         // 用户名
	RealName     string   `json:"realName" binding:"required,min=1,max=10" description:"真实姓名"`    // 真实姓名
	CollegeId    uint     `json:"collegeId" binding:"required,min=1" description:"学院id"`          // 学院id
	Role         dao.Role `json:"role" description:"角色"`                                          // 角色
	EntranceYear uint     `json:"entranceYear" binding:"required,min=1980" description:"入学/入职年份"` // 入学/入职年份
}

// UpdateUserPassword 更新用户密码
type UpdateUserPassword struct {
	Password string `json:"password" binding:"required,password" description:"新密码"`
}
