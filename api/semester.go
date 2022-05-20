package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/model/req"
	"github.com/se2022-qiaqia/course-system/model/resp"
	S "github.com/se2022-qiaqia/course-system/services"
)

type Semester struct{}

// QuerySemester
// @Summary					获取学期。
// @Tags					学期
// @Accept					json
// @Produce					json
// @Security				ApiKeyAuth
// @Success 				200 			{array} resp.Semester
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/semester [get]
func (api *Semester) QuerySemester(c *gin.Context) {
	if semesters, err := S.Services.Semester.GetSemesters(); err == nil {
		ss := make([]*resp.Semester, len(semesters))
		for i, semester := range semesters {
			ss[i] = resp.NewSemester(semester)
		}
		resp.Ok(semesters, c)
		return
	} else {
		resp.FailJust("获取失败", c)
		return
	}
}

// CreateSemester
// @Summary					创建学期。
// @Tags					学期
// @Accept					json
// @Produce					json
// @Security				ApiKeyAuth
// @Param					new				body	 req.Semester	true		"学期信息"
// @Success 				200 			{object} resp.Semester "创建好的学期"
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/semester [post]
func (api *Semester) CreateSemester(c *gin.Context) {
	var reqSemester req.Semester
	if !req.BindAndValidate(c, &reqSemester) {
		return
	}

	if semester, err := S.Services.Semester.CreateSemester(&reqSemester); err == nil {
		resp.Ok(resp.NewSemester(semester), c)
		return
	} else if errors.Is(err, S.ErrConflict) {
		resp.FailJust("该学期冲突", c)
		return
	} else {
		resp.FailJust("创建失败", c)
	}
}

// GetCurrentSemester
// @Summary					获取当前学期。
// @Tags					学期
// @Accept					json
// @Produce					json
// @Security				ApiKeyAuth
// @Success 				200 			{object} 	resp.Semester "当前学期"
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/semester/curr [get]
func (api *Semester) GetCurrentSemester(c *gin.Context) {
	if semester, err := S.Services.Semester.GetCurrentSemester(); err == nil {
		resp.Ok(resp.NewSemester(semester), c)
		return
	} else if errors.Is(err, S.ErrNotFound) {
		resp.FailJust("未设置当前学期，请联系管理员", c)
		return
	} else {
		resp.FailJust("获取失败", c)
	}
}

// SetCurrentSemester
// @Summary					设置当前学期。
// @Tags					学期
// @Accept					json
// @Produce					json
// @Security				ApiKeyAuth
// @Param					id				body		req.IdReq		true		"学期id"
// @Success 				200 			{object} 	resp.Semester 	"更新后的当前学期"
// @Failure 				400 			{object} 	resp.ErrorResponse
// @Router					/semester/curr [post]
func (api *Semester) SetCurrentSemester(c *gin.Context) {
	var reqId req.IdReq
	if !req.BindAndValidate(c, &reqId) {
		return
	}

	if semester, err := S.Services.Semester.SetCurrentSemester(reqId.Id); err == nil {
		resp.Ok(resp.NewSemester(semester), c)
		return
	} else if errors.Is(err, S.ErrNotFound) {
		resp.FailJust("无法找到对应学期", c)
		return
	} else /* if errors.Is(err,S.ErrUpdateFailed) */ {
		resp.FailJust("更新当前学期失败", c)
	}
}
