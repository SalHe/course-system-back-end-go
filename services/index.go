package services

type RootServices struct {
	Public
	Start
	Course
	User
	College
}

var Services = &RootServices{}
