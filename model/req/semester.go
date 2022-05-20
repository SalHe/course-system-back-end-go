package req

type Semester struct {
	Year uint `json:"year" binding:"required,min=1980"`    // 年份
	Term uint `json:"term" binding:"required,min=1,max=4"` // 对应年份第几学期
}
