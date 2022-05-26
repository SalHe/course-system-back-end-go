package dao

import "gorm.io/gorm"

// CourseCommon 课程公共信息
type CourseCommon struct {
	Model
	Name string `gorm:"not null"` // 课程名

	Credits float32 // 学分
	Hours   uint    // 学时，单位不一定是小时

	CollegeId uint     `gorm:"not null"` // 开课学院ID
	College   *College // 开课学院

	CourseSpecifics []*CourseSpecific // 课程具体信息
}

// CourseSpecific 具体课头，指具体开给某一个老师的课程
// TODO 考虑记录课头状态（可选课、已结课等）
type CourseSpecific struct {
	Model
	CourseCommonId  uint              `gorm:"not null"` // 课程公共信息ID
	CourseCommon    *CourseCommon     // 课程公共信息
	TeacherId       uint              `gorm:"not null"`             // 授课教师ID
	Teacher         *User             `gorm:"foreignKey:TeacherId"` // 授课教师
	Location        string            // 上课地点
	Quota           uint              // 人数配额
	QuotaUsed       uint              // 课程已选人数
	SemesterId      uint              // 学期ID
	Semester        *Semester         // 学期
	CourseSchedules []*CourseSchedule `gorm:"many2many:course_specific_course_schedule;"`
}

// CourseSchedule 上课时间
type CourseSchedule struct {
	Model
	DayOfWeek   uint // 每周第几天
	HoursId     uint // 第几节课
	HoursCount  uint // 课程时长
	StartWeekId uint // 起始周次
	EndWeekId   uint // 结束周次
}

// StudentCourse 学生和具体某门课的关系
type StudentCourse struct {
	Model
	StudentId    uint            // 学生ID
	Student      User            `gorm:"foreignKey:StudentId"` // 学生
	CourseId     uint            // 课程ID
	Course       *CourseSpecific `gorm:"foreignKey:CourseId"` // 具体课程
	CourseStatus CourseStatus    // 课程状态
	Score        float32         // 学生成绩
}

type CourseScheduleWithCourseSpecific struct {
	CourseSchedule
	CourseSpecific   *CourseSpecific `gorm:"foreignKey:CourseSpecificId"`
	CourseSpecificId uint
}

func CourseSchedule_CourseSpecific(tx *gorm.DB) *gorm.DB {
	return tx.Table("course_schedules").
		Joins("JOIN course_specific_course_schedule ON course_specific_course_schedule.course_schedule_id = course_schedules.id").
		Joins("JOIN course_specifics ON course_specifics.id = course_specific_id").
		Select("course_schedules.*, course_specific_course_schedule.course_specific_id")
}

func CourseSchedule_CourseSpecific_Student(tx *gorm.DB) *gorm.DB {
	return CourseSchedule_CourseSpecific(tx).
		Joins("JOIN student_courses ON student_courses.course_id = course_specific_course_schedule.course_specific_id")
}

type CourseSpecificWithStudent struct {
	CourseSpecific
	StudentCourses []*StudentCourse `gorm:"many2many:student_courses;"`
}
