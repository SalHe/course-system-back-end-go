package services

import (
	"github.com/se2022-qiaqia/course-system/dao"
	"gorm.io/gorm/clause"
)

type QueryCoursesServices struct {
	Page
	Name        string `json:"name"`
	Semester    uint   `json:"semester"`
	TeacherName string `json:"teacherName"`
	CollegesId  []uint `json:"collegesId"`
}

func (q QueryCoursesServices) Query() (courseCommons []*dao.CourseCommon, err error) {
	db := dao.DB.Preload("CourseSpecifics").Preload("College").Model(&dao.CourseCommon{}).
		Preload("Teacher").Model(&dao.CourseSpecific{})
	if len(q.Name) > 0 {
		db = db.Where("name like (?)", "%"+q.Name+"%")
	}
	if len(q.CollegesId) > 0 {
		db = db.Where("college_id in (?)", q.CollegesId)
	}

	{
		var conditions []interface{}
		if q.Semester > 0 {
			conditions = append(conditions, "semester_id = ?", q.Semester)
		}
		if len(q.TeacherName) > 0 {
			conditions = append(conditions, "teacher_id in (?)", dao.DB.Table("users").Where("real_name like ?", "%"+q.TeacherName+"%").Select("id"))
		}
		if len(conditions) > 0 {
			db = db.Preload("CourseSpecifics", conditions...)
		}
	}

	err = db.Offset(q.Offset()).Limit(q.ActualSize()).Find(&courseCommons).Error
	return
}

type NewCourseService struct {
	Name      string  `json:"name"`
	CollegeId uint    `json:"collegeId"`
	Credits   float32 `json:"credits"`
	Hours     uint    `json:"hours"`
}

func (n NewCourseService) NewCourse() (courseCommon *dao.CourseCommon, err error) {
	courseCommon = &dao.CourseCommon{
		Name:      n.Name,
		Credits:   n.Credits,
		Hours:     n.Hours,
		CollegeId: n.CollegeId,
	}
	err = dao.DB.Create(courseCommon).Error
	if err == nil {
		dao.DB.Preload(clause.Associations).First(courseCommon)
	}
	return
}

type OpenCourseService struct {
	CourseCommonId  uint                  `json:"courseCommonId"`
	SemesterId      uint                  `json:"semesterId"`
	TeacherId       uint                  `json:"teacherId"`
	Location        string                `json:"location"`
	Quota           uint                  `json:"quota"`
	CourseSchedules []*dao.CourseSchedule `json:"courseSchedules"`
}

func (o OpenCourseService) OpenCourse() (course dao.CourseSpecific, err error) {
	course = dao.CourseSpecific{
		CourseCommonId:  o.CourseCommonId,
		TeacherId:       o.TeacherId,
		Location:        o.Location,
		Quota:           o.Quota,
		QuotaUsed:       0,
		SemesterId:      o.SemesterId,
		CourseSchedules: o.CourseSchedules,
	}
	err = dao.DB.Create(&course).Error
	if err == nil {
		err = dao.DB.Model(&course).Association("CourseSchedules").Append(o.CourseSchedules)
		if err == nil {
			var t dao.CourseSpecific
			dao.DB.Preload(clause.Associations).Find(&t, course.ID)
			course = t
		}
	}
	return
}
