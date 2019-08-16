package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	// 静态资源加载，本例为css,js以及资源图片
	router.StaticFS("/", http.Dir("./apps/dist"))
	//router.StaticFile("/index.html", "./apps/dist/index.html")
	// Listen and serve on 0.0.0.0:80
	router.Run(":8080")

}
