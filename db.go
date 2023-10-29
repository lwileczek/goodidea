package goodidea

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

var (
	DB *pgxpool.Pool
)

func Connect() error {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
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
