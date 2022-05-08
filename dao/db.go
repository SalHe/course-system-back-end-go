package dao

import (
	"github.com/se2022-qiaqia/course-system/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type Role = uint

const (
	RoleStudent = Role(iota)
	RoleTeacher
	RoleAdmin
)

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
		panic("连接数据库失败")
	} else {
		DB = db
	}
	if config.Config.Debug {
		DB.Logger.LogMode(logger.Info)
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
}

// User 定义用户的信息
type User struct {
	gorm.Model
	Username     string  `gorm:"unique;not null;"` // 用户名，可以自定义
	RealName     string  `gorm:"not null"`         // 真实姓名
	Password     string  `gorm:"not null"`         // 密码
	Role         Role    `gorm:"not null"`         // 角色，本系统中，一个用户只能有用一种角色
	CollegeId    uint    `gorm:"not null"`         // 所属学院ID
	College      College // 所属学院
	EntranceYear uint    // 入驻学校年份
}

// College 学院
type College struct {
	gorm.Model
	Name string `gorm:"unique;not null;"` // 学院名
}

// CourseCommon 课程公共信息
type CourseCommon struct {
	gorm.Model
	Name string `gorm:"not null"` // 课程名

	Credits float32 // 学分
	Hours   uint    // 学时，单位不一定是小时

	CollegeId uint    `gorm:"not null"` // 开课学院ID
	College   College // 开课学院

	CourseSpecifics []CourseSpecific // 课程具体信息
}

// CourseSpecific 具体课头，指具体开给某一个老师的课程
type CourseSpecific struct {
	gorm.Model
	CourseCommonId  uint              `gorm:"not null"` // 课程公共信息ID
	CourseCommon    CourseCommon      // 课程公共信息
	TeacherId       uint              `gorm:"not null"`             // 授课教师ID
	Teacher         User              `gorm:"foreignkey:TeacherId"` // 授课教师
	Location        string            // 上课地点
	Quota           uint              // 人数配额
	QuotaUsed       uint              // 课程已选人数
	SemesterId      uint              // 学期ID
	Semester        Semester          // 学期
	CourseSchedules []*CourseSchedule `gorm:"many2many:course_specific_course_schedule;"`
}

// Semester 学期
type Semester struct {
	gorm.Model
	Year uint // 年份
	Term uint // 对应年份第几学期
}

// CourseSchedule 上课时间
type CourseSchedule struct {
	gorm.Model
	DayOfWeek uint // 每周第几天
	HoursId   uint // 第几节课
}

// StudentCourse 学生和具体某门课的关系
type StudentCourse struct {
	gorm.Model
	StudentId uint           // 学生ID
	Student   User           `gorm:"foreignKey:StudentId"` // 学生
	CourseId  uint           // 课程ID
	Course    CourseSpecific `gorm:"foreignKey:CourseId"` // 具体课程
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
