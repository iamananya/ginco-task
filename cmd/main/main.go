package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamananya/ginco-task/pkg/middlewares"
	"github.com/iamananya/ginco-task/pkg/routes"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterUserRoute(r)
	// Apply authentication middleware to all routes
	authenticatedRouter := middlewares.AuthenticationMiddleware(r.ServeHTTP)

	// http.Handle("/", authenticatedRouter)
	log.Fatal(http.ListenAndServe("localhost:9010", authenticatedRouter))
}
