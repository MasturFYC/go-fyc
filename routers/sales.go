package routers

import (
	"go-fyc/middleware"

	"github.com/gorilla/mux"
)

func SalesRouter(router *mux.Router) {

	router.HandleFunc("", middleware.GetSalesmans).Methods("GET")
	router.HandleFunc("/{id}", middleware.GetSalesman).Methods("GET")
	router.HandleFunc("/{id}", middleware.DeleteSalesman).Methods("DELETE")
	router.HandleFunc("", middleware.CreateSalesman).Methods("POST")
	router.HandleFunc("/{id}", middleware.UpdateSalesman).Methods("PUT")

}
