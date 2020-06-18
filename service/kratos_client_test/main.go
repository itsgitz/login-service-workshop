package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello!",
		})
	})

	r.GET("/registration", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "http://127.0.0.1:4433/self-service/browser/flows/registration")
	})

	r.GET("/auth/registration", func(c *gin.Context) {
		log.Println(c)
	})

	r.Run(":9080")
}
