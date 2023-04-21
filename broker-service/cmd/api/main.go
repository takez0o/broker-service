package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const web_port = "8080"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	rabbit_conn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbit_conn.Close()
	app := Config{
		Rabbit: rabbit_conn,
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", web_port),
		Handler: app.routes(),
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("Failed to connect to RabbitMQ")
			counts++
		} else {
			fmt.Println("Connected to RabbitMQ")
			connection = c
			break
		}
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}
		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backoff)
		continue
	}
	return connection, nil
}
