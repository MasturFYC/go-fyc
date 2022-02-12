package middleware

import (
	"encoding/json"

	"go-fyc/models"

	"net/http"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	db, _ := CreateConnection()

	defer db.Close()

	var products []models.Product
	db.Find(&products)
	json.NewEncoder(w).Encode(&products)
}
