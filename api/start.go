package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/model/req"
	"github.com/se2022-qiaqia/course-system/model/resp"
	S "github.com/se2022-qiaqia/course-system/services"
)

type Start struct{}

// IsInitialized
// @Summary					检测是否已初始化系统。
// @Description
// @Tags					初始化
// @Accept					json
// @Produce					json
// @Success 				200 			{object}	resp.Response{data=boolean}
// @Router					/init			[get]
func (api *Start) IsInitialized(c *gin.Context) {
	resp.Ok(S.Services.IsInitialized(), c)
}

// InitSystem
// @Summary					初始化系统。
// @Description
// @Tags					初始化
// @Accept					json
// @Produce					json
// @Param					params			body		req.InitRequest		true	"初始化信息"
// @Success 				200 			{object}	resp.Response{data=boolean}
// @Router					/init			[post]
func (api *Start) InitSystem(c *gin.Context) {
	var b req.InitRequest
	if !req.BindAndValidate(c, &b) {
		return
	}

	err := S.Services.Start.InitSystem(b)
	if err == nil {
		resp.Ok(true, c)
		return
	} else if errors.Is(err, S.ErrConflict) {
		resp.FailJust("系统已初始化", c)
		return
	} else if err != nil {
		resp.FailJust("初始化失败", c)
		return
	}
}
