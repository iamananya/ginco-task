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

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Authenticate user using x-token in headers to get user details-----
	token := r.Header.Get("X-Token")
	user := models.GetAllUsers(token)

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := json.Marshal(user)
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

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Authenticate user using x-token in headers to get user details--------

	var updateUser models.User
	utils.ParseBody(r, &updateUser)
	token := r.Header.Get("X-Token")
	user := models.GetUserByToken(token)
	fmt.Print(user)

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if updateUser.Name != "" {
		user[0].Name = updateUser.Name
	}
	if updateUser.Token != "" {
		user[0].Token = updateUser.Token
	}

	err := models.UpdateUser(&user[0])
	if err != nil {
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
	characterPool := generateCharacterPool(characters)
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
func generateCharacterPool(characters []models.Character) []models.Character {
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
