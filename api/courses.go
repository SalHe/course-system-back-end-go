package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/middleware"
	"github.com/se2022-qiaqia/course-system/model/req"
	"github.com/se2022-qiaqia/course-system/model/resp"
	S "github.com/se2022-qiaqia/course-system/services"
	"github.com/se2022-qiaqia/course-system/token"
	"github.com/se2022-qiaqia/course-system/utils"
	"gorm.io/gorm"
	"strconv"
)

type Course struct{}

// GetCourseList
// @Summary					获取课程列表。
// @Description				可根据课程名关键字、学期、教师名字、学院ID筛选，字符串类参数为模糊搜索，填空代表不筛选对应条件。
// @Tags					课程
// @Accept					json
// @Produce					json
// @Param					params			body		req.QueryCoursesRequest	true	"筛选条件"
// @Security				ApiKeyAuth
// @Success 				200 			{array} 	resp.CourseCommonWithSpecifics
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/course/list 	[post]
func (api *Course) GetCourseList(c *gin.Context) {
	var b req.QueryCoursesRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	if courseCommons, err := S.Services.Course.Query(b); err != nil {
		resp.FailJust("查询失败", c)
		return
	} else {
		results := make([]*resp.CourseCommonWithSpecifics, len(courseCommons))
		for i, courseCommon := range courseCommons {
			results[i] = resp.NewCourseCommonWithSpecifics(courseCommon)
		}
		resp.Ok(results, c)
		return
	}
}

// NewCourse
// @Summary					添加课程。
// @Description				添加新课程。
// @Tags					课程
// @Accept					json
// @Produce					json
// @Param					params			body		req.NewCourseRequest	true	"课程信息"
// @Security				ApiKeyAuth
// @Success 				200 			{object} 	resp.CourseCommon
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/course		 	[post]
func (api *Course) NewCourse(c *gin.Context) {
	var b req.NewCourseRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	if courseCommon, err := S.Services.Course.NewCourse(b); err != nil {
		resp.FailJust("创建失败", c)
		return
	} else {
		resp.Ok(resp.NewCourseCommon(courseCommon), c)
		return
	}
}

// OpenCourse
// @Summary					开课。
// @Description				添加新课头。
// @Tags					课程
// @Accept					json
// @Produce					json
// @Param					params			body		req.OpenCourseRequest	true	"课程信息"
// @Security				ApiKeyAuth
// @Success 				200 			{object} 	resp.CourseSpecific
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/course/open	[post]
func (api *Course) OpenCourse(c *gin.Context) {
	var b req.OpenCourseRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	var teacher dao.User
	if err := dao.DB.Model(&dao.User{}).Where("id = ?", b.TeacherId).First(&teacher).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(resp.ErrCodeNotFound, fmt.Sprintf("未找到对应教师: id=%v", b.TeacherId), c)
		return
	}
	if teacher.Role != dao.RoleTeacher {
		resp.Fail(resp.ErrCodeNotFound, fmt.Sprintf("对应用户不是教师: id=%v", b.TeacherId), c)
		return
	}

	if utils.IsScheduleConflict(b.CourseSchedules, false) {
		resp.FailJust("课程时间冲突", c)
		return
	}

	if course, err := S.Services.Course.OpenCourse(b); err != nil {
		resp.FailJust("开课失败", c)
		return
	} else {
		resp.Ok(resp.NewCourseSpecific(&course), c)
		return
	}
}

// UpdateCourseCommon
// @Summary					更新课程。
// @Description				更新课程的公共信息。
// @Tags					课程
// @Accept					json
// @Produce					json
// @Param 					id				path		int 							true 	"课程ID"
// @Param					new				body		req.UpdateCourseCommonRequest	true	"新课程信息"
// @Security				ApiKeyAuth
// @Success 				200 			{array} 	resp.CourseCommon	"更新后的课程信息"
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/course/{id} 	[put]
func (api *Course) UpdateCourseCommon(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var b req.UpdateCourseCommonRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	if updated, err := S.Services.Course.UpdateCourseCommon(uint(id), b); err == nil {
		resp.Ok(resp.NewCourseCommon(updated), c)
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, S.ErrNotFound) {
		resp.Fail(resp.ErrCodeNotFound, fmt.Sprintf("未找到对应课程: %v", id), c)
		return
	} else {
		resp.FailJust("更新失败", c)
		return
	}
}

// UpdateCourseSpecific
// @Summary					更新课头。
// @Description
// @Tags					课程
// @Accept					json
// @Produce					json
// @Param 					id				path		int 							true 	"课头ID"
// @Param					new				body		req.UpdateCourseSpecificRequest	true	"新课头信息"
// @Security				ApiKeyAuth
// @Success 				200 			{array} 	resp.CourseCommon	"更新后的课程信息"
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/course/spec/{id} 	[put]
func (api *Course) UpdateCourseSpecific(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var b req.UpdateCourseSpecificRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	// TODO 考虑抽取这部分逻辑
	var teacher dao.User
	if err := dao.DB.Model(&dao.User{}).Where("id = ?", b.TeacherId).First(&teacher).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(resp.ErrCodeNotFound, fmt.Sprintf("未找到对应教师: id=%v", b.TeacherId), c)
		return
	}
	if teacher.Role != dao.RoleTeacher {
		resp.Fail(resp.ErrCodeNotFound, fmt.Sprintf("对应用户不是教师: id=%v", b.TeacherId), c)
		return
	}

	if utils.IsScheduleConflict(b.CourseSchedules, false) {
		resp.FailJust("课程时间冲突", c)
		return
	}

	if updated, err := S.Services.Course.UpdateCourseSpecific(uint(id), b); err == nil {
		resp.Ok(resp.NewCourseSpecific(updated), c)
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, S.ErrNotFound) {
		resp.Fail(resp.ErrCodeNotFound, fmt.Sprintf("未找到对应课头: id=%v", id), c)
		return
	} else {
		resp.FailJust("更新失败", c)
		return
	}
}

// SelectCourse
// @Summary					选课。
// @Description
// @Tags					课程
// @Accept					json
// @Produce					json
// @Param					id				body		req.SelectCourseRequest	true "选课信息"
// @Security				ApiKeyAuth
// @Success 				200 			{array} 	resp.CourseSpecific	"选中的课头"
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/course/select 	[post]
func (api *Course) SelectCourse(c *gin.Context) {
	// TODO 待精心测试该API

	var b req.SelectCourseRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	cla, _ := c.Get(middleware.ClaimsKey)
	claims := cla.(*token.JwtClaims)

	if b.StudentId <= 0 {
		b.StudentId = claims.ID
	}

	if selected, err := S.Services.Course.SelectCourse(&b, claims); err == nil {
		resp.Ok(resp.NewCourseSpecific(selected), c)
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, S.ErrNotFound) {
		resp.Fail(resp.ErrCodeNotFound, fmt.Sprintf("未找到对应课头: id=%v", b.CourseId), c)
		return
	} else if errors.Is(err, S.ErrUnauthorized) {
		resp.Fail(resp.ErrCodeUnauthorized, "没有权限", c)
		return
	} else if errors.Is(err, S.ErrQuotaExceeded) {
		resp.Fail(resp.ErrCodeQuotaExceeded, "余量不足", c)
		return
	} else if errors.Is(err, S.ErrConflict) {
		resp.Fail(resp.ErrCodeConflict, "课程与当前课表时间冲突，或者您已选该课头", c)
		return
	}

	resp.FailJust("选课失败", c)
	return
}

// UnSelectCourse
// @Summary					撤课。
// @Description
// @Tags					课程
// @Accept					json
// @Produce					json
// @Param					id				body		req.SelectCourseRequest	true "撤课信息"
// @Security				ApiKeyAuth
// @Success 				200 			{array} 	resp.CourseSpecific	"撤掉的课头信息"
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/course/select 	[delete]
func (api *Course) UnSelectCourse(c *gin.Context) {
	// TODO 待精心测试该API

	var b req.SelectCourseRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	cla, _ := c.Get(middleware.ClaimsKey)
	claims := cla.(*token.JwtClaims)

	if b.StudentId <= 0 {
		b.StudentId = claims.ID
	}

	if unselected, err := S.Services.Course.UnSelectCourse(&b, claims); err == nil {
		resp.Ok(resp.NewCourseSpecific(unselected), c)
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, S.ErrNotFound) {
		resp.Fail(resp.ErrCodeNotFound, fmt.Sprintf("未找到对应课头: id=%v", b.CourseId), c)
		return
	} else if errors.Is(err, S.ErrUnauthorized) {
		resp.Fail(resp.ErrCodeUnauthorized, "没有权限", c)
		return
	}

	resp.FailJust("撤课失败", c)
	return
}

// GetCourseSchedules
// @Summary					查询课程安排。
// @Description
// @Tags					课程
// @Accept					json
// @Produce					json
// @Param					params				body		req.GetSchedulesRequest	true "查询信息"
// @Security				ApiKeyAuth
// @Success 				200 				{array} 	dao.CourseScheduleWithCourseSpecific
// @Failure 				400 				{object} 	resp.ErrorResponse
// @Router					/course/schedules 	[post]
func (api *Course) GetCourseSchedules(c *gin.Context) {
	var b req.GetSchedulesRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	cla, _ := c.Get(middleware.ClaimsKey)
	claims := cla.(*token.JwtClaims)

	if b.UserId <= 0 {
		b.UserId = claims.ID
	}

	if b.UserId != claims.ID && !claims.IsAdmin() {
		resp.Fail(resp.ErrCodeUnauthorized, "您不可查询他人课程安排", c)
		return
	}

	var user dao.User
	if err := dao.FindUserById(dao.DB, b.UserId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Fail(resp.ErrCodeNotFound, "未找到对应用户", c)
			return
		}
		resp.Fail(resp.ErrCodeNotFound, "查询用户失败", c)
		return
	}

	if b.UserId == claims.ID && claims.IsAdmin() {
		// 管理员自身不会有课程安排
		resp.Ok([]dao.CourseScheduleWithCourseSpecific{}, c)
		return
	}

	var schedules []*dao.CourseScheduleWithCourseSpecific
	var err error
	if user.Role == dao.RoleStudent {
		schedules, err = S.Services.Course.GetStudentSchedules(&b)
	} else if user.Role == dao.RoleAdmin {
		schedules, err = S.Services.Course.GetTeacherSchedules(&b)
	}

	if err == nil {
		res := make([]*resp.CourseScheduleWithCourseSpecific, len(schedules))
		for i, schedule := range schedules {
			res[i] = &resp.CourseScheduleWithCourseSpecific{
				CourseSchedule: resp.CourseSchedule{
					Model:        schedule.Model,
					DayOfWeek:    schedule.DayOfWeek,
					StartHoursId: schedule.HoursId,
					EndHoursId:   schedule.HoursId + schedule.HoursCount - 1,
					StartWeekId:  schedule.StartWeekId,
					EndWeekId:    schedule.EndWeekId,
				},
				CourseSpecific: resp.NewCourseSpecific(schedule.CourseSpecific),
			}
		}
		resp.Ok(res, c)
		return
	} else {
		resp.FailJust("查询课程安排失败", c)
		return
	}
}
