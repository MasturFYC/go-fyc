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

func GetSuppliers(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	suppliers, err := getAllSupplier()

	if err != nil {
		log.Fatalf("Unable to get all suppliers. %v", err)
	}

	json.NewEncoder(w).Encode(&suppliers)
}

func GetSupplier(w http.ResponseWriter, r *http.Request) {

	EnableCors(&w)

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	suppliers, err := getSupplier(&id)

	if err != nil {
		log.Fatalf("Unable to get category. %v", err)
	}

	json.NewEncoder(w).Encode(&suppliers)
}

func DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	deletedRows := deleteSupplier(&id)

	msg := fmt.Sprintf("Supplier deleted successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := Response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func CreateSupplier(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var suppliers models.Supplier

	err := json.NewDecoder(r.Body).Decode(&suppliers)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	id, err := createSupplier(&suppliers)

	if err != nil {
		log.Fatalf("Nama supplier tidak boleh sama.  %v", err)
	}

	suppliers.ID = id

	json.NewEncoder(w).Encode(&suppliers)

}

func UpdateSupplier(w http.ResponseWriter, r *http.Request) {

	EnableCors(&w)

	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	// create the postgres db connection

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	var suppliers models.Supplier

	err := json.NewDecoder(r.Body).Decode(&suppliers)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	updatedRows := updateSupplier(&id, &suppliers)

	msg := fmt.Sprintf("Supplier updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := Response{
		ID:      updatedRows,
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func getSupplier(id *int) (models.Supplier, error) {

	var s models.Supplier

	var sqlStatement = `SELECT id, name, sales_name, street, city, phone, cell, zip, email FROM suppliers WHERE id=$1`

	rs := Sql().QueryRow(sqlStatement, id)

	err := rs.Scan(&s.ID, &s.Name, &s.SalesName, &s.Street, &s.City, &s.Phone, &s.Cell, &s.Zip, &s.Email)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return s, nil
	case nil:
		return s, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty user on error
	return s, err
}

func getAllSupplier() ([]models.Supplier, error) {

	var suppliers []models.Supplier

	var sqlStatement = `SELECT id, name, sales_name, street, city, phone, cell, zip, email FROM suppliers ORDER BY name`

	rs, err := Sql().Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute suppliers query %v", err)
	}

	defer rs.Close()

	for rs.Next() {
		var s models.Supplier

		err := rs.Scan(&s.ID, &s.Name, &s.SalesName, &s.Street, &s.City, &s.Phone, &s.Cell, &s.Zip, &s.Email)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		suppliers = append(suppliers, s)
	}

	return suppliers, err
}

func deleteSupplier(id *int) int64 {
	// create the delete sql query
	sqlStatement := `DELETE FROM suppliers WHERE id=$1`

	// execute the sql statement
	res, err := Sql().Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to delete supplier. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	return rowsAffected
}

func createSupplier(cust *models.Supplier) (int, error) {

	sqlStatement := `INSERT INTO suppliers 
	(name, sales_name, street, city, phone, cell, zip, email) 
	VALUES 
	($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id`

	var id int

	err := Sql().QueryRow(sqlStatement,
		cust.Name,
		cust.SalesName,
		cust.Street,
		cust.City,
		cust.Phone,
		cust.Cell,
		cust.Zip,
		cust.Email,
	).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to create supplier. %v", err)
	}

	return id, err
}

func updateSupplier(id *int, cust *models.Supplier) int64 {

	sqlStatement := `UPDATE suppliers SET
	name=$2, sales_name=$3 street=$4, city=$5, phone=$6, cell=$7, zip=$8, email=$9
	WHERE id=$1`

	res, err := Sql().Exec(sqlStatement,
		id,
		cust.Name,
		cust.SalesName,
		cust.Street,
		cust.City,
		cust.Phone,
		cust.Cell,
		cust.Zip,
		cust.Email,
	)

	if err != nil {
		log.Fatalf("Unable to update supplier. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while updating supplier. %v", err)
	}

	return rowsAffected
}
