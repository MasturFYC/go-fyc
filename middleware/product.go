package middleware

import (
	"encoding/json"
	"log"

	"go-fyc/models"

	"net/http"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	products, err := getAllProducts()

	if err != nil {
		log.Fatalf("Unable to get all products. %v", err)
	}

	json.NewEncoder(w).Encode(&products)
}

func getAllProducts() ([]models.Product, error) {
	// defer Sql().Close()
	var products []models.Product

	sqlStatement := `SELECT
		id, name, spec, base_unit,
		base_weight, base_price, first_stock,
		stock, is_active, is_sale, category_id
	FROM products 
	ORDER BY name`

	rows, err := Sql().Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute product query %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Spec,
			&product.BaseUnit,
			&product.BaseWeight,
			&product.BasePrice,
			&product.FirstStock,
			&product.Stock,
			&product.IsActive,
			&product.IsSale,
			&product.CategoryID)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		products = append(products, product)
	}

	return products, err
}
