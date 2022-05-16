package req

type LoginCredit struct {
	Username string `json:"username" binding:"required" description:"用户名"`
	Password string `json:"password" binding:"required" description:"密码"`
}

type RegisterInfo struct {
	Username string `json:"username" binding:"required,username" description:"用户名"`
	Password string `json:"password" binding:"required,password" description:"密码"`
	Id       uint   `json:"id"`
}

// InitRequest 初始化系统参数
type InitRequest struct {
	Id       uint   `json:"id" binding:"required" description:"管理员ID"`                      // 管理员ID
	Username string `json:"username" binding:"required,username" description:"管理员用户名"`      // 管理员用户名
	Password string `json:"password" binding:"required,password" description:"管理员密码"`       // 管理员密码
	RealName string `json:"realName" binding:"required,min=2,max=10" description:"管理员真实姓名"` // 管理员真实姓名
}

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
	DayOfWeek uint `json:"dayOfWeek" binding:"required,min=0,max=6"` // 每周第几天
	HoursId   uint `json:"hoursId" binding:"required,min=1"`         // 第几节课
}

// OpenCourseRequest 开设课头。
type OpenCourseRequest struct {
	CourseCommonId  uint              `json:"courseCommonId" binding:"required" description:"课程id"` // 课程id
	SemesterId      uint              `json:"semesterId" binding:"required" description:"学期id"`     // 学期id
	TeacherId       uint              `json:"teacherId" binding:"required" description:"教师id"`      // 教师id
	Location        string            `json:"location" binding:"required" description:"上课地点"`       // 上课地点
	Quota           uint              `json:"quota" binding:"required" description:"容量"`            // 容量
	CourseSchedules []*CourseSchedule `json:"courseSchedules" description:"上课时间"`                   // 上课时间
}

// NewCollegeService 新学院信息
type NewCollegeService struct {
	Name string `json:"name" binding:"required" description:"学院名称"` // 学院名称
}

// QueryCollegesService 查询学院信息的筛选条件
type QueryCollegesService struct {
	Name string `json:"name" description:"学院名称"` // 学院名称，将会模糊搜索
}
