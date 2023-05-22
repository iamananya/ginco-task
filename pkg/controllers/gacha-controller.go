package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iamananya/ginco-task/pkg/models"
	"github.com/iamananya/ginco-task/pkg/utils"
)

var NewUser models.User

/*
TODO----(Task completed)

[*] Handle case for empty name- return error (Error code 400), non exisiting user ID return error
[*] Autogenerate Token Don't take it from user
[*] For UpdateUser don't show response
[*] Check swagger yaml for responses.

*/

func GetUser(w http.ResponseWriter, r *http.Request) {
	newUsers := models.GetAllUsers()
	res, _ := json.Marshal(newUsers)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	ID, err := strconv.ParseInt(userId, 0, 0)
	if err != nil {
		fmt.Println("error while parsing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Invalid ID error handling done here

	userDetails, db := models.GetUserById(ID)
	if db.Error != nil {
		if db.RecordNotFound() {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	res, err := json.Marshal(userDetails)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(requestBody, user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//  Empty Username case has been handled here

	if user.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u := user.CreateUser()

	// Marshal the user object into JSON
	res, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

// UpdateUser does not show any reponse.

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var updateUser = &models.User{}
	utils.ParseBody(r, updateUser)
	vars := mux.Vars(r)
	userId := vars["userId"]
	ID, err := strconv.ParseInt(userId, 0, 0)
	if err != nil {
		fmt.Println("error while parsing")
	}
	userDetails, db := models.GetUserById(ID)
	if updateUser.Name != "" {
		userDetails.Name = updateUser.Name
	}
	if updateUser.Token != "" {
		userDetails.Token = updateUser.Token
	}
	db.Save(&userDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}
