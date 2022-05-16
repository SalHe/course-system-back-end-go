package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/model/resp"
	"github.com/se2022-qiaqia/course-system/token"
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
		resp.Fail(resp.ErrCodeUnauthorized, "请先登录", c)
		return
	}
	c.Next()
}

func AuthorizedRoleRequired(roles ...dao.Role) func(c *gin.Context) {
	return func(c *gin.Context) {
		if claims := GetClaims(c); claims != nil {
			for _, r := range roles {
				if r == claims.Role {
					c.Next()
					return
				}
			}
		}
		resp.Fail(resp.ErrCodeUnauthorized, "您的权限不足", c)
	}
}
