package req

// QueryCoursesRequest 查询课程列表
type QueryCoursesRequest struct {
	Page
	Name        string `json:"name" description:"课程名称"`          // 课程名称
	Semester    uint   `json:"semester" description:"学期"`        // 学期id
	TeacherName string `json:"teacherName" description:"教师名称"`   // 教师名称
	CollegesId  []uint `json:"collegesId" description:"包含的学院id"` // 包含的学院id
}

// NewCourseRequest 新增课程，对应于课程的公共信息部分。
type NewCourseRequest struct {
	Name      string  `json:"name" binding:"required" description:"课程名称"`      // 课程名称
	CollegeId uint    `json:"collegeId" binding:"required" description:"学院id"` // 学院id
	Credits   float32 `json:"credits" binding:"required" description:"学分"`     // 学分
	Hours     uint    `json:"hours" binding:"required" description:"学时"`       // 学时
}

type CourseSchedule struct {
	DayOfWeek    uint `json:"dayOfWeek" binding:"min=0,max=6,oneof=0 1 2 3 4 5 6" description:"每周第几天"`   // 每周第几天
	StartHoursId uint `json:"startHoursId" binding:"required,min=1,max=24" description:"第几节课开始（包含该节课）"`  // 第几节课开始（包含该节课）
	EndHoursId   uint `json:"endHoursId" binding:"required,min=1,max=24" description:"第几节课结束（包含该节课）"`    // 第几节课结束（包含该节课）
	StartWeekId  uint `json:"startWeekId" binding:"required,min=1" description:"起始周"`                    // 起始周次
	EndWeekId    uint `json:"endWeekId" binding:"required,min=1,gtefield=StartWeekId" description:"结束周"` // 结束周次
}

// OpenCourseRequest 开设课头。
type OpenCourseRequest struct {
	CourseCommonId  uint              `json:"courseCommonId" binding:"required" description:"课程id"` // 课程id
	SemesterId      uint              `json:"semesterId" binding:"required" description:"学期id"`     // 学期id
	TeacherId       uint              `json:"teacherId" binding:"required" description:"教师id"`      // 教师id
	Location        string            `json:"location" binding:"required" description:"上课地点"`       // 上课地点
	Quota           uint              `json:"quota" binding:"required" description:"容量"`            // 容量
	CourseSchedules []*CourseSchedule `json:"courseSchedules" binding:"dive" description:"上课时间"`    // 上课时间
}

// UpdateCourseCommonRequest 更新课程。
type UpdateCourseCommonRequest struct {
	Name      string  `json:"name" binding:"required" description:"课程名称"`      // 课程名称
	CollegeId uint    `json:"collegeId" binding:"required" description:"学院id"` // 学院id
	Credits   float32 `json:"credits" binding:"required" description:"学分"`     // 学分
	Hours     uint    `json:"hours" binding:"required" description:"学时"`       // 学时
}

type UpdateCourseSpecificRequest struct {
	TeacherId       uint              `json:"teacherId" binding:"required" description:"教师id"`   // 教师id
	Location        string            `json:"location" binding:"required" description:"上课地点"`    // 上课地点
	Quota           uint              `json:"quota" binding:"required" description:"容量"`         // 容量
	CourseSchedules []*CourseSchedule `json:"courseSchedules" binding:"dive" description:"上课时间"` // 上课时间
}
