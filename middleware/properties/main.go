package properties

import (
	"encoding/json"
	mid "go-fyc/middleware"
	"go-fyc/models"
	"log"
	"net/http"
)

func GetCategoryProps(w http.ResponseWriter, r *http.Request) {
	mid.EnableCors(&w)

	props, err := getProperties("categories")

	if err != nil {
		log.Fatalf("Unable to get all property. %v", err)
	}

	json.NewEncoder(w).Encode(&props)
}

func GetProductsProps(w http.ResponseWriter, r *http.Request) {
	mid.EnableCors(&w)

	props, err := getProductProps()

	if err != nil {
		log.Fatalf("Unable to get all property. %v", err)
	}

	json.NewEncoder(w).Encode(&props)
}

func getProperties(table string) ([]models.Property, error) {
	// defer Sql().Close()
	var props []models.Property

	sqlStatement := "SELECT id, name FROM " + table + " ORDER BY name"

	rows, err := mid.Sql().Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute products query %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Property

		err := rows.Scan(&p.ID, &p.Name)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		props = append(props, p)
	}

	return props, err
}

func getProductProps() ([]models.Property, error) {
	// defer Sql().Close()
	var props []models.Property

	sqlStatement := `SELECT
		id, name, spec
	FROM products 
	ORDER BY name`

	rows, err := mid.Sql().Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute products query %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Property
		var id int64
		var name string
		var spec models.NullString

		err := rows.Scan(&id, &name, &spec)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		p.ID = id
		p.Name = mid.CombineString(name, string(spec))

		props = append(props, p)
	}

	return props, err
}
