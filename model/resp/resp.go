package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Msg  string      `json:"msg"`  // 消息，当发生错误时，该提示信息可以用于向用户展示（也可以用于调试），约定此信息必须是用户友好型。
	Code ErrCode     `json:"code"` // 状态码，约定：当为0时表示操作成功，为负数时表示错误，有对应的错误码；暂时未定义正数状态码。
	Data interface{} `json:"data"` // 实际响应数据。
}

type Page struct {
	PageNo   int         `json:"pageNo"`
	PageSize int         `json:"pageSize"`
	Total    int64       `json:"total"`
	Contents interface{} `json:"contents"`
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
	c.JSON(http.StatusOK, Response{
		Msg:  msg,
		Code: ErrCodeOk,
		Data: data,
	})
}

func FailJust(msg string, c *gin.Context) {
	Fail(ErrCodeFail, msg, c)
}

func Fail(errCode ErrCode, msg string, c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, Response{
		Msg:  msg,
		Code: errCode,
		Data: nil,
	})
}

func OkPage(data interface{}, pageNo int, pageSize int, total int64, c *gin.Context) {
	Resp(Page{
		PageNo:   pageNo,
		PageSize: pageSize,
		Total:    total,
		Contents: data,
	}, "ok", c)
}
