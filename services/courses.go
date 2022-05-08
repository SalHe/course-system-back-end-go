package services

import (
	"github.com/se2022-qiaqia/course-system/dao"
	"gorm.io/gorm/clause"
)

type QueryCoursesServices struct {
	Page
	Name        *string `json:"name"`
	Semester    *uint   `json:"semester"`
	TeacherName *string `json:"teacherName"`
	CollegesId  *[]uint `json:"collegesId"`
	CommonId    *uint   `json:"commonId"`
}

func (q QueryCoursesServices) Query() ([]*dao.CourseCommon, error) {
	var courseCommons []*dao.CourseCommon
	// TODO 条件筛选
	err := dao.DB.Preload(clause.Associations).Offset(q.Offset()).Limit(q.Size).Find(&courseCommons).Error
	return courseCommons, err
}
