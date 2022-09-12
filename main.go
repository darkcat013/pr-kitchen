package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.POST("/order", func(ctx *gin.Context) {
		ctx.JSON(200, "Ordered")
	})
	router.Run(":8080")
}
