package api

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/model/req"
	"github.com/se2022-qiaqia/course-system/model/resp"
	"github.com/se2022-qiaqia/course-system/services"
)

type College struct{}

func (api College) ListColleges(c *gin.Context) {
	var b services.QueryCollegesService
	if !req.BindAndValidate(c, &b) {
		return
	}

	results := b.Query()
	actual := make([]*resp.College, len(results))
	for i, result := range results {
		actual[i] = resp.NewCollege(&result)
	}
	resp.Ok(actual, c)
	return
}

func (api College) NewCollege(c *gin.Context) {
	var b services.NewCollegeService
	if !req.BindAndValidate(c, &b) {
		return
	}

	if college, err := b.NewCollege(); err != nil {
		resp.FailJust("创建失败", c)
		return
	} else {
		resp.Ok(resp.NewCollege(college), c)
		return
	}
}
