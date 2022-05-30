package test

import (
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/server"
	"gopkg.in/natefinch/lumberjack.v2"
	"math/rand"
)

func InitTest() {
	config.Config = config.RootConfig{
		Debug: true,
		Log: config.Log{
			Level:   "debug",
			Console: true,
			Logger:  &lumberjack.Logger{},
		},
		Server: config.Server{
			Port: uint(rand.Int31n(10000) + 10000),
		},
		Database: config.Database{
			Sqlite: &config.Sqlite{
				Filename: "./test.db",
			},
		},
		Token: config.Token{
			Storage: config.TokenStorage{
				Redis: false,
				InMemory: &config.InMemory{
					File: "tokens-test.json",
				},
			},
		},
		Redis: nil,
	}
	server.Init()
}
