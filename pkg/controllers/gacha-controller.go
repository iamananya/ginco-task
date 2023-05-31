package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/iamananya/ginco-task/pkg/models"
	"github.com/iamananya/ginco-task/pkg/utils"
)

// var NewUser models.User
// var NewCharacter models.Character

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Authenticate user using x-token in headers to get user details-----
	user := r.Context().Value("user").(models.User)

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
	fmt.Print(user.Name, user.Token)
	u, err := user.CreateUser()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

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
	// Retrieve the user from the context
	user := r.Context().Value("user").(models.User)

	var updateUser models.User
	utils.ParseBody(r, &updateUser)
	fmt.Print(user)

	if updateUser.Name != "" {
		user.Name = updateUser.Name
	}
	if updateUser.Token != "" {
		user.Token = updateUser.Token
	}

	err := models.UpdateUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func ListCharacters(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Print(user.Name)
	characters, err := models.GetAllCharacters()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(characters)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func HandleGachaDraw(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Print(user.Name)
	var reqBody models.GachaDrawRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received request: %+v\n", reqBody)
	characters, err := models.GetAllCharacters()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	characterPool := generateCharacterPool(characters)
	response := models.GachaDrawResponse{
		Results: []models.CharacterResponse{},
	}
	// fmt.Println(reqBody.NumTrials)
	for i := 0; i < reqBody.Times; i++ {
		character, err := models.DrawCharacter(characters, characterPool) // Simulate drawing a character
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
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
