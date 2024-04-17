package goodidea

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/url"
	"os"
)

var (
	DB *pgxpool.Pool
)

func Connect() error {
	uri := makeConnStr()
	dbpool, err := pgxpool.New(context.Background(), uri)
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		return err
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		log.Println("open_connection: failed to ping database")
		return err
	}

	DB = dbpool
	return nil
}

func makeConnStr() string {
	d := os.Getenv("DATABASE_URL")
	if d != "" {
		return d
	}
	u := "fairy"
	if os.Getenv("DB_USER") != "" {
		u = os.Getenv("DB_USER")
	}
	pw := "goodidea"
	if os.Getenv("DB_PASS") != "" {
		pw = url.QueryEscape(os.Getenv("DB_PASS"))
	}
	port := "5555"
	if os.Getenv("DB_PORT") != "" {
		pw = os.Getenv("DB_PORT")
		//TODO check it's all numbers
	}
	h := "localhost"
	if os.Getenv("DB_HOST") != "" {
		h = os.Getenv("DB_HOST")
	}
	n := "tasks"
	if os.Getenv("DB_NAME") != "" {
		n = os.Getenv("DB_NAME")
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", u, pw, h, port, n)
}
