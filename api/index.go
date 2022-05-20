package api

type RootApi struct {
	Public
	Start
	User
	Course
	College
	Semester
}

var Api = &RootApi{}
