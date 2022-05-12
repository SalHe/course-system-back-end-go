package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorResponse struct {
	Msg  string  `json:"msg"`
	Code ErrCode `json:"code"`
}

type ErrCode int

const (
	_ = iota
	ErrCodeFail
	ErrCodeNotFound
	ErrCodeConflict
	ErrCodeUnauthorized
	ErrCodeInternal
)

var errAndHttpCode = map[ErrCode]int{
	ErrCodeFail:         http.StatusBadRequest,
	ErrCodeNotFound:     http.StatusNotFound,
	ErrCodeConflict:     http.StatusConflict,
	ErrCodeUnauthorized: http.StatusUnauthorized,
	ErrCodeInternal:     http.StatusInternalServerError,
}

func Ok(data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, data)
}

func FailJust(msg string, c *gin.Context) {
	Fail(http.StatusBadRequest, msg, c)
}

func Fail(errCode ErrCode, msg string, c *gin.Context) {
	var status int
	var find bool
	if status, find = errAndHttpCode[errCode]; !find {
		status = http.StatusBadRequest
	}
	FailWithHttpStatus(status, errCode, msg, c)
}

func FailWithHttpStatus(httpCode int, errCode ErrCode, msg string, c *gin.Context) {
	c.AbortWithStatusJSON(httpCode, ErrorResponse{
		Msg:  msg,
		Code: errCode,
	})
}
