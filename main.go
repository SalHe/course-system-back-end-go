package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/dao"
	docs "github.com/se2022-qiaqia/course-system/docs/swagger"
	"github.com/se2022-qiaqia/course-system/log"
	"github.com/se2022-qiaqia/course-system/middleware"
	"github.com/se2022-qiaqia/course-system/router"
	"github.com/se2022-qiaqia/course-system/token"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
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
	config.Init()
	log.Init()

	token.Init()
	defer token.WhenExit()

	dao.Init()
	dao.Migrate()

	// 引用一下，不然不会生成swagger文档
	docs.SwaggerInfo.InstanceName()

	server := createServer()
	runServer(server)

	// 设置可控退出逻辑

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.L.Error().Err(err).Msg("关闭服务器出错")
	}

	log.L.Info().Msg("正在关闭服务器...")
	select {
	case <-ctx.Done():
		log.L.Info().Msg("已关闭服务器")
	}
	log.L.Info().Msg("See you.")

}

func runServer(server *http.Server) {
	go func() {
		// _ = r.Run(fmt.Sprintf(":%d", config.Config.Server.Port))
		log.L.Info().
			Str("addr", server.Addr).
			Msg("服务将运行")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.L.Fatal().Err(err).Msg("服务器启动失败")
		}
	}()
}

func createServer() *http.Server {
	baseEngine := gin.New()

	logger := log.L
	baseEngine.Use(middleware.LoggerWithZerolog(logger))
	baseEngine.Use(middleware.RecoveryWithZerolog(logger, true))

	r := router.NewRouter(baseEngine)

	r.GET("/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.PersistAuthorization(true),
		),
	)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Config.Server.Port),
		Handler: r,
	}
	return server
}
