package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Collection *mongo.Collection

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(os.Getenv("uri"))
		fmt.Println(os.Getenv("dbname"))
		fmt.Println(os.Getenv("colname"))
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(os.Getenv("uri"))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Printf("An error occured while connecting to mongodb. Error(%v)", err)
	}
	fmt.Println("MongoDB connection successful")

	Collection = client.Database(os.Getenv("dbname")).Collection(os.Getenv("colname"))
	fmt.Println("Collection instance is ready.")
}
