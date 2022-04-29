package dao

import (
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var dialector gorm.Dialector
	database := config.Config.Database
	switch {
	case database.Sqlite != nil:
		dialector = sqlite.Open(database.Sqlite.Filename)
	default:
		panic("未正确配置数据库!")
	}
	if db, err := gorm.Open(dialector, &gorm.Config{}); err != nil {
		log.Logger.Fatalln("连接数据库失败")
	} else {
		DB = db
	}
}

func Migrate() {
	DB.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null;"`
	Password string `gorm:"not null"`
	IsAdmin  bool   `gorm:"not null"`
}

func (u *User) SetPassword(password string) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		panic("加密出错")
	}
	u.Password = string(bytes)
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil
}
