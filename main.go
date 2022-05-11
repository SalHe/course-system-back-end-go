package main

import (
	"fmt"
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/log"
	"github.com/se2022-qiaqia/course-system/router"
)

func main() {
	config.Init()
	dao.Init()
	dao.Migrate()

	// TODO 规范响应体
	// TODO 响应体脱敏

	_ = router.NewRouter().Run(fmt.Sprintf(":%d", config.Config.Server.Port))

	log.Logger.Printf("Server is running on port %d\n", config.Config.Server.Port)
}
