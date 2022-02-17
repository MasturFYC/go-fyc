package routers

import (
	"go-fyc/middleware"

	"github.com/gorilla/mux"
)

func CustomerRouter(router *mux.Router) {

	router.HandleFunc("", middleware.GetCustomers).Methods("GET")
	router.HandleFunc("/{id}", middleware.GetCustomer).Methods("GET")
	router.HandleFunc("/{id}", middleware.DeleteCustomer).Methods("DELETE")
	router.HandleFunc("", middleware.CreateCustomer).Methods("POST")
	router.HandleFunc("/{id}", middleware.UpdateCustomer).Methods("PUT")

}
