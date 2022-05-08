package colleges

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api/resp"
	"github.com/se2022-qiaqia/course-system/services"
	"net/http"
)

func ListColleges(c *gin.Context) {
	var b services.QueryCollegesService
	if err := c.BindJSON(&b); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Response{Msg: "请输入正确的参数"})
		return
	}
	c.JSON(http.StatusOK, resp.Response{Data: b.Query()})
}

func NewCollege(c *gin.Context) {
	var b services.NewCollegeService
	if err := c.BindJSON(&b); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Response{Msg: "请输入正确的参数"})
		return
	}

	if college, err := b.NewCollege(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Response{Msg: "创建失败"})
		return
	} else {
		c.JSON(http.StatusOK, resp.Response{Data: college})
	}
}
