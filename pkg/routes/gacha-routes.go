package routes

import (
	"github.com/gorilla/mux"
	"github.com/iamananya/ginco-task/pkg/controllers"
	"github.com/iamananya/ginco-task/pkg/middlewares"
)

var RegisterUserRoute = func(router *mux.Router) {
	router.HandleFunc("/user/", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/user/", middlewares.AuthenticationMiddleware(controllers.GetUser)).Methods("GET")
	router.HandleFunc("/user/", middlewares.AuthenticationMiddleware(controllers.UpdateUser)).Methods("PUT")

	router.HandleFunc("/characters/list/", controllers.ListCharacters).Methods("GET")
	router.HandleFunc("/gacha/draw/", controllers.HandleGachaDraw).Methods("POST")
	// router.HandleFunc("/user/characters/", controllers.CreateUserCharacter).Methods("POST")
}
