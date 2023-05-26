package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/iamananya/ginco-task/pkg/models"
	"github.com/iamananya/ginco-task/pkg/utils"
)

var NewUser models.User
var NewCharacter models.Character

/*
TODO----(Task completed)

[*] Handle case for empty name- return error (Error code 400), non exisiting user ID return error
[*] Autogenerate Token Don't take it from user
[*] For UpdateUser don't show response
[*] Check swagger yaml for responses.

*/

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Extract the X-Token value from the request headers
	token := r.Header.Get("X-Token")

	// Use the X-Token value to query the user from the database
	user := models.GetAllUsers(token)

	if user == nil {
		// User not found, return an appropriate error response
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Return the user data as a JSON response
	res, err := json.Marshal(user)
	if err != nil {
		// Error while marshaling JSON, return an error response
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
	var updateUser models.User
	utils.ParseBody(r, &updateUser)

	// Extract the X-Token value from the request headers
	token := r.Header.Get("X-Token")

	// Retrieve the user based on the token
	user := models.GetUserByToken(token)
	fmt.Print(user)

	if user == nil {
		// User not found, return an appropriate error response
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Update the user details if the fields are not empty
	if updateUser.Name != "" {
		user[0].Name = updateUser.Name
	}
	if updateUser.Token != "" {
		user[0].Token = updateUser.Token
	}

	// Save the updated user details
	err := models.UpdateUser(&user[0])
	if err != nil {
		// Error occurred while updating user, return an appropriate error response
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func ListCharacters(w http.ResponseWriter, r *http.Request) {

	characters := models.GetAllCharacters()

	res, _ := json.Marshal(characters)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func HandleGachaDraw(w http.ResponseWriter, r *http.Request) {

	var reqBody models.GachaDrawRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received request: %+v\n", reqBody)
	characters := models.GetAllCharacters()
	characterPool := generatecharacterPool(characters)
	response := models.GachaDrawResponse{
		Results: []models.CharacterResponse{},
	}
	// fmt.Println(reqBody.NumTrials)
	for i := 0; i < reqBody.Times; i++ {
		character := models.DrawCharacter(characters, characterPool) // Simulate drawing a character
		fmt.Println(character)
		response.Results = append(response.Results, models.CharacterResponse{
			CharacterID: fmt.Sprintf("Character-%d", character.ID),
			Name:        character.Name,
		})

	}

	respBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)

}
func generatecharacterPool(characters []models.Character) []models.Character {
	var characterPool []models.Character

	for _, character := range characters {
		rarity := character.Rarity

		// Assign the probability based on rarity
		var probability int
		switch rarity {
		case "SSR":
			probability = 5
		case "SR":
			probability = 15
		case "R":
			probability = 80
		}

		// Add the character to the pool multiple times based on its probability
		poolSize := probability
		for i := 0; i < poolSize; i++ {
			characterPool = append(characterPool, character)
		}
	}

	return characterPool
}
