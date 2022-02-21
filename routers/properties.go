package routers

import (
	"go-fyc/middleware/properties"

	"github.com/gorilla/mux"
)

func PropertyRouter(router *mux.Router) {

	router.HandleFunc("/product/", properties.GetProductsProps).Methods("GET")
	router.HandleFunc("/category/", properties.GetCategoryProps).Methods("GET")

}
