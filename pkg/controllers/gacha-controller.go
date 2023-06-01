package controllers

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
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
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var updateUser models.User
	err := utils.ParseBody(r, &updateUser)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if updateUser.Name != "" {
		user.Name = updateUser.Name
	}
	if updateUser.Token != "" {
		user.Token = updateUser.Token
	}

	err = models.UpdateUser(&user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
	characterPool := GenerateCharacterPool(characters)
	response := models.GachaDrawResponse{
		Results: make([]models.CharacterResponse, reqBody.Times),
	}

	for i := 0; i < reqBody.Times; i++ {
		character, err := DrawCharacter(&characterPool)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		characterResponse := models.CharacterResponse{
			CharacterID: fmt.Sprintf("Character-%d", character.ID),
			Name:        character.Name,
			Rarity:      character.Rarity,
		}
		response.Results[i] = characterResponse
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

func GenerateCharacterPool(characters []models.Character) []*models.Character {
	characterPool := make([]*models.Character, 0)
	probabilityMap := make(map[uint]int)

	for _, character := range characters {
		rarity := character.Rarity
		var probability int

		switch rarity {
		case "SSR":
			probability = 5
		case "SR":
			probability = 15
		case "R":
			probability = 80
		}

		probabilityMap[character.ID] = probability
	}

	for characterID, probability := range probabilityMap {
		character := FindCharacterByID(characters, characterID)
		for i := 0; i < probability; i++ {
			characterPool = append(characterPool, character)
		}
	}

	return characterPool
}

func FindCharacterByID(characters []models.Character, characterID uint) *models.Character {
	for _, character := range characters {
		if character.ID == characterID {
			return &character
		}
	}
	return nil
}

func DrawCharacter(characterPool *[]*models.Character) (*models.Character, error) {
	poolSize := len(*characterPool)
	if poolSize == 0 {
		return nil, errors.New("empty character pool")
	}

	// Randomly select a character from the pool
	randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(poolSize)))
	if err != nil {
		return nil, errors.New("failed to generate random number")
	}
	index := int(randIndex.Int64())
	character := (*characterPool)[index]

	// Remove the selected character from the pool
	*characterPool = append((*characterPool)[:index], (*characterPool)[index+1:]...)

	return character, nil
}
