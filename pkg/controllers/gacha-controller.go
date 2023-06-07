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
		http.Error(w, "Internal Server Error to retrieve user", http.StatusInternalServerError)
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

	gacha, err := models.GetGachaByID(reqBody.GachaID)
	if err != nil {
		http.Error(w, "Invalid gacha ID", http.StatusBadRequest)
		return
	}

	characters, err := models.GetAllCharacters()
	if err != nil {
		http.Error(w, "Internal Server Error in characters", http.StatusInternalServerError)
		return
	}
	var characterPointers []*models.Character
	for _, character := range characters {
		charac := character
		characterPointers = append(characterPointers, &charac)
	}
	characterPool := generateCharacterPool(characterPointers, gacha.RarityR, gacha.RaritySR, gacha.RaritySSR, reqBody.Times)
	response := models.GachaDrawResponse{
		Results: make([]models.CharacterResponse, reqBody.Times),
	}
	// Create a slice to store the user characters for batch insert
	userCharacters := make([]*models.UserCharacter, reqBody.Times)

	for i := 0; i < reqBody.Times; i++ {
		character, err := drawCharacter(&characterPool)
		if err != nil {
			http.Error(w, "Internal Server Error in drawing characters", http.StatusInternalServerError)
			return
		}
		userCharacter := models.UserCharacter{
			UserID:            user.ID,
			CharacterID:       character.ID,
			GachaID:           gacha.ID,
			AttackPower:       character.AttackPower,
			Defense:           character.Defense,
			Speed:             character.Speed,
			HitPoints:         character.HitPoints,
			CriticalHitRate:   character.CriticalHitRate,
			ElementalAffinity: character.ElementalAffinity,
			Rarity:            character.Rarity,
			Synergy:           character.Synergy,
			Evolution:         character.Evolution,
		}
		// Store the user character in the batch slice
		userCharacters[i] = &userCharacter

		characterResponse := models.CharacterResponse{
			CharacterID: fmt.Sprintf("Character-%d", character.ID),
			Name:        character.Name,
			Rarity:      character.Rarity,
		}
		response.Results[i] = characterResponse
	}
	// Batch insert user characters into the database
	if err := models.CreateUserCharacterBatch(userCharacters); err != nil {
		http.Error(w, "Internal Server Error in creating batch", http.StatusInternalServerError)
		return
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

func generateCharacterPool(characters []*models.Character, rarityR, raritySR, raritySSR float64, numDraws int) []*models.Character {
	characterPool := make([]*models.Character, 0)

	characters_all := make([]string, 0)

	rarityProbabilities := map[string]float64{
		"SSR": raritySSR,
		"SR":  raritySR,
		"R":   rarityR,
	}

	for _, character := range characters {
		rarity := character.Rarity
		draws := rarityProbabilities[rarity]
		rarityDraws := int(draws * float64(numDraws))
		for i := 0; i < rarityDraws; i++ {
			characters_all = append(characters_all, character.Name)
		}
	}
	poolSize := len(characters_all)
	for i := 0; i < poolSize; i++ {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(poolSize-i)))
		if err != nil {
			fmt.Printf("Failed to generate random number: %s\n", err.Error())
			return characterPool
		}
		index := int(randIndex.Int64())
		characterName := characters_all[index]
		character := findCharacterByName(characters, characterName)
		characterPool = append(characterPool, character)
		// Remove the selected character from the pool
		characters_all[index] = characters_all[poolSize-i-1]
		characters_all = characters_all[:poolSize-i-1]

	}

	return characterPool
}
func findCharacterByName(characters []*models.Character, name string) *models.Character {
	for _, character := range characters {
		if character.Name == name {
			return character
		}
	}
	return nil
}
func drawCharacter(characterPool *[]*models.Character) (*models.Character, error) {
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
	// Used the reference of character pool to update the size of the pool
	*characterPool = append((*characterPool)[:index], (*characterPool)[index+1:]...)
	return character, nil
}
