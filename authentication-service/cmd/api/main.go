package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const webPort = "80"

var dbRetriesCount int32

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Printf("Starting authentication service on port %s\n", webPort)
	
	//connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Fatal("Could not connect to Postgres")
	}
	defer conn.Close()

	app := Config{
		DB:     conn,
        Models: data.New(conn),
	}

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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
        return nil, err
    }

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		db, err := openDB(dsn)
        if err != nil {
			log.Println("Postgres not yet ready...")
            dbRetriesCount++
        } else {
			log.Printf("Connected to Postgres after %d retries\n", dbRetriesCount)
			return db
		}

		if dbRetriesCount >= 10 {
            log.Fatal("Unable to connect to Postgres after 10 retries", err)
			return nil;
        }
		
        log.Println("Retrying connection to Postgres...")
		time.Sleep(time.Second * 2)
		continue
	}
}
