package dao

import "gorm.io/gorm"

// User 定义用户的信息
type User struct {
	Model
	Username     string  `gorm:"unique;not null;"` // 用户名，可以自定义
	RealName     string  `gorm:"not null"`         // 真实姓名
	Password     string  `gorm:"not null"`         // 密码
	Role         Role    `gorm:"not null"`         // 角色，本系统中，一个用户只能有用一种角色
	CollegeId    uint    `gorm:"not null"`         // 所属学院ID
	College      College // 所属学院
	EntranceYear uint    // 入驻学校年份
}

func FindUserById(db *gorm.DB, id uint) *gorm.DB {
	return db.Model(&User{}).Where("id = ?", id)
}
