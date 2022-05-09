package courses

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api/resp"
	"github.com/se2022-qiaqia/course-system/services"
	"net/http"
)

func GetCourseList(c *gin.Context) {
	var b services.QueryCoursesServices
	if err := c.BindJSON(&b); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Response{Msg: "请输入参数"})
		return
	}
	if courseCommons, err := b.Query(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Response{Msg: "查询失败"})
		return
	} else {
		c.JSON(http.StatusOK, resp.Response{Data: courseCommons})
		return
	}
}

func NewCourse(c *gin.Context) {
	var b services.NewCourseService
	if err := c.BindJSON(&b); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Response{Msg: "请输入正确的参数"})
		return
	}
	if courseCommon, err := b.NewCourse(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Response{Msg: "创建失败"})
		return
	} else {
		c.JSON(http.StatusOK, resp.Response{Msg: "创建成功", Data: courseCommon})
		return
	}
}

func OpenCourse(c *gin.Context) {
	var b services.OpenCourseService
	if err := c.BindJSON(&b); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Response{Msg: "请输入正确的参数"})
		return
	}
	if course, err := b.OpenCourse(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Response{Msg: "开课失败"})
		return
	} else {
		c.JSON(http.StatusOK, resp.Response{Msg: "开课成功", Data: course})
		return
	}
}

func UpdateCourse(c *gin.Context) {}
