package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseContext struct {
	Connection *pgxpool.Pool
}

func Init() (*DatabaseContext, error) {
	dbUrl := os.Getenv("DATABASE_URL")

	conn, err := pgxpool.New(context.Background(), dbUrl)

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database successfully.")

	return &DatabaseContext{
		Connection: conn,
	}, nil
}
