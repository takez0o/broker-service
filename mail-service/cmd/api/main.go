package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Mailer Mail
}

const web_port = "80"

func main() {
	app := Config{
		Mailer: createMail(),
	}
	log.Println("Starting mail service on port ,", web_port)

	srv := &http.Server{
		Addr:    ":" + web_port,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	return Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("MAIL_FROM_NAME"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
	}
}
