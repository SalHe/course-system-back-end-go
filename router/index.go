package router

// type IRouter interface {
// 	Init(r *gin.RouterGroup)
// }

type RootRouter struct {
	Public
	Course
	Start
	User
	College
	Semester
}

var Router = &RootRouter{}
