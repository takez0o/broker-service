package main

import (
	"fmt"
	"log"
	"net/http"
)

const web_port = "80"

type Config struct{}

func main() {
	app := Config{}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", web_port),
		Handler: app.routes(),
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
