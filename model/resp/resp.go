package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorResponse struct {
	Msg  string  `json:"msg"`
	Code ErrCode `json:"code"`
}

type OkResponse struct {
	Msg  string      `json:"msg"`
	Code ErrCode     `json:"code"`
	Data interface{} `json:"data"`
}

type ErrCode int

const (
	ErrCodeOk = -iota
	ErrCodeFail
	ErrCodeNotFound
	ErrCodeConflict
	ErrCodeUnauthorized
	ErrCodeInternal
	ErrCodeQuotaExceeded
)

func Ok(data interface{}, c *gin.Context) {
	Resp(data, "ok", c)
}

func Resp(data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, OkResponse{
		Msg:  msg,
		Code: ErrCodeOk,
		Data: data,
	})
}

func FailJust(msg string, c *gin.Context) {
	Fail(ErrCodeFail, msg, c)
}

func Fail(errCode ErrCode, msg string, c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, ErrorResponse{
		Msg:  msg,
		Code: errCode,
	})
}
