package routes

import (
	"github.com/gorilla/mux"
	"github.com/iamananya/ginco-task/pkg/controllers"
)

var RegisterUserRoute = func(router *mux.Router) {
	router.HandleFunc("/user/", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/user/", controllers.GetUser).Methods("GET")
	router.HandleFunc("/user/{userId}", controllers.GetUserById).Methods("GET")
	router.HandleFunc("/user/{userId}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/characters/", controllers.CreateCharacter).Methods("POST")
	router.HandleFunc("/characters/list/", controllers.GetCharacters).Methods("GET")
	// router.HandleFunc("/user/characters/", controllers.CreateUserCharacter).Methods("POST")

	// router.HandleFunc("/user/{userId}", controllers.DeleteUser).Methods("DELETE")

}
