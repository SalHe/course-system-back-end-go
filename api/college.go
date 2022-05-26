package api

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/model/req"
	"github.com/se2022-qiaqia/course-system/model/resp"
	S "github.com/se2022-qiaqia/course-system/services"
)

type College struct{}

// ListColleges
// @Summary					获取学院。
// @Description				获取学院，可以根据关键字模糊查询。
// @Tags					学院
// @Accept					json
// @Produce					json
// @Param					queryFilter		body	req.QueryCollegesService	true	"查询条件"
// @Success 				200 			{object} resp.OkResponse{data=[]resp.College}
// @Router					/college/list [post]
func (api *College) ListColleges(c *gin.Context) {
	var b req.QueryCollegesService
	if !req.BindAndValidate(c, &b) {
		return
	}

	results := S.Services.College.Query(b)
	actual := make([]*resp.College, len(results))
	for i, result := range results {
		actual[i] = resp.NewCollege(&result)
	}
	resp.Ok(actual, c)
	return
}

// NewCollege
// @Summary					添加学院。
// @Tags					学院
// @Accept					json
// @Produce					json
// @Param					new				body		req.QueryCollegesService	true	"新学院信息"
// @Security				ApiKeyAuth
// @Success 				200 			{object} 	resp.OkResponse{data=resp.College}
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/college/new 	[post]
func (api *College) NewCollege(c *gin.Context) {
	var b req.NewCollegeService
	if !req.BindAndValidate(c, &b) {
		return
	}

	if college, err := S.Services.College.NewCollege(b); err != nil {
		resp.FailJust("存在同名学院！", c)
		return
	} else {
		resp.Ok(resp.NewCollege(college), c)
		return
	}
}
