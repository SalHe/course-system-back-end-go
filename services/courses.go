package services

import (
	"errors"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/model/req"
	"github.com/se2022-qiaqia/course-system/token"
	"github.com/se2022-qiaqia/course-system/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Course struct{}

func (c *Course) Query(q req.QueryCoursesRequest) (count int64, courseCommons []*dao.CourseCommon, err error) {
	db := dao.DB.Preload("CourseSpecifics").Preload("College").Model(&dao.CourseCommon{})
	if len(q.Name) > 0 {
		db = db.Where("name like (?)", "%"+q.Name+"%")
	}
	if len(q.CollegesId) > 0 {
		db = db.Where("college_id in (?)", q.CollegesId)
	}

	{
		// TODO 待优化
		query := ""
		var args []interface{}
		joinQuery := func(q string, arg ...interface{}) {
			if len(query) > 0 {
				query += " and "
			}
			query += q
			args = append(args, arg...)
		}
		if q.Semester > 0 {
			joinQuery("semester_id = ?", q.Semester)
		}
		if len(q.TeacherName) > 0 {
			joinQuery("teacher_id in (?)", dao.DB.Table("users").Where("real_name like ?", "%"+q.TeacherName+"%").Select("id"))
		}
		if q.OnlyLeftQuota {
			joinQuery("quota_used < quota")
		}
		if len(query) > 0 {
			fargs := append([]interface{}{}, query)
			fargs = append(fargs, args...)
			db = db.Preload("CourseSpecifics", fargs...).
				Preload("CourseSpecifics." + clause.Associations).
				Preload("CourseSpecifics.Teacher.College")
		}
	}

	db.Count(&count)
	err = db.Offset(q.Offset()).Limit(q.ActualSize()).Find(&courseCommons).Error
	return
}

func (c *Course) NewCourse(n req.NewCourseRequest) (courseCommon *dao.CourseCommon, err error) {
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

func (c *Course) OpenCourse(o req.OpenCourseRequest) (course dao.CourseSpecific, err error) {
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
			HoursId:     s.StartHoursId,
			HoursCount:  s.EndHoursId - s.StartHoursId + 1,
			StartWeekId: s.StartWeekId,
			EndWeekId:   s.EndWeekId,
		})
	}
	return schedules
}

func (c *Course) UpdateCourseCommon(id uint, b req.UpdateCourseCommonRequest) (*dao.CourseCommon, error) {
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

func (c *Course) UpdateCourseSpecific(u uint, b req.UpdateCourseSpecificRequest) (*dao.CourseSpecific, error) {
	var course dao.CourseSpecific
	if err := dao.DB.Model(&dao.CourseSpecific{}).Where("id = ?", u).First(&course).Error; err != nil {
		return nil, err
	}

	course.TeacherId = b.TeacherId
	course.Location = b.Location
	course.Quota = b.Quota

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

func (c *Course) SelectCourse(b *req.SelectCourseRequest, operator *token.JwtClaims) (*dao.CourseSpecific, error) {
	// 需要判断的：
	// 0. 操作用户是否有权限（普通用户不得为他人选课）
	// 1. 判断是否已选同课程下的任何课头
	// 2. 是否超过限额
	// 3. 时间是否冲突

	// 选课具体进行的操作
	// 1. 添加选课记录
	// 2. 更新课程已用额度（管理员添加要不要影响额度？），考量了一下，管理员操作还是影响额度

	if !operator.IsAdmin() && operator.ID != b.StudentId {
		// 无权限
		return nil, ErrUnauthorized
	}

	var course dao.CourseSpecific
	if err := dao.DB.Model(&dao.CourseSpecific{}).Where("id = ?", b.CourseId).First(&course).Error; err != nil {
		return nil, err
	}

	var studentCourse dao.StudentCourse
	err := dao.DB.Table("student_courses").
		Select("student_courses.*").
		Joins("JOIN course_specifics ON course_specifics.id = student_courses.course_id").
		Where("student_courses.student_id = ? AND course_specifics.course_common_id = ?", b.StudentId, course.CourseCommonId).
		First(&studentCourse).Error
	if studentCourse.CourseId != b.CourseId || studentCourse.CourseStatus == dao.CourseStatusWithdraw || errors.Is(err, gorm.ErrRecordNotFound) {
		// 对应课程未选或已撤掉

		if !operator.IsAdmin() {
			var courseSpecific dao.CourseSpecific
			if err := dao.DB.Model(&dao.CourseSpecific{}).Where("id = ?", b.CourseId).Find(&courseSpecific).Error; err != nil {
				return nil, err
			} else if courseSpecific.QuotaUsed >= courseSpecific.Quota {
				return nil, ErrQuotaExceeded
			}
		}
		// 余量充足（或管理员选课不考虑余量）

		var schedules []*dao.CourseSchedule
		var targetSchedule []*dao.CourseSchedule
		tx := dao.DB.Table("course_schedules").Select("course_schedules.*").
			Joins("JOIN course_specific_course_schedule ON course_specific_course_schedule.course_schedule_id = course_schedules.id").
			Joins("JOIN student_courses ON student_courses.course_id = course_specific_course_schedule.course_specific_id AND student_courses.course_status IN (?)", []dao.CourseStatus{dao.CourseStatusNormal})
		tx.Where("student_courses.student_id = ?", b.StudentId).
			Find(&schedules)
		tx.Where("course_specific_course_schedule.course_specific_id = ?", b.CourseId).
			Find(&targetSchedule)

		schedulesAll := make([]*req.CourseSchedule, len(schedules)+len(targetSchedule))
		for i, schedule := range schedules {
			schedulesAll[i] = req.NewCourseSchedule(schedule)
		}
		for i, schedule := range targetSchedule {
			schedulesAll[i+len(schedules)] = req.NewCourseSchedule(schedule)
		}

		if utils.IsScheduleConflict(schedulesAll, false) {
			// 时间安排冲突
			return nil, ErrConflict
		}

		// 权限足、没选过、余量足、不冲突
		if err := dao.DB.Transaction(func(tx *gorm.DB) error {
			if studentCourse.ID > 0 && studentCourse.CourseId == b.CourseId {
				studentCourse.CourseStatus = dao.CourseStatusNormal
				if err := tx.Save(&studentCourse).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Create(&dao.StudentCourse{
					StudentId: b.StudentId,
					CourseId:  b.CourseId,
				}).Error; err != nil {
					return err
				}
			}

			if err := tx.Model(&dao.CourseSpecific{}).Where("id = ?", b.CourseId).Update("quota_used", gorm.Expr("quota_used + 1")).Error; err != nil {
				return err
			}

			return nil
		}); err != nil {
			return nil, err
		}
		dao.DB.Preload(clause.Associations).
			Preload("CourseCommon.College").
			Preload("Teacher.College").
			Model(&dao.CourseSpecific{}).Where("id = ?", b.CourseId).First(&course)
		return &course, nil
	} else {
		// 已选
		return nil, ErrConflict
	}

}

func (c *Course) UnSelectCourse(b *req.SelectCourseRequest, operator *token.JwtClaims) (*dao.CourseSpecific, error) {

	if !operator.IsAdmin() && operator.ID != b.StudentId {
		// 无权限
		return nil, ErrUnauthorized
	}

	var studentCourse dao.StudentCourse
	if err := dao.DB.Model(&dao.StudentCourse{}).Where("student_id = ? AND course_id = ?", b.StudentId, b.CourseId).First(&studentCourse).Error; err != nil {
		return nil, err
	}

	if studentCourse.CourseStatus == dao.CourseStatusWithdraw {
		return nil, ErrCannotOperate
	}

	if err := dao.DB.Transaction(func(tx *gorm.DB) error {
		studentCourse.CourseStatus = dao.CourseStatusWithdraw
		if err := tx.Save(studentCourse).Error; err != nil {
			return err
		}

		if err := tx.Model(&dao.CourseSpecific{}).Where("id = ?", b.CourseId).Update("quota_used", gorm.Expr("quota_used - 1")).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	var course dao.CourseSpecific
	dao.DB.Preload(clause.Associations).
		Preload("CourseCommon.College").
		Preload("Teacher.College").
		Model(&dao.CourseSpecific{}).Where("id = ?", b.CourseId).First(&course)
	return &course, nil
}

func (c *Course) GetStudentSchedules(b *req.GetSchedulesRequest) ([]*dao.CourseScheduleWithCourseSpecific, error) {
	tx := dao.CourseSchedule_CourseSpecific_Student(dao.DB)
	if len(b.SemesterIds) > 0 {
		tx = tx.Where("semester_id IN (?)", b.SemesterIds)
	}
	tx = tx.Select("course_schedules.*, course_specific_course_schedule.course_specific_id")

	tx = tx.Where("student_id = ? AND student_courses.course_status in (?)", b.UserId, []dao.CourseStatus{dao.CourseStatusNormal})

	var schedules []*dao.CourseScheduleWithCourseSpecific
	if err := tx.Preload(clause.Associations).
		Preload("CourseSpecific.CourseCommon").
		Preload("CourseSpecific.CourseCommon." + clause.Associations).
		Preload("CourseSpecific." + clause.Associations).
		Preload("CourseSpecific.Teacher." + clause.Associations).
		Preload("CourseSpecific.Semester").
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (c *Course) GetTeacherSchedules(b *req.GetSchedulesRequest) ([]*dao.CourseScheduleWithCourseSpecific, error) {
	tx := dao.CourseSchedule_CourseSpecific_Student(dao.DB)
	if len(b.SemesterIds) > 0 {
		tx = tx.Where("semester_id IN (?)", b.SemesterIds)
	}
	tx = tx.Select("course_schedules.*, course_specific_course_schedule.course_specific_id")

	tx = tx.Where("teacher_id = ?", b.UserId)

	var schedules []*dao.CourseScheduleWithCourseSpecific
	if err := tx.Preload(clause.Associations).
		Preload("CourseSpecific.CourseCommon").
		Preload("CourseSpecific.CourseCommon." + clause.Associations).
		Preload("CourseSpecific." + clause.Associations).
		Preload("CourseSpecific.Teacher." + clause.Associations).
		Preload("CourseSpecific.Semester").
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (c *Course) GetCourseSpecificDetails(id uint) (*dao.CourseSpecificWithStudent, error) {
	var res dao.CourseSpecificWithStudent
	if err := dao.DB.Preload(clause.Associations).
		Preload("CourseCommon.College").
		Preload("Teacher.College").
		Model(&dao.CourseSpecific{}).
		Where("id = ?", id).
		First(&res.CourseSpecific).Error; err != nil {
		return nil, err
	}

	if err := dao.DB.Model(&dao.StudentCourse{}).
		Preload(clause.Associations).
		Preload("Student."+clause.Associations).
		Where("course_id = ?", id).
		Find(&res.StudentCourses).Error; err != nil {
		return nil, err
	}

	return &res, nil
}
