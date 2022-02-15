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

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	customers, err := getAllCustomer()

	if err != nil {
		log.Fatalf("Unable to get all customers. %v", err)
	}

	json.NewEncoder(w).Encode(&customers)
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {

	EnableCors(&w)

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	customers, err := getCustomer(&id)

	if err != nil {
		log.Fatalf("Unable to get category. %v", err)
	}

	json.NewEncoder(w).Encode(&customers)
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	deletedRows := deleteCustomer(&id)

	msg := fmt.Sprintf("Customer deleted successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := Response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var customers models.Customer

	err := json.NewDecoder(r.Body).Decode(&customers)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	id, err := createCustomer(&customers)

	if err != nil {
		log.Fatalf("Nama customers tidak boleh sama.  %v", err)
	}

	customers.ID = id

	json.NewEncoder(w).Encode(&customers)

}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {

	EnableCors(&w)

	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	// create the postgres db connection

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	var customers models.Customer

	err := json.NewDecoder(r.Body).Decode(&customers)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	updatedRows := updateCustomer(&id, &customers)

	msg := fmt.Sprintf("Customer updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := Response{
		ID:      updatedRows,
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func getCustomer(id *int) (models.Customer, error) {

	var customers models.Customer

	var sqlStatement = `SELECT id, name, street, city, phone, cell, zip, email FROM customers WHERE id=$1`

	rs := Sql().QueryRow(sqlStatement, id)

	err := rs.Scan(&customers.ID, &customers.Name, &customers.Street, &customers.City, &customers.Phone, &customers.Cell, &customers.Zip, &customers.Email)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return customers, nil
	case nil:
		return customers, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty user on error
	return customers, err
}

func getAllCustomer() ([]models.Customer, error) {

	var customers []models.Customer

	var sqlStatement = `SELECT id, name, street, city, phone, cell, zip, email FROM customers ORDER BY name`

	rs, err := Sql().Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute customers query %v", err)
	}

	defer rs.Close()

	for rs.Next() {
		var s models.Customer

		err := rs.Scan(&s.ID, &s.Name, &s.Street, &s.City, &s.Phone, &s.Cell, &s.Zip, &s.Email)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		customers = append(customers, s)
	}

	return customers, err
}

func deleteCustomer(id *int) int64 {
	// create the delete sql query
	sqlStatement := `DELETE FROM customers WHERE id=$1`

	// execute the sql statement
	res, err := Sql().Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to delete customer. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	return rowsAffected
}

func createCustomer(cust *models.Customer) (int, error) {

	sqlStatement := `INSERT INTO customers 
	(name, street, city, phone, cell, zip, email) 
	VALUES 
	($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`

	var id int

	err := Sql().QueryRow(sqlStatement,
		cust.Name,
		cust.Street,
		cust.City,
		cust.Phone,
		cust.Cell,
		cust.Zip,
		cust.Email,
	).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to create customer. %v", err)
	}

	return id, err
}

func updateCustomer(id *int, cust *models.Customer) int64 {

	sqlStatement := `UPDATE customers SET
	name=$2, street=$3, city=$4, phone=$5, cell=$6, zip=$7, email=$8
	WHERE id=$1`

	res, err := Sql().Exec(sqlStatement,
		id,
		cust.Name,
		cust.Street,
		cust.City,
		cust.Phone,
		cust.Cell,
		cust.Zip,
		cust.Email,
	)

	if err != nil {
		log.Fatalf("Unable to update customer. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while updating customer. %v", err)
	}

	return rowsAffected
}
