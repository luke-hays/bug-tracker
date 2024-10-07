package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type DatabaseContext struct {
	Connection *pgx.Conn
}

func Init() (*DatabaseContext, error) {
	dbUrl := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(context.Background(), dbUrl)

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database successfully.")

	return &DatabaseContext{
		Connection: conn,
	}, nil
}
