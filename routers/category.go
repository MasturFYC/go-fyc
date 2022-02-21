package routers

import (
	"go-fyc/middleware"

	"github.com/gorilla/mux"
)

func CategoryRouter(router *mux.Router) {

	router.HandleFunc("/", middleware.GetCategories).Methods("GET")
	router.HandleFunc("/{id}/", middleware.GetCategory).Methods("GET")
	router.HandleFunc("/{id}/", middleware.DeleteCategory).Methods("DELETE")
	router.HandleFunc("/", middleware.CreateCategory).Methods("POST")
	router.HandleFunc("/{id}/", middleware.UpdateCategory).Methods("PUT")

}
