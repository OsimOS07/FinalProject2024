package main

import (
	"test/app/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(Cors)
	
	r.POST("/signup", controller.SignUp)
	r.POST("/login",controller.Login)
	r.POST("/add",controller.AddProduct)
	r.GET("/getList",controller.List)
	r.POST("/delete",controller.Delete)
	r.Run(":2024")
}
func Cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://192.168.43.70:5500")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	}

	c.Next()
}

