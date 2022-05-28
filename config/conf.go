package config

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

var Config RootConfig

type Server struct {
	Port uint `yaml:"port"`
}

type InMemory struct {
	File string `yaml:"file"`
}

type TokenStorage struct {
	Redis    bool      `yaml:"redis"`
	InMemory *InMemory `yaml:"in-memory"`
}

type Token struct {
	SignKey string       `yaml:"sign-key"`
	Expire  string       `yaml:"expire"`
	Storage TokenStorage `yaml:"storage"`
}

var defaultExpireDuration, _ = time.ParseDuration("30d")

func (t *Token) ExpireDuration() time.Duration {
	duration, err := time.ParseDuration(t.Expire)
	if err != nil {
		return defaultExpireDuration
	}
	return duration
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
	Redis    *redis   `yaml:"redis"`
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
