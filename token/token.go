package token

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/dao"
)

var (
	signMethod = jwt.SigningMethodHS512
)

var Storage TokenStorage

type JwtClaims struct {
	jwt.StandardClaims
	*dao.User
}

func Init() {
	tokenConf := config.Config.Token
	if tokenConf.Storage.InMemory != nil {
		Storage = NewInMemoryTokenStorage()
		Storage.(*InMemoryTokenStorage).Load(tokenConf.Storage.InMemory.File, false)
	}
}

func WhenExit() {
	if ts, ok := Storage.(*InMemoryTokenStorage); ok {
		ts.Save(config.Config.Token.Storage.InMemory.File)
	}
}

func NewJwt(user *dao.User) string {
	signingString, err := jwt.NewWithClaims(signMethod, &JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Id: strconv.Itoa(int(user.ID)),
		},
		User: user,
	}).SignedString([]byte(config.Config.Token.SignKey))
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidKey) {
			panic("请检查您的sign-key是否有效")
		}
		return ""
	}
	return signingString
}

func NewToken(user *dao.User) string {
	newJwt := NewJwt(user)
	token := FromJwt(newJwt)
	Storage.Set(token, newJwt)
	return token
}

func md5Hash(content string) string {
	b := md5.Sum([]byte(content))
	token := md5.Sum([]byte(string(b[:]) + "hellasdgasjdYGSDsayufas"))
	return fmt.Sprintf("%x", token)
}

func FromJwt(jwt string) string {
	token := md5Hash(jwt)
	return token
}

func ClaimsFromJwt(jwtString string) *JwtClaims {
	claims, err := jwt.ParseWithClaims(jwtString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != signMethod {
			return "", fmt.Errorf("该算法不正确：%v", token.Method)
		}
		return []byte(config.Config.Token.SignKey), nil
	})
	if err != nil {
		return nil
	}
	return claims.Claims.(*JwtClaims)
}

func ToClaims(token string) *JwtClaims {
	jwt2, _ := Storage.Get(token)
	return ClaimsFromJwt(jwt2)
}

func (c *JwtClaims) IsAdmin() bool {
	return c.Role == dao.RoleAdmin
}

func (c *JwtClaims) IsTeacher() bool {
	return c.Role == dao.RoleTeacher
}

func (c *JwtClaims) IsStudent() bool {
	return c.Role == dao.RoleStudent
}
