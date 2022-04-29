package token

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/dao"
	"strconv"
)

var (
	tokenJwt = make(map[string]string)
	jwtToken = make(map[string]string)

	signMethod = jwt.SigningMethodHS512
)

type JwtClaims struct {
	jwt.StandardClaims
	Admin    bool   `json:"admin"`
	Username string `json:"username"`
}

func NewJwt(user *dao.User) string {
	signingString, err := jwt.NewWithClaims(signMethod, &JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Id: strconv.Itoa(int(user.ID)),
		},
		Username: user.Username,
		Admin:    user.IsAdmin,
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
	tokenJwt[token] = newJwt
	return token
}

func md5Hash(content string) string {
	b := md5.Sum([]byte(content))
	token := md5.Sum([]byte(string(b[:]) + "hellasdgasjdYGSDsayufas"))
	return fmt.Sprintf("%x", token)
}

func FromJwt(jwt string) string {
	if s, found := jwtToken[jwt]; found {
		return s
	} else {
		token := md5Hash(jwt)
		jwtToken[jwt] = token
		return token
	}
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
	return ClaimsFromJwt(tokenJwt[token])
}
