package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/model/req"
	"github.com/se2022-qiaqia/course-system/model/resp"
	S "github.com/se2022-qiaqia/course-system/services"
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
func (api Course) GetCourseList(c *gin.Context) {
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
func (api Course) NewCourse(c *gin.Context) {
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
func (api Course) OpenCourse(c *gin.Context) {
	var b req.OpenCourseRequest
	if !req.BindAndValidate(c, &b) {
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
func (api Course) UpdateCourseCommon(c *gin.Context) {
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
func (api Course) UpdateCourseSpecific(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var b req.UpdateCourseSpecificRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	var count int64
	if err := dao.DB.Model(&dao.User{}).Where("id = ?", b.TeacherId).Count(&count).Error; errors.Is(err, gorm.ErrRecordNotFound) || count == 0 {
		resp.Fail(resp.ErrCodeNotFound, fmt.Sprintf("未找到对应教师: id=%v", b.TeacherId), c)
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
