package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api/resp"
	"github.com/se2022-qiaqia/course-system/api/token"
	"net/http"
	"strings"
)

const bearerString = "Bearer"
const ClaimsKey = "claims"

func GetClaims(c *gin.Context) *token.JwtClaims {
	if claims, exists := c.Get(ClaimsKey); exists {
		return claims.(*token.JwtClaims)
	}
	return nil
}

func Authorize(context *gin.Context) {
	authorization := context.GetHeader("Authorization")
	bearer := strings.SplitN(authorization, " ", 2)
	if len(bearer) >= 2 && bearer[0] == bearerString {
		t := bearer[1]
		claims := token.ToClaims(t)
		context.Set(ClaimsKey, claims)
	}
	context.Next()
}

func AuthorizedRequired(c *gin.Context) {
	if GetClaims(c) == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, resp.Response{Msg: "请先登录"})
		return
	}
	c.Next()
}
