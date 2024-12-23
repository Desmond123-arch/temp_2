package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Products struct {
    bun.BaseModel `bun:"table:products"`

    ID         uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
    Name       string    `bun:"name,notnull"`
    CategoryID uuid.UUID `bun:"category_id,type:uuid,notnull"`
    Category   Category  `bun:"rel:belongs-to,join:category_id=id"`
    Price      float64   `bun:"price,notnull"`
    Quantity   int       `bun:"quantity,notnull"`
    ImageURL   string    `bun:"image_url"`
    SupplierID uuid.UUID `bun:"supplier_id,type:uuid,notnull"`
    Supplier   Supplier  `bun:"rel:belongs-to,join:supplier_id=id"`
}

type Category struct {
    bun.BaseModel `bun:"table:categories"`

    ID   uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
    Name string    `bun:"name,notnull,unique"`
}

type Supplier struct {
    bun.BaseModel `bun:"table:suppliers"`

    ID    uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
    Name  string    `bun:"name,notnull,unique"`
    Email string    `bun:"email,unique"`
    Phone string    `bun:"phone"`
}
type Status string

func (s *Status) Scan(value interface{}) error {
	*s = Status(fmt.Sprintf("%s", value))
	return nil
}

func (s Status) Value() (driver.Value, error) {
	return string(s), nil
}

type Orders  struct{
	bun.BaseModel `bun:"table:orders"`

	Id uuid.UUID `bun:",pk,type:uuid,default:gen_random_uuid()"` 
	OrderDate   time.Time `bun:"order_date"`
	Status      Status    `bun:"status,type:order_status"`
	TotalAmount float64   `bun:"total_amount"`                           
}

type OrderItem struct {
	bun.BaseModel `bun:"table:order_items,alias:oi"`

	ID        uuid.UUID `bun:",pk,type:uuid,default:gen_random_uuid()"` // Primary key
	OrderID   uuid.UUID `bun:"order_id,notnull"`                       // Foreign key to Orders
	ProductID uuid.UUID `bun:"product_id,notnull"`                     // Foreign key to Products
	Quantity  int       `bun:"quantity,notnull"`                       // Quantity
	Price     float64   `bun:"price,notnull"`                          // Price
}