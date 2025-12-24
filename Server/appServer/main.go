package main

import (
	"fmt"

	"github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/database"
	"github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {

	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "hello, magicstream movies")
	})

	var client *mongo.Client = database.Connect()

	routes.SetupUnProtectedRoutes(router, client)
	routes.SetupProtectedRoutes(router, client)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server:", err)
	}

}
