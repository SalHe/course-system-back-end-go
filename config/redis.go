package config

import "strconv"

type redis struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Db       int    `json:"db"`
	Prefix   string `json:"prefix"`
}

func (r *redis) Addr() string {
	return r.Host + ":" + strconv.Itoa(r.Port)
}
