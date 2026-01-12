package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

const DataBaseUrl string = "postgresql://postgres:root@localhost:5432/subscribly"
var DBConn *pgxpool.Pool

func makeDBPool() (*pgxpool.Pool ,error) {
	pool, err := pgxpool.New(context.TODO(), DataBaseUrl)
	if err != nil {
		return pool, err
	}
	return pool, nil
}

func InitializePool() {
	pool, err := makeDBPool()
	if err != nil {
		log.Fatalf("failed to initialize pool : %v", err)
	}

	DBConn = pool
}

func CloseDBPool(pool *pgxpool.Pool) {
	pool.Close()
}

