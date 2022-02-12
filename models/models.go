package models

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
}
