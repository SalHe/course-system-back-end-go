package api

type RootApi struct {
	Public
	Start
	User
	Course
	College
}

var Api = &RootApi{}
