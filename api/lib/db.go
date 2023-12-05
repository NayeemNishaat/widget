package lib

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() *pgxpool.Pool {
	conn := connectToDB()

	if conn == nil {
		log.Panic("Can't connect to the DB.")
	}

	return conn
}

func connectToDB() *pgxpool.Pool {
	count := 0

	dsn := os.Getenv("DSN")

	for {
		conn, err := openDB(dsn)

		if err != nil {
			log.Println("DB is not yet ready.")
		} else {
			log.Print("Connected to DB.")
			return conn
		}

		if count > 5 {
			return nil
		}

		log.Print("Backing off for 1 second.")
		time.Sleep(1 * time.Second)

		count++
		continue
	}
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())

	if err != nil {
		return nil, err
	}

	return db, nil
}
