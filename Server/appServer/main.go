package main

import (
	"fmt"

	controller "github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/controllers"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "hello, magicstream movies")
	})

	router.GET("/movies", controller.GetMovies())
	router.GET("/movie/:imdb_id", controller.GetMovie())
	router.POST("/addmovie", controller.AddMovie())
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server:", err)
	}

}
