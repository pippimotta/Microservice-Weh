package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Mailer Mail
}

const webPort = "80"

func main() {
	app := Config{
		Mailer: newMail(),
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func newMail() Mail {
	port , _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail {
		Domain: os.Getenv("MAIL_DOMAIN") ,
		Host: os.Getenv("HOST") ,
		Port: port,
		Username: os.Getenv("USER_NAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
		Encryption:os.Getenv("MAIL_ENCRYPTION") ,
		FromName: os.Getenv("MAIL_FROMNAME"),
		FromAddress: os.Getenv("MAIL_FROMADDRESS"),
	}
	return m
} 
