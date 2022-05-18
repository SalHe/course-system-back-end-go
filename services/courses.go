package services

import (
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/model/req"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Course struct{}

func (c Course) Query(q req.QueryCoursesRequest) (courseCommons []*dao.CourseCommon, err error) {
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

func (c Course) NewCourse(n req.NewCourseRequest) (courseCommon *dao.CourseCommon, err error) {
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

func (c Course) OpenCourse(o req.OpenCourseRequest) (course dao.CourseSpecific, err error) {
	schedules := newDaoCourseSchedule(o.CourseSchedules)

	course = dao.CourseSpecific{
		CourseCommonId:  o.CourseCommonId,
		TeacherId:       o.TeacherId,
		Location:        o.Location,
		Quota:           o.Quota,
		QuotaUsed:       0,
		SemesterId:      o.SemesterId,
		CourseSchedules: schedules,
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

func newDaoCourseSchedule(o []*req.CourseSchedule) []*dao.CourseSchedule {
	var schedules []*dao.CourseSchedule
	for _, s := range o {
		schedules = append(schedules, &dao.CourseSchedule{
			DayOfWeek:   s.DayOfWeek,
			HoursId:     s.HoursId,
			HoursCount:  s.HoursCount,
			StartWeekId: s.StartWeekId,
			EndWeekId:   s.EndWeekId,
		})
	}
	return schedules
}

func (c Course) UpdateCourseCommon(id uint, b req.UpdateCourseCommonRequest) (*dao.CourseCommon, error) {
	var course dao.CourseCommon
	if err := dao.DB.Model(&dao.CourseCommon{}).Where("id = ?", id).First(&course).Error; err != nil {
		return nil, err
	}

	course.Name = b.Name
	course.CollegeId = b.CollegeId
	course.Credits = b.Credits
	course.Hours = b.Hours

	if err := dao.DB.Model(&course).Updates(course).Error; err != nil {
		return nil, err
	}

	dao.DB.Preload(clause.Associations).Model(&dao.CourseCommon{}).First(&course)

	return &course, nil
}

func (c Course) UpdateCourseSpecific(u uint, b req.UpdateCourseSpecificRequest) (*dao.CourseSpecific, error) {
	var course dao.CourseSpecific
	if err := dao.DB.Model(&dao.CourseSpecific{}).Where("id = ?", u).First(&course).Error; err != nil {
		return nil, err
	}

	course.TeacherId = b.TeacherId
	course.Location = b.Location
	course.Quota = b.Quota

	// TODO 处理课程安排冲突问题

	if err := dao.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&course).Updates(course).Error; err != nil {
			return err
		}

		if err := tx.Model(&course).Association("CourseSchedules").Replace(newDaoCourseSchedule(b.CourseSchedules)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	dao.DB.
		Preload(clause.Associations).
		Preload("Teacher.College").
		Preload("CourseCommon.College").
		Model(&dao.CourseSpecific{}).First(&course)
	return &course, nil
}
