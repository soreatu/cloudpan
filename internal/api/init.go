package api

import (
	"cloudpan/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	// 跨域
	r.Use(middleware.Cors())

	// 路由
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", Ping)

		// 用户登录
		v1.POST("/login", Login)
		// 用户注册
		v1.POST("/register", Register)

		user := v1.Group("/user")
		user.Use(middleware.AuthRequired())
		{
			// 用户注销
			user.GET("/logout", Logout)

			// 上传文件
			user.POST("/upload", Upload)
			// 下载文件
			user.GET("/download/:id", Download)

			// 获取所有文件
			user.GET("/files", GetFiles)
			// 删除文件
			user.DELETE("/file/:id", DeleteFile)
		}
	}

	return r
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "pong"})
}
