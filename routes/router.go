package routes

import (
	v1 "ginblog/api/v1"
	"ginblog/middleware"
	"ginblog/utils"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRoutor() {
	gin.SetMode(utils.AppMode)
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	// 添加prometheus中间件
	gp := middleware.New(r)
	r.Use(gp.Middleware())

	// 添加prometheus metrics采样
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	auth := r.Group("api/v1")
	auth.Use(middleware.JwtToken())
	{
		//用户模块的路由接口
		auth.PUT("user/:id", v1.EditUser)
		auth.DELETE("user/:id", v1.DeleteUser)

		//分类模块的路由接口
		auth.POST("cate/add", v1.AddCategory)
		auth.PUT("cate/:id", v1.EditCategory)
		auth.DELETE("cate/:id", v1.DeleteCategory)
		//文章模块的路由接口
		auth.POST("art/add", v1.AddArticle)
		auth.PUT("art/:id", v1.EditArticle)
		auth.DELETE("art/:id", v1.DeleteArticle)
		//上传文件
		auth.POST("upload", v1.UpLoad)
	}

	router := r.Group("api/v1")
	{
		router.GET("users", v1.GetUser)
		router.GET("cate", v1.GetCategory)
		router.GET("art", v1.GetArticle)
		router.GET("art/list/:id", v1.GetCateArt)
		router.GET("art/info/:id", v1.GetArtInfo)
		router.POST("login", v1.Login)
		router.POST("user/add", v1.AddUser)
	}

	// gp := ginprometheus.New(r)
	// r.Use(gp.Middleware())
	// // metrics采样
	// r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Run(utils.HttpPort)
}
