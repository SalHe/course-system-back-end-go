package req

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/model/resp"
)

func BindAndValidate(c *gin.Context, obj interface{}) bool {
	// TODO 处理校验错误信息
	if err := c.BindJSON(obj); err != nil {
		resp.FailJust("请输入正确的参数", c)
		return false
	}
	return true
}
