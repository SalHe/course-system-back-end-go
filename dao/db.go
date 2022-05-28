package dao

import (
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

type Role = uint
type CourseStatus = uint

const (
	RoleStudent = Role(0)
	RoleTeacher = Role(100)
	RoleAdmin   = Role(200)
)

const (
	CourseStatusNormal   = CourseStatus(0)   // 正常（已选上，正常上课中）
	CourseStatusWithdraw = CourseStatus(100) // 撤销（已退课）
	// CourseStatusClosed                        // 已结课
	// 结课状态应该表示在课程中，而不是`学生-课程`中
)

func Init() {
	var dialector gorm.Dialector
	database := config.Config.Database
	switch {
	case database.Postgres != nil:
		dialector = postgres.Open(database.Postgres.DSN())
	case database.Mysql != nil:
		dialector = mysql.Open(database.Mysql.DSN())
	case database.Sqlite != nil:
		dialector = sqlite.Open(database.Sqlite.Filename)
	default:
		panic("未正确配置数据库!")
	}

	if db, err := gorm.Open(dialector, &gorm.Config{Logger: &gormZeroLogger{Logger: log.L}}); err != nil {
		panic("连接数据库失败")
	} else {
		DB = db
	}
	if config.Config.Debug {
		DB = DB.Debug()
	}
}

func Migrate() {
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&College{})
	DB.AutoMigrate(&CourseCommon{})
	DB.AutoMigrate(&CourseSpecific{})
	DB.AutoMigrate(&Semester{})
	DB.AutoMigrate(&CourseSchedule{})
	DB.AutoMigrate(&StudentCourse{})
	DB.AutoMigrate(&Setting{})
}

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// College 学院
type College struct {
	Model
	Name string `gorm:"unique;not null;"` // 学院名
}

// Semester 学期
type Semester struct {
	Model
	Year uint // 年份
	Term uint // 对应年份第几学期
}

func (u *User) SetPassword(password string) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		panic("加密出错")
	}
	u.Password = string(bytes)
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
