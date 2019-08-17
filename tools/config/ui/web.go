package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.StaticFS("/", http.Dir("./apps/dist"))
	router.Run(":8080")

}
