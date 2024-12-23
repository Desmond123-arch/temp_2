package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var (
	DB  *bun.DB
	Ctx context.Context
)

func ConnectDb() (*bun.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL not set in environment")
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize tables
	if err := InitializeTables(db, context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	// Initialize global variables
	DB = db
	Ctx = context.Background()

	return db, nil
}

func GetDB() *bun.DB {
	if DB == nil {
		var err error
		DB, err = ConnectDb()
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	}
	return DB
}

func GetContext() context.Context {
	if Ctx == nil {
		Ctx = context.Background()
	}
	return Ctx
}
