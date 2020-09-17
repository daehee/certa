package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddDomain(d string) {
	client, ctx, _ := mongoConnect()

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"domain", d}}
	update := bson.D{{"$set", bson.D{
		{Key: "domain", Value: d},
	}}}
	_, err := client.Database("recon").Collection("domains").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		sugar.Fatal(err)
	}

	client.Disconnect(ctx)
}

func mongoConnect() (*mongo.Client, context.Context, context.CancelFunc) {
	mongoEndpoint := os.Getenv("MONGODB_ENDPOINT") // export MONGODB_ENDPOINT=localhost:27017
	if mongoEndpoint == "" {
		log.Fatal("MONGODB_ENDPOINT not set as env variable")
		os.Exit(1)
	}
	mongoConnect := fmt.Sprintf("mongodb://%s", mongoEndpoint)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoConnect))
	if err != nil {
		sugar.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		sugar.Fatal(err)
	}
	// Force a connection to verify our connection string
	// err = client.Ping(ctx, nil)
	// if err != nil {
	// 	sugar.Fatalf("Failed to ping cluster: %v", err)
	// }

	return client, ctx, cancel
}
