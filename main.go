package main

import (
	"fmt"
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/dao"
	docs "github.com/se2022-qiaqia/course-system/docs/swagger"
	"github.com/se2022-qiaqia/course-system/log"
	"github.com/se2022-qiaqia/course-system/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	dao.Init()
	dao.Migrate()

	// 引用一下，不然不会生成swagger文档
	docs.SwaggerInfo.InstanceName()

	r := router.NewRouter()
	r.GET("/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.PersistAuthorization(true),
		),
	)
	_ = r.Run(fmt.Sprintf(":%d", config.Config.Server.Port))

	log.Logger.Printf("Server is running on port %d\n", config.Config.Server.Port)
}
