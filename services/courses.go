package services

import (
	"github.com/se2022-qiaqia/course-system/dao"
	"gorm.io/gorm/clause"
)

type QueryCoursesServices struct {
	Page
	Name        string `json:"name" description:"课程名称"`
	Semester    uint   `json:"semester" description:"学期"`
	TeacherName string `json:"teacherName" description:"教师名称"`
	CollegesId  []uint `json:"collegesId" description:"包含的学院id"`
}

func (q QueryCoursesServices) Query() (courseCommons []*dao.CourseCommon, err error) {
	db := dao.DB.Preload("CourseSpecifics").Preload("College").Model(&dao.CourseCommon{})
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
			db = db.Preload("CourseSpecifics", conditions...).
				Preload("CourseSpecifics." + clause.Associations).
				Preload("CourseSpecifics.Teacher.College")
		}
	}

	err = db.Offset(q.Offset()).Limit(q.ActualSize()).Find(&courseCommons).Error
	return
}

type NewCourseService struct {
	Name      string  `json:"name" description:"课程名称"`
	CollegeId uint    `json:"collegeId" description:"学院id"`
	Credits   float32 `json:"credits" description:"学分"`
	Hours     uint    `json:"hours" description:"学时"`
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
	CourseCommonId  uint                  `json:"courseCommonId" description:"课程id"`
	SemesterId      uint                  `json:"semesterId" description:"学期id"`
	TeacherId       uint                  `json:"teacherId" description:"教师id"`
	Location        string                `json:"location" description:"上课地点"`
	Quota           uint                  `json:"quota" description:"容量"`
	CourseSchedules []*dao.CourseSchedule `json:"courseSchedules" description:"上课时间"`
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
			dao.DB.Preload(clause.Associations).
				Preload("Teacher.College").
				Preload("CourseCommon.College").
				Find(&t, course.ID)
			course = t
		}
	}
	return
}
