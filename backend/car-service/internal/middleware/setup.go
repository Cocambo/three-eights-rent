package middleware

import "github.com/gin-gonic/gin"

func Setup(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
}
