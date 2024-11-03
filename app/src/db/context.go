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
	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database successfully.")

	return &DatabaseContext{
		Connection: conn,
	}, nil
}
