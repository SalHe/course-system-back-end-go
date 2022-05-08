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
