package services

type RootServices struct {
	*Public
	*Start
	*Course
	*User
	*College
	*Semester
	Setting
}

var Services = &RootServices{
	Setting: &settingInDB{},
}
