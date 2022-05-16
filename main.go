package main

import (
	"context"
	"fmt"
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/dao"
	docs "github.com/se2022-qiaqia/course-system/docs/swagger"
	"github.com/se2022-qiaqia/course-system/log"
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
		log.Logger.Fatalln(err)
	}

	log.Logger.Println("正在关闭服务器...")
	select {
	case <-ctx.Done():
		log.Logger.Println("已关闭服务器.")
	}
	log.Logger.Println("See you.")

}

func runServer(server *http.Server) {
	go func() {
		// _ = r.Run(fmt.Sprintf(":%d", config.Config.Server.Port))
		log.Logger.Println("服务将运行于：", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger.Fatalln(err)
		}
	}()
}

func createServer() *http.Server {
	r := router.NewRouter()
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
