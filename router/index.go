package router

import "github.com/gin-gonic/gin"

type IRouter interface {
	Init(r *gin.RouterGroup)
}

type RootRouter struct {
	Public
	Course
	Start
	User
	College
	Semester
}

var Router = &RootRouter{}
