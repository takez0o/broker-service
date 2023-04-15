package main

import (
	"authentication-service/cmd/api/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const web_port = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting server on port", web_port)

	conn := connectToDB()
	if conn == nil {
		log.Panic("Cant connect to DB")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", web_port),
		Handler: app.routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready")
			counts++
		} else {
			log.Println("Postgres connected")
			return connection
		}
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Retrying in 2 seconds")
		time.Sleep(2 * time.Second)
		continue
	}

}
