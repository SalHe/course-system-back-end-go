package api

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/model/req"
	"github.com/se2022-qiaqia/course-system/model/resp"
	S "github.com/se2022-qiaqia/course-system/services"
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

// UpdateCourse TODO
// @Summary					TODO。
// @Description				TODO。
// @Tags					课程
// @Accept					json
// @Produce					json
// @Param					new				body		resp.ErrorResponse	true	"TODO"
// @Security				ApiKeyAuth
// @Success 				200 			{array} 	resp.ErrorResponse TODO
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/course		 	[put]
func (api Course) UpdateCourse(c *gin.Context) {
	// TODO 实现 UpdateCourse
}
