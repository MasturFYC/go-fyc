package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"go-fyc/models"

	"net/http"

	"github.com/gorilla/mux"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	products, err := getAllProducts()

	if err != nil {
		log.Fatalf("Unable to get all products. %v", err)
	}

	json.NewEncoder(w).Encode(&products)
}

func GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	params := mux.Vars(r)
	var err error

	id, err := strconv.Atoi(params["id"])

	fmt.Printf("test: %v\n", id)

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	var prod []models.Product

	if id == 0 {
		prod, err = getAllProducts()
	} else {
		prod, err = getProductsByCategory(&id)
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		//		json.NewEncoder(w).Encode(http.StatusText(http.StatusNotFound))
		//		return
	}

	json.NewEncoder(w).Encode(&prod)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	params := mux.Vars(r)

	id, err := strconv.ParseInt(params["id"], 10, 64)

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	prod, err := getProduct(&id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		//		json.NewEncoder(w).Encode(http.StatusText(http.StatusNotFound))
		//		return
	}

	json.NewEncoder(w).Encode(&prod)

}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)

	id, err := strconv.ParseInt(params["id"], 10, 64)

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	deletedRows := deleteProduct(&id)

	msg := fmt.Sprintf("Product deleted successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := Response{
		ID:      id,
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)

}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var p models.Product

	err := json.NewDecoder(r.Body).Decode(&p)

	id := createProduct(&p)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	p.ID = id
	json.NewEncoder(w).Encode(&p)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	EnableCors(&w)

	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	// create the postgres db connection

	params := mux.Vars(r)

	id, _ := strconv.ParseInt(params["id"], 10, 64)

	var p models.Product

	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	updatedRows := updateProduct(&id, &p)

	msg := fmt.Sprintf("Product updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := Response{
		ID:      updatedRows,
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)

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

func getProduct(id *int64) (models.Product, error) {
	var p models.Product

	sqlStatement := `SELECT
		id, name, spec, base_unit,
		base_weight, base_price, first_stock,
		stock, is_active, is_sale, category_id
	FROM products 
	WHERE id=$1`

	//defer stmt.Close()
	row := Sql().QueryRow(sqlStatement, id)

	err := row.Scan(&p.ID,
		&p.Name,
		&p.Spec,
		&p.BaseUnit,
		&p.BaseWeight,
		&p.BasePrice,
		&p.FirstStock,
		&p.Stock,
		&p.IsActive,
		&p.IsSale,
		&p.CategoryID)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("Product not found!")
		return p, err
	case nil:
		units, _ := getUnitsByProduct(id)
		p.Units = units
		return p, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty user on error
	return p, err
}

func getUnitsByProduct(id *int64) ([]models.Unit, error) {
	// defer Sql().Close()
	var units []models.Unit

	sqlStatement := `SELECT
		id, name, barcode, content,
		buy_price, margin, price,
		is_default, product_id
	FROM units 
	WHERE product_id=$1`

	rows, err := Sql().Query(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute product query %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var u models.Unit
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Barcode,
			&u.Content,
			&u.BuyPrice,
			&u.Margin,
			&u.Price,
			&u.IsDefault,
			&u.ProductID)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		units = append(units, u)
	}

	return units, err
}

func deleteProduct(id *int64) int64 {
	// create the delete sql query
	sqlStatement := `DELETE FROM products WHERE id=$1`

	// execute the sql statement
	res, err := Sql().Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to delete product. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	return rowsAffected
}

func createProduct(prod *models.Product) int64 {

	sqlStatement := `INSERT INTO products
		(name, spec, base_unit, base_weight, base_price, first_stock, stock, is_active, is_sale, category_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7,$8, $9, $10) RETURNING id`

	var id int64

	err := Sql().QueryRow(sqlStatement,
		prod.Name,
		prod.Spec,
		prod.BaseUnit,
		prod.BaseWeight,
		prod.BasePrice,
		prod.FirstStock,
		prod.Stock,
		prod.IsActive,
		prod.IsSale,
		prod.CategoryID).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to create product. %v", err)
	}

	return id
}

func updateProduct(id *int64, prod *models.Product) int64 {

	sqlStatement := `UPDATE products SET 
		name=$2, spec=$3, base_unit=$4,
		base_weight=$5, base_price=$6, first_stock=$7,
		is_active=$8, is_sale=$9, category_id=$10
	WHERE id=$1`

	res, err := Sql().Exec(sqlStatement,
		id,
		prod.Name,
		prod.Spec,
		prod.BaseUnit,
		prod.BaseWeight,
		prod.BasePrice,
		prod.FirstStock,
		prod.IsActive,
		prod.IsSale,
		prod.CategoryID,
	)

	if err != nil {
		log.Fatalf("Unable to update product. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while updating product. %v", err)
	}

	return rowsAffected
}
