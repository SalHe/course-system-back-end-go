package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/cache"
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/log"
	"github.com/se2022-qiaqia/course-system/middleware"
	"github.com/se2022-qiaqia/course-system/router"
	"github.com/se2022-qiaqia/course-system/token"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

func Init() {
	log.Init()
	token.Init()
	dao.Init()
	dao.Migrate()

	// Redis 目前是非必须的
	if config.Config.Redis != nil {
		cache.Init()
	}
}

func RunServer(server *http.Server) {
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

func CreateServer() *http.Server {
	baseEngine := gin.New()

	logger := log.L
	baseEngine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	baseEngine.Use(middleware.LoggerWithZerolog(logger))
	baseEngine.Use(middleware.RecoveryWithZerolog(logger, true, config.Config.Debug))

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
