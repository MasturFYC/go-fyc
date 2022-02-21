package routers

import (
	"go-fyc/middleware"

	"github.com/gorilla/mux"
)

func ProductRouter(router *mux.Router) {

	router.HandleFunc("/", middleware.GetProducts).Methods("GET")
	router.HandleFunc("/", middleware.CreateProduct).Methods("POST")
	router.HandleFunc("/{id}/", middleware.GetProduct).Methods("GET")
	router.HandleFunc("/category/{id}/", middleware.GetProductsByCategory).Methods("GET")
	router.HandleFunc("/{id}/", middleware.UpdateProduct).Methods("PUT")
	router.HandleFunc("/{id}/", middleware.DeleteProduct).Methods("DELETE")

}
