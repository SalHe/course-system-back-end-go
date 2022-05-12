package api

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/model/req"
	"github.com/se2022-qiaqia/course-system/model/resp"
	"github.com/se2022-qiaqia/course-system/services"
)

type Course struct{}

func (api Course) GetCourseList(c *gin.Context) {
	var b services.QueryCoursesServices
	if !req.BindAndValidate(c, &b) {
		return
	}

	if courseCommons, err := b.Query(); err != nil {
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

func (api Course) NewCourse(c *gin.Context) {
	var b services.NewCourseService
	if !req.BindAndValidate(c, &b) {
		return
	}

	if courseCommon, err := b.NewCourse(); err != nil {
		resp.FailJust("创建失败", c)
		return
	} else {
		resp.Ok(resp.NewCourseCommon(courseCommon), c)
		return
	}
}

func (api Course) OpenCourse(c *gin.Context) {
	var b services.OpenCourseService
	if !req.BindAndValidate(c, &b) {
		return
	}

	if course, err := b.OpenCourse(); err != nil {
		resp.FailJust("开课失败", c)
		return
	} else {
		resp.Ok(resp.NewCourseSpecific(&course), c)
		return
	}
}

func (api Course) UpdateCourse(c *gin.Context) {
	// TODO 实现 UpdateCourse
}
