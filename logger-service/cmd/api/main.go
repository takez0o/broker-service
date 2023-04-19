package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"logger-service/data"
)

const (
	web_port  = "80"
	rpc_port  = "5001"
	mongo_url = "mongodb://mongo:27017"
	grpc_port = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	log.Println("Starting server on port: ", web_port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", web_port),
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server: ", err)
		log.Panic()
	}
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongo_url)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	c, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println("Error connecting to mongo: ", err)
		return nil, err
	}
	return c, nil
}
