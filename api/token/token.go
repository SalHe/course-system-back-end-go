package token

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/golang-jwt/jwt"
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/dao"
)

var (
	tokenJwt = make(map[string]string)
	jwtToken = make(map[string]string)
	mux      = &sync.Mutex{}

	signMethod = jwt.SigningMethodHS512
)

type JwtClaims struct {
	jwt.StandardClaims
	Role     uint   `json:"role"`
	Username string `json:"username"`
}

func NewJwt(user *dao.User) string {
	signingString, err := jwt.NewWithClaims(signMethod, &JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Id: strconv.Itoa(int(user.ID)),
		},
		Username: user.Username,
		Role:     user.Role,
	}).SignedString([]byte(config.Config.Token.SignKey))
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidKey) {
			panic("请检查您的sign-key是否有效")
		}
		return ""
	}
	return signingString
}

func saveToken(token string, jwt string) {
	mux.Lock()
	defer mux.Unlock()

	tokenJwt[token] = jwt
}

func getJwt(token string) (string, bool) {
	mux.Lock()
	defer mux.Unlock()

	jwt2, ok := tokenJwt[token]
	return jwt2, ok
}

func NewToken(user *dao.User) string {
	newJwt := NewJwt(user)
	token := FromJwt(newJwt)
	saveToken(token, newJwt)
	return token
}

func md5Hash(content string) string {
	b := md5.Sum([]byte(content))
	token := md5.Sum([]byte(string(b[:]) + "hellasdgasjdYGSDsayufas"))
	return fmt.Sprintf("%x", token)
}

func saveJwt(jwt string, token string) {
	mux.Lock()
	defer mux.Unlock()
	jwtToken[jwt] = token
}

func getToken(jwt string) (string, bool) {
	mux.Lock()
	defer mux.Unlock()

	s, b := jwtToken[jwt]
	return s, b
}

func FromJwt(jwt string) string {
	if s, found := getToken(jwt); found {
		return s
	} else {
		token := md5Hash(jwt)
		saveJwt(jwt, token)
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
	jwt2, _ := getJwt(token)
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
