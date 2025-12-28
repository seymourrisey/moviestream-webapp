package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/database"
	"github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	log.Println("Starting MagicStream API...")

	// Load env (local only, Railway injects env automatically)
	_ = godotenv.Load()

	router := gin.Default()

	// Health check / sanity endpoint
	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, MagicStreamMovies!")
	})

	// ---- CORS ----
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	var origins []string

	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
			log.Println("Allowed Origin:", origins[i])
		}
	} else {
		// SAFE DEFAULT for production demo
		origins = []string{"*"}
		log.Println("Allowed Origin: *")
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ---- MongoDB ----
	var mongoClient *mongo.Client
	var mongoAvailable bool = true

	log.Println("Connecting to MongoDB...")
	mongoClient = database.Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Println("WARNING: MongoDB unavailable, running without DB:", err)
		mongoAvailable = false
	}

	if mongoAvailable {
		defer func() {
			if err := mongoClient.Disconnect(context.Background()); err != nil {
				log.Println("Mongo disconnect error:", err)
			}
		}()
	}

	// ---- Routes ----
	routes.SetupUnProtectedRoutes(router, mongoClient)

	if mongoAvailable {
		routes.SetupProtectedRoutes(router, mongoClient)
	} else {
		log.Println("Protected routes disabled (MongoDB not connected)")
	}

	// ---- PORT (Railway REQUIRED) ----
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server listening on port", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
