package main

import (
	"fmt"
	"github.com/se2022-qiaqia/course-system/api"
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/log"
)

func main() {
	config.Init()
	dao.Init()
	dao.Migrate()

	_ = api.NewRouter().Run(fmt.Sprintf(":%d", config.Config.Server.Port))

	log.Logger.Printf("Server is running on port %d\n", config.Config.Server.Port)
}
