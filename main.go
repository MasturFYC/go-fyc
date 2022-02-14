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
	mainRouter := mux.NewRouter()

	routers.CategoryRouter(mainRouter.PathPrefix("/categories").Subrouter())
	routers.ProductRouter(mainRouter.PathPrefix("/products").Subrouter())
	routers.SalesRouter(mainRouter.PathPrefix("/sales").Subrouter())

	handler := cors.Default().Handler(mainRouter)

	fmt.Println("web server run at local: http://localhost:8080/")
	fmt.Println("web server run at: http://pixel.id:8080/")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
