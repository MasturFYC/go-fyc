package models

import "time"

type Category struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Products []Product `json:"products,omitempty"`
}

type Product struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Spec       string  `json:"spec"`
	BaseUnit   string  `json:"baseUnit"`
	BaseWeight float32 `json:"baseWeight"`
	BasePrice  float64 `json:"basePrice"`
	FirstStock float32 `json:"firstStock"`
	Stock      float32 `json:"stock"`
	IsActive   bool    `json:"isActive"`
	IsSale     bool    `json:"isSale"`
	CategoryID int     `json:"categoryId"`
	Units      []Unit  `json:"units,omitempty"`
}

type Unit struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Barcode   string  `json:"barcode"`
	Content   float32 `json:"content"`
	BuyPrice  float64 `json:"buyPrice"`
	Margin    float32 `json:"margin"`
	Price     float64 `json:"price"`
	IsDefault bool    `json:"isDefault"`
	ProductID int64   `json:"productId"`
}

type Customer struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Street string  `json:"street"`
	City   string  `json:"city"`
	Phone  string  `json:"phone"`
	Cell   string  `json:"cell"`
	Zip    string  `json:"zip"`
	Email  string  `json:"email"`
	Orders []Order `json:"orders,omitempty"`
}

type Facturer struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Descriptions string    `json:"descriptions"`
	Instructions string    `json:"instructions"`
	Total        float64   `json:"total"`
	Qty          float32   `json:"qty"`
	Price        float64   `json:"price"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Salesman struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Street string  `json:"street"`
	City   string  `json:"city"`
	Phone  string  `json:"phone"`
	Cell   string  `json:"cell"`
	Zip    string  `json:"zip"`
	Email  string  `json:"email"`
	Orders []Order `json:"orders,omitempty"`
}

type Order struct {
	ID            int64         `json:"id"`
	Total         float64       `json:"total"`
	Cash          float64       `json:"cash"`
	Payment       float64       `json:"payment"`
	RemainPayment float64       `json:"remainPayment"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	CustomerID    int           `json:"customerId"`
	SalesID       int           `json:"salesId"`
	OrderDetails  []OrderDetail `json:"details,omitempty"`
}

type OrderDetail struct {
	ID        int64   `json:"id"`
	Qty       float32 `json:"qty"`
	Content   float32 `json:"content"`
	UnitName  string  `json:"unitName"`
	RealQty   float32 `json:"realQty"`
	Price     float64 `json:"price"`
	BuyPrice  float64 `json:"buyPrice"`
	Discount  float64 `json:"discount"`
	Subtotal  float64 `json:"subtotal"`
	OrderID   int64   `json:"orderId"`
	ProductID int64   `json:"productId"`
	UnitID    int64   `json:"unitId"`
}
