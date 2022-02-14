package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"go-fyc/models"

	"net/http"

	"strconv"

	"github.com/gorilla/mux"
)

func GetSalesmans(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	sales, err := getAllSales()

	if err != nil {
		log.Fatalf("Unable to get all sales. %v", err)
	}

	json.NewEncoder(w).Encode(&sales)
}

func GetSalesman(w http.ResponseWriter, r *http.Request) {

	EnableCors(&w)

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	sales, err := getSales(&id)

	if err != nil {
		log.Fatalf("Unable to get category. %v", err)
	}

	json.NewEncoder(w).Encode(&sales)
}

func DeleteSalesman(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	deletedRows := deleteSales(&id)

	msg := fmt.Sprintf("Sales deleted successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := Response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func CreateSalesman(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var sales models.Salesman

	err := json.NewDecoder(r.Body).Decode(&sales)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	id, err := createSales(&sales)

	if err != nil {
		log.Fatalf("Nama sales tidak boleh sama.  %v", err)
	}

	sales.ID = id

	json.NewEncoder(w).Encode(&sales)

}

func UpdateSalesman(w http.ResponseWriter, r *http.Request) {

	EnableCors(&w)

	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	// create the postgres db connection

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	var sales models.Salesman

	err := json.NewDecoder(r.Body).Decode(&sales)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	updatedRows := updateSales(&id, &sales)

	msg := fmt.Sprintf("Sales updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := Response{
		ID:      updatedRows,
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func getSales(id *int) (models.Salesman, error) {

	var sales models.Salesman

	var sqlStatement = `SELECT id, name, street, city, phone, cell, zip, email FROM salesmans WHERE id=$1`

	rs := Sql.QueryRow(sqlStatement, id)

	err := rs.Scan(&sales.ID, &sales.Name, &sales.Street, &sales.City, &sales.Phone, &sales.Cell, &sales.Zip, &sales.Email)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return sales, nil
	case nil:
		return sales, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty user on error
	return sales, err
}

func getAllSales() ([]models.Salesman, error) {

	var sales []models.Salesman

	var sqlStatement = `SELECT id, name, street, city, phone, cell, zip, email FROM salesmans ORDER BY name`

	rs, err := Sql.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute sales query %v", err)
	}

	defer rs.Close()

	for rs.Next() {
		var s models.Salesman

		err := rs.Scan(&s.ID, &s.Name, &s.Street, &s.City, &s.Phone, &s.Cell, &s.Zip, &s.Email)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		sales = append(sales, s)
	}

	return sales, err
}

func deleteSales(id *int) int64 {
	// create the delete sql query
	sqlStatement := `DELETE FROM salesmans WHERE id=$1`

	// execute the sql statement
	res, err := Sql.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to delete sales. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	return rowsAffected
}

func createSales(sales *models.Salesman) (int, error) {

	sqlStatement := `INSERT INTO salesmans 
	(name, street, city, phone, cell, zip, email) 
	VALUES 
	($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`

	var id int

	err := Sql.QueryRow(sqlStatement,
		sales.Name,
		sales.Street,
		sales.City,
		sales.Phone,
		sales.Cell,
		sales.Zip,
		sales.Email,
	).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to create sales. %v", err)
	}

	return id, err
}

func updateSales(id *int, sales *models.Salesman) int64 {

	sqlStatement := `UPDATE salesmans SET
	name=$2, street=$3, city=$4, phone=$5, cell=$6, zip=$7, email=$8
	WHERE id=$1`

	res, err := Sql.Exec(sqlStatement,
		id,
		sales.Name,
		sales.Street,
		sales.City,
		sales.Phone,
		sales.Cell,
		sales.Zip,
		sales.Email,
	)

	if err != nil {
		log.Fatalf("Unable to update sales. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while updating category. %v", err)
	}

	return rowsAffected
}
