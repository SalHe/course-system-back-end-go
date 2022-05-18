package resp

import "github.com/se2022-qiaqia/course-system/dao"

type User struct {
	dao.Model
	Username     string   `json:"username"`
	RealName     string   `json:"realName"`
	Role         dao.Role `json:"role"`
	College      College  `json:"college"`
	EntranceYear uint     `json:"entranceYear"`
}

func NewUser(user *dao.User) *User {
	return &User{
		Model:        user.Model,
		Username:     user.Username,
		RealName:     user.RealName,
		Role:         user.Role,
		College:      *NewCollege(&user.College),
		EntranceYear: user.EntranceYear,
	}
}

type College struct {
	dao.Model
	Name string `json:"name"`
}

func NewCollege(college *dao.College) *College {
	return &College{
		Model: college.Model,
		Name:  college.Name,
	}
}

type CourseCommon struct {
	dao.Model
	Name string `json:"name"` // 课程名

	Credits float32 `json:"credits"` // 学分
	Hours   uint    `json:"hours"`   // 学时，单位不一定是小时

	College College `json:"college"` // 开课学院
}

func NewCourseCommon(course *dao.CourseCommon) *CourseCommon {
	return &CourseCommon{
		Model:   course.Model,
		Name:    course.Name,
		Credits: course.Credits,
		Hours:   course.Hours,
		College: *NewCollege(&course.College),
	}
}

type CourseCommonWithSpecifics struct {
	dao.Model
	CourseCommon
	CourseSpecifics []CourseSpecificWithoutCommon `json:"courseSpecifics"`
}

func NewCourseCommonWithSpecifics(course *dao.CourseCommon) *CourseCommonWithSpecifics {
	specifics := make([]CourseSpecificWithoutCommon, len(course.CourseSpecifics))
	for i, specific := range course.CourseSpecifics {
		specifics[i] = *NewCourseSpecificWithoutCommon(&specific)
	}
	return &CourseCommonWithSpecifics{
		Model:           course.Model,
		CourseCommon:    *NewCourseCommon(course),
		CourseSpecifics: specifics,
	}
}

type CourseSpecificWithoutCommon struct {
	dao.Model
	Teacher         User             `json:"teacher"`   // 授课教师
	Location        string           `json:"location"`  // 上课地点
	Quota           uint             `json:"quota"`     // 人数配额
	QuotaUsed       uint             `json:"quotaUsed"` // 课程已选人数
	Semester        Semester         `json:"semester"`  // 学期
	CourseSchedules []CourseSchedule `json:"courseSchedules"`
}

func NewCourseSpecificWithoutCommon(course *dao.CourseSpecific) *CourseSpecificWithoutCommon {
	schedules := make([]CourseSchedule, len(course.CourseSchedules))
	for i, schedule := range course.CourseSchedules {
		schedules[i] = *NewCourseSchedule(schedule)
	}
	return &CourseSpecificWithoutCommon{
		Model:           course.Model,
		Teacher:         *NewUser(&course.Teacher),
		Location:        course.Location,
		Quota:           course.Quota,
		QuotaUsed:       course.QuotaUsed,
		Semester:        *NewSemester(&course.Semester),
		CourseSchedules: schedules,
	}
}

type CourseSpecific struct {
	CourseSpecificWithoutCommon
	CourseCommon CourseCommon `json:"courseCommon"` // 课程公共信息
}

func NewCourseSpecific(course *dao.CourseSpecific) *CourseSpecific {
	return &CourseSpecific{
		CourseSpecificWithoutCommon: *NewCourseSpecificWithoutCommon(course),
		CourseCommon:                *NewCourseCommon(&course.CourseCommon),
	}
}

type CourseSchedule struct {
	dao.Model    `json:"-"`
	DayOfWeek    uint `json:"dayOfWeek"`    // 每周第几天
	StartHoursId uint `json:"startHoursId"` // 第几节课开始（包含该节课）
	EndHoursId   uint `json:"endHoursId"`   // 第几节课结束（包含该节课）
	StartWeekId  uint `json:"startWeekId"`  // 起始周次
	EndWeekId    uint `json:"endWeekId"`    // 结束周次
}

func NewCourseSchedule(courseSchedule *dao.CourseSchedule) *CourseSchedule {
	return &CourseSchedule{
		Model:        courseSchedule.Model,
		DayOfWeek:    courseSchedule.DayOfWeek,
		StartHoursId: courseSchedule.HoursId,
		EndHoursId:   courseSchedule.HoursId + courseSchedule.HoursCount - 1,
		StartWeekId:  courseSchedule.StartWeekId,
		EndWeekId:    courseSchedule.EndWeekId,
	}
}

type Semester struct {
	dao.Model
	Year uint `json:"year"` // 年份
	Term uint `json:"term"` // 对应年份第几学期
}

func NewSemester(semester *dao.Semester) *Semester {
	return &Semester{
		Model: semester.Model,
		Year:  semester.Year,
		Term:  semester.Term,
	}
}
