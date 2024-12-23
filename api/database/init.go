package database

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	// "api/models"
)

func InitializeTables(db *bun.DB, ctx context.Context) error {
	// Create tables in the correct order with proper constraints
	tables := []string{
		`CREATE TABLE IF NOT EXISTS categories (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS suppliers (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR NOT NULL UNIQUE,
			email VARCHAR UNIQUE,
			phone VARCHAR
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR NOT NULL,
			category_id UUID NOT NULL REFERENCES categories(id),
			price FLOAT NOT NULL,
			quantity INTEGER NOT NULL,
			image_url VARCHAR,
			supplier_id UUID NOT NULL REFERENCES suppliers(id)
		)`,
	}

	for _, query := range tables {
		_, err := db.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
} 