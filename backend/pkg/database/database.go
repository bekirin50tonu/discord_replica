package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Initialize(config DatabaseConfig) (*mongo.Client, *mongo.Database) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Url))

	if err != nil {
		log.Fatalf("Connection not Success.Err:%e", err)
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatalf("Ping Not Successed.Err:%v", err)
	}
	db := client.Database(config.DatabaseName)
	return client, db
}
