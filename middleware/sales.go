package middleware

import (
	"encoding/json"
	"log"

	"go-fyc/models"

	"net/http"

	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func GetSalesmans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db, _ := CreateConnection()

	defer db.Close()

	var salesmans []models.Salesman

	db.Find(&salesmans)

	json.NewEncoder(w).Encode(&salesmans)
}

func GetSalesman(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, _ := CreateConnection()

	defer db.Close()

	params := mux.Vars(r)

	var sales models.Salesman

	var orders []models.Order

	db.First(&sales, params["id"])

	db.Model(&sales).Related(&orders)

	sales.Orders = orders

	json.NewEncoder(w).Encode(&sales)
}

func DeleteSalesman(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	db, _ := CreateConnection()

	defer db.Close()

	var sales models.Salesman

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	db.First(&sales, id)

	db.Delete(&sales)

	var msg string

	if db.RowsAffected > 0 {
		msg = "Salesman deleted successfully"
	} else {
		msg = "Salesman can not be deleted."
	}

	res := Response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(&res)

}

func CreateSalesman(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create the postgres db connection
	db, _ := CreateConnection()

	// close the db connection
	defer db.Close()

	var sales models.Salesman

	err := json.NewDecoder(r.Body).Decode(&sales)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	db.Create(&sales)

	log.Printf("Inserted a single record %v", sales.ID)

	// return the inserted id
	json.NewEncoder(w).Encode(&sales)
}

func UpdateSalesman(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// create the postgres db connection

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	db, _ := CreateConnection()

	// close the db connection
	defer db.Close()

	var sales, body models.Salesman

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// get category first
	db.First(&sales, id)

	// update category
	updateErr := db.Model(&sales).Updates(models.Salesman{
		Name:   body.Name,
		Street: body.Street,
		City:   body.City,
		Phone:  body.Phone,
		Cell:   body.Cell,
		Zip:    body.Zip,
		Email:  body.Email,
	}).Error //.Exec(sqlStatement, cat.Name).Scan(&cat)

	if updateErr != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(Response{
			ID:      500,
			Message: "Duplicate sales name",
		})
		return
		//log.Fatalf("Unable to update category. %v", updateErr)
	}

	log.Printf("Update a single record %v\n", id)

	// return the inserted id
	json.NewEncoder(w).Encode(&sales)
}
