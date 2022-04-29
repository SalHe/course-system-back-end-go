package config

import (
	"github.com/se2022-qiaqia/course-system/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var Config RootConfig

type Server struct {
	Port uint `yaml:"port"`
}

type Sqlite struct {
	Filename string `yaml:"filename"`
}

type Database struct {
	Sqlite *Sqlite `yaml:"sqlite"`
}

type Token struct {
	SignKey string `yaml:"sign-key"`
}

type RootConfig struct {
	Debug    bool     `yaml:"debug"`
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Token    Token    `yaml:"token"`
}

func Init() {
	var bytes []byte
	var err error
	if bytes, err = ioutil.ReadFile("./config.yml"); err != nil {
		log.Logger.Fatalln("读取配置文件失败!")
		return
	}
	if err = yaml.Unmarshal(bytes, &Config); err != nil {
		log.Logger.Fatalln("读取配置出错!")
		return
	}
}
