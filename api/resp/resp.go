package resp

import "github.com/se2022-qiaqia/course-system/dao"

type Response struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type User struct {
	Id           uint     `json:"id"`
	Username     string   `json:"username"`
	RealName     string   `json:"realName"`
	Role         dao.Role `json:"role"`
	College      College  `json:"college"`
	EntranceYear uint     `json:"entranceYear"`
}

type College struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
