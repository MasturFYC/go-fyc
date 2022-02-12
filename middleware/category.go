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

func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db, _ := CreateConnection()

	defer db.Close()

	var categories []models.Category

	db.Find(&categories)

	json.NewEncoder(w).Encode(&categories)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, _ := CreateConnection()

	defer db.Close()

	params := mux.Vars(r)

	var category models.Category

	var products []models.Product

	db.First(&category, params["id"])

	db.Model(&category).Related(&products)

	category.Products = products

	json.NewEncoder(w).Encode(&category)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	db, _ := CreateConnection()

	defer db.Close()

	var cat models.Category

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	db.First(&cat, id)

	db.Delete(&cat)

	var msg string

	if db.RowsAffected > 0 {
		msg = "Category deleted successfully"
	} else {
		msg = "Category can not be deleted."
	}

	res := Response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(&res)

}

func CreateCategory(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// create the postgres db connection
	db, _ := CreateConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning userid will return the id of the inserted user
	//sqlStatement := `INSERT INTO categories (name) VALUES ($1) RETURNING *`

	// the inserted id will store in this id
	var cat models.Category

	err := json.NewDecoder(r.Body).Decode(&cat)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// execute the sql statement
	// Scan function will save the insert id in the id
	db.Create(&cat) //.Exec(sqlStatement, cat.Name).Scan(&cat)

	// if err != nil {
	// 	log.Fatalf("Unable to insert category. %v", err)
	// }

	log.Printf("Inserted a single record %v", cat.ID)

	// return the inserted id
	json.NewEncoder(w).Encode(&cat)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {

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

	var cat, body models.Category

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// get category first
	db.First(&cat, id)

	// update category
	updateErr := db.Model(&cat).Updates(models.Category{Name: body.Name}).Error //.Exec(sqlStatement, cat.Name).Scan(&cat)

	if updateErr != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(Response{
			ID:      500,
			Message: "Duplicate category name",
		})
		return
		//log.Fatalf("Unable to update category. %v", updateErr)
	}

	log.Printf("Update a single record %v\n", id)

	// return the inserted id
	json.NewEncoder(w).Encode(&cat)
}
