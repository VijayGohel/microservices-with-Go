package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %s\n", webPort)

	//define http server
	serv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	//start server
	err := serv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
