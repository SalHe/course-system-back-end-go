package api

import (
	"github.com/gin-gonic/gin"
	"github.com/se2022-qiaqia/course-system/api/admin"
	"github.com/se2022-qiaqia/course-system/api/colleges"
	"github.com/se2022-qiaqia/course-system/api/courses"
	"github.com/se2022-qiaqia/course-system/api/user"
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/middleware"
)

func NewRouter() *gin.Engine {
	engine := gin.Default()
	gin.SetMode(gin.DebugMode)
	if !config.Config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine.Use(middleware.Authorize)

	r := engine.Group("/api/v1")

	// 不需要认证的API
	{
		r := r.Group("/")
		r.POST("/login", user.Login)
		r.POST("/register", user.Register)

		a := r.Group("/admin")
		a.GET("/init", admin.IsInitialized)
		a.POST("/init", admin.InitSystem)
	}

	// 需要认证的API
	r = r.Group("/")
	r.Use(middleware.AuthorizedRequired)
	{
		u := r.Group("/user")
		u.GET("/info", user.GetUserInfo)

	}

	{
		r := r.Group("/admin")
		r.Use(middleware.AuthorizedRoleRequired(dao.RoleAdmin))
		{
			u := r.Group("/user")
			u.GET("/list/:page/:size", admin.GetUserList)
			u.GET("/:id", admin.GetUser)
			u.POST("/:id", admin.UpdateUser)
			u.DELETE("/:id", admin.DeleteUser)
			u.POST("/new", admin.NewUser)
		}
	}
	{
		r := r.Group("/college")
		r.GET("/list", colleges.ListColleges)

		r = r.Group("/")
		r.Use(middleware.AuthorizedRoleRequired(dao.RoleAdmin))
		r.POST("/new", colleges.NewCollege)
	}
	{
		r := r.Group("/course")
		{
			c := r.Group("/")
			c.POST("/list", courses.GetCourseList)

			c = r.Group("")
			c.Use(middleware.AuthorizedRoleRequired(dao.RoleAdmin))
			c.POST("", courses.NewCourse)
			c.PUT("/:id", courses.UpdateCourse)
			c.POST("/open", courses.OpenCourse)
		}
	}
	return engine
}
