package config

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
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

type InMemory struct {
	File string `yaml:"file"`
}

type TokenStorage struct {
	InMemory *InMemory `yaml:"in-memory"`
}

type Token struct {
	SignKey string       `yaml:"sign-key"`
	Storage TokenStorage `yaml:"storage"`
}

type Log struct {
	Level   string             `yaml:"level"`
	Console bool               `yaml:"console"`
	Logger  *lumberjack.Logger `yaml:"logger"`
}

type RootConfig struct {
	Debug    bool     `yaml:"debug"`
	Log      Log      `yaml:"log"`
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Token    Token    `yaml:"token"`
}

func Init() {
	var bytes []byte
	var err error
	if bytes, err = ioutil.ReadFile("./config.yml"); err != nil {
		log.Error().Err(err).Msg("读取配置文件失败!")
		return
	}
	if err = yaml.Unmarshal(bytes, &Config); err != nil {
		log.Error().Err(err).Msg("解析配置文件失败!")
		return
	}
}
