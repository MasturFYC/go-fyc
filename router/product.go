package routers

import (
	"go-fyc/middleware"

	"github.com/gorilla/mux"
)

func ProductRouter(router *mux.Router) {

	router.HandleFunc("", middleware.GetProducts).Methods("GET")

}
