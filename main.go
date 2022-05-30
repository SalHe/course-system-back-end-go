package main

import (
	"context"
	"fmt"
	"github.com/se2022-qiaqia/course-system/config"
	docs "github.com/se2022-qiaqia/course-system/docs/swagger"
	"github.com/se2022-qiaqia/course-system/flags"
	"github.com/se2022-qiaqia/course-system/log"
	"github.com/se2022-qiaqia/course-system/server"
	"github.com/se2022-qiaqia/course-system/token"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

//go:generate swag init -o docs/swagger

// @title 								选课系统
// @version 							1.0
// @description 						软件工程课设——选课系统.
// @termsOfService 						http://swagger.io/terms/

// @contact.name 						Qiaqia
// @contact.url 						https://github.com/se2022-qiaqia

// @BasePath 							/api/v1

// @securityDefinitions.apikey 			ApiKeyAuth
// @in 									header
// @name 								Authorization
// ----@description 						Bearer <Token>
// 我也不知道这里搞什么玩意儿他不给我生成认证的描述，很是无语

func main() {
	flags.Parse()
	config.Init(*flags.ConfigPath)
	abs, _ := filepath.Abs(*flags.ConfigPath)
	fmt.Printf("启用日志文件：%s\n", abs)

	server.Init()
	defer token.WhenExit()

	// 引用一下，不然不会生成swagger文档
	docs.SwaggerInfo.InstanceName()

	svr := server.CreateServer()
	server.RunServer(svr)

	// 设置可控退出逻辑

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		log.L.Error().Err(err).Msg("关闭服务器出错")
	}

	log.L.Info().Msg("正在关闭服务器...")
	select {
	case <-ctx.Done():
		log.L.Info().Msg("已关闭服务器")
	}
	log.L.Info().Msg("See you.")

}
