package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-fyc/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	categories, err := getAllCategories()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	json.NewEncoder(w).Encode(&categories)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	cat, err := getCategory(&id)

	if err != nil {
		log.Fatalf("Unable to get category. %v", err)
	}

	json.NewEncoder(w).Encode(&cat)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	deletedRows := deleteCategory(&id)

	msg := fmt.Sprintf("Category deleted successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := Response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)

}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var cat models.Category

	err := json.NewDecoder(r.Body).Decode(&cat)

	id := createCategory(cat.Name)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	cat.ID = id
	json.NewEncoder(w).Encode(&cat)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {

	EnableCors(&w)

	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	// create the postgres db connection

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	var cat models.Category

	err := json.NewDecoder(r.Body).Decode(&cat)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	updatedRows := updateCategory(&id, &cat)

	msg := fmt.Sprintf("Category updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := Response{
		ID:      updatedRows,
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)

}

func getAllCategories() ([]models.Category, error) {
	// defer Sql.Close()
	var categories []models.Category

	sqlStatement := `SELECT id, name FROM categories ORDER BY name`

	rows, err := Sql.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute category query %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var cat models.Category
		err := rows.Scan(&cat.ID, &cat.Name)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		categories = append(categories, cat)
	}

	return categories, err
}

func getCategory(id *int) (models.Category, error) {
	var cat models.Category

	sqlStatement := `SELECT c.id, c.name FROM categories c WHERE c.id=$1`
	//stmt, _ := Sql.Prepare(sqlStatement)

	//defer stmt.Close()
	row := Sql.QueryRow(sqlStatement, id)

	err := row.Scan(&cat.ID, &cat.Name)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return cat, nil
	case nil:
		products, _ := getProductsByCategory(id)
		cat.Products = products
		return cat, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty user on error
	return cat, err
}

func getProductsByCategory(id *int) ([]models.Product, error) {
	// defer Sql.Close()

	var products []models.Product

	sqlStatement := `SELECT
		p.id, p.name, p.spec, p.base_unit,
		p.base_weight, p.base_price, p.first_stock,
		p.stock, p.is_active, p.is_sale, p.category_id
	FROM products AS p
	WHERE p.category_id=$1
	ORDER BY p.name`

	rows, err := Sql.Query(sqlStatement, id)

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
			&product.CategoryID,
		)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		products = append(products, product)
	}

	return products, err
}

func deleteCategory(id *int) int64 {
	// create the delete sql query
	sqlStatement := `DELETE FROM categories WHERE id=$1`

	// execute the sql statement
	res, err := Sql.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to delete category. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	return rowsAffected
}

func createCategory(catName string) int {

	sqlStatement := `INSERT INTO categories (name) VALUES ($1) RETURNING id`

	var id int

	err := Sql.QueryRow(sqlStatement, catName).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to create category. %v", err)
	}

	return id
}

func updateCategory(id *int, cat *models.Category) int64 {

	sqlStatement := `UPDATE categories SET name=$2 WHERE id=$1`

	res, err := Sql.Exec(sqlStatement, id, cat.Name)

	if err != nil {
		log.Fatalf("Unable to update category. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while updating category. %v", err)
	}

	return rowsAffected
}
