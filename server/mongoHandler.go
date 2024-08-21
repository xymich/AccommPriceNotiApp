package server

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)
var mongoPassVar string;

// DB connection
func MongoConnect() (mongoClient *mongo.Client) {
	if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
	mongoPass, exists := os.LookupEnv("MONGOPASS")
	if exists {
		mongoPassVar = mongoPass
		}
	url := fmt.Sprintf("mongodb+srv://user1:%s@cluster0.iwkktm9.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0", mongoPassVar)
	ctx = context.TODO()
	// ? Connect to MongoDB
	mongoconn := options.Client().ApplyURI(url)
	mongoclient, err := mongo.Connect(ctx, mongoconn)
	if err != nil {
		panic(err)
	}
	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")
	return mongoclient
}