package database

import (
	"context"
	"fmt"

	"devops.zedeks.com/TheHiddenDeveloper/ims-zedeks/api/models"
	"github.com/uptrace/bun"
	// "api/models"
)

func InitializeTables(db *bun.DB, ctx context.Context) error {
	// Create tables in the correct order
	models := []interface{}{
		(*models.Category)(nil),
		(*models.Supplier)(nil),
		(*models.Products)(nil),
		// Add other models here
	}

	for _, model := range models {
		_, err := db.NewCreateTable().
			Model(model).
			IfNotExists().
			WithForeignKeys().
			Exec(ctx)
		
		if err != nil {
			return fmt.Errorf("failed to create table for %T: %w", model, err)
		}
	}
	
	return nil
} 