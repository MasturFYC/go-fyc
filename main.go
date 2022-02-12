package main

import (
	"fmt"

	"log"

	routers "go-fyc/router"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rs/cors"
)

func main() {
	mainRouter := mux.NewRouter()

	categoryRouter := mainRouter.PathPrefix("/categories").Subrouter()
	routers.CategoryRouter(categoryRouter)
	productRouter := mainRouter.PathPrefix("/products").Subrouter()
	routers.ProductRouter(productRouter)
	handler := cors.Default().Handler(mainRouter)

	fmt.Println("web server run at: http://pixel.id:8080/")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
