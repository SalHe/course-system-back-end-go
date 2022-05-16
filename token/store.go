package token

type TokenStorage interface {
	Get(token string) (string, bool)      // 根据token获取对应的jwt
	Set(token string, jwt string)         // 设置token和jwt
	SetExpire(token string, expire int64) // 设置token过期时间
	Delete(token string)                  // 删除token
}
