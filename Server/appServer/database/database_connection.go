package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	MongoDb := os.Getenv("MONGODB_URI")
	if MongoDb == "" {
		log.Fatal("MONGODB_URI NOT SET in .env file")
	}

	fmt.Println("Mongo DB URI:", MongoDb)

	clientOptions := options.Client().ApplyURI(MongoDb)
	Client, err := mongo.Connect(nil, clientOptions)
	if err != nil {
		return nil
	}

	return Client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(collectionName string) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	databaseName := os.Getenv("DATABASE_NAME")
	fmt.Println("Database Name:", databaseName)

	collection := Client.Database(databaseName).Collection(collectionName)

	if collection == nil {
		return nil
	}
	return collection
}
