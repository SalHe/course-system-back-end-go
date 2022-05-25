package req

type LoginCredit struct {
	Username string `json:"username" binding:"required" description:"用户名"`
	Password string `json:"password" binding:"required" description:"密码"`
}

type RegisterInfo struct {
	Username  string `json:"username" binding:"required,username" description:"用户名"`
	Password  string `json:"password" binding:"required,password" description:"密码"`
	CollegeId uint   `json:"collegeId" binding:"required" description:"学院id"`
	RealName  string `json:"realName" binding:"required,min=1,max=10" description:"真实姓名"`
	Id        uint   `json:"id"`
}

// InitRequest 初始化系统参数
type InitRequest struct {
	Id       uint   `json:"id" binding:"required" description:"管理员ID"`                      // 管理员ID
	Username string `json:"username" binding:"required,username" description:"管理员用户名"`      // 管理员用户名
	Password string `json:"password" binding:"required,password" description:"管理员密码"`       // 管理员密码
	RealName string `json:"realName" binding:"required,min=2,max=10" description:"管理员真实姓名"` // 管理员真实姓名
}

// NewCollegeService 新学院信息
type NewCollegeService struct {
	Name string `json:"name" binding:"required" description:"学院名称"` // 学院名称
}

// QueryCollegesService 查询学院信息的筛选条件
type QueryCollegesService struct {
	Name string `json:"name" description:"学院名称"` // 学院名称，将会模糊搜索
}

type IdReq struct {
	Id uint `json:"id" binding:"required"`
}

type EnableRegisterRequest struct {
	Enable bool `json:"enable" binding:"required"`
}
