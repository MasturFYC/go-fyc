package main

import (
	"fmt"

	"log"

	"go-fyc/routers"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rs/cors"
)

func main() {

	cor := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8081",
			"http://localhost:8080",
			"http://192.168.100.2:8081",
			"http://192.168.100.2:8080",
			"http://127.0.0.1:8081",
			"http://127.0.0.1:8080",
			"http://yoga.pixel.id:8080",
			"http://yoga.pixel.id",
			"http://yoga.pixel.id:8081"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"Accept", "Accept-Language", "Content-Type"},
		//AllowCredentials: true,
		//Debug: true,
	})

	mainRouter := mux.NewRouter()

	routers.CategoryRouter(mainRouter.PathPrefix("/api/categories/").Subrouter())
	routers.ProductRouter(mainRouter.PathPrefix("/api/products/").Subrouter())
	routers.SalesRouter(mainRouter.PathPrefix("/api/sales/").Subrouter())
	routers.PropertyRouter(mainRouter.PathPrefix("/api/properties/").Subrouter())
	routers.CustomerRouter(mainRouter.PathPrefix("/api/customers/").Subrouter())

	handler := cor.Handler(mainRouter)

	fmt.Println("web server run at local: http://localhost:8080/")
	fmt.Println("web server run at: http://pixel.id:8080/")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
