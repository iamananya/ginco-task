package models

import (
	"math/rand"
	"time"

	"github.com/iamananya/ginco-task/pkg/config"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type User struct {
	gorm.Model
	Name       string          `gorm:"type:varchar(30);size:30" json:"name"`
	Token      string          `gorm:"type:char(30)" json:"token"`
	Characters []UserCharacter `gorm:"foreignKey:UserID" json:"characters"`
}
type Character struct {
	gorm.Model
	Name              string          `gorm:"type:varchar(30);size:30" json:"name"`
	AttackPower       int             `gorm:"column:attack_power" json:"attack_power"`
	Defense           int             `gorm:"column:defense" json:"defense"`
	Speed             int             `json:"speed"`
	HitPoints         int             `gorm:"column:hit_points" json:"hit_points"`
	CriticalHitRate   float64         `gorm:"column:critical_hit_rate" json:"critical_hit_rate"`
	ElementalAffinity string          `json:"elemental_affinity"`
	Rarity            string          `json:"rarity"`
	Synergy           bool            `json:"synergy"`
	Evolution         bool            `json:"evolution"`
	Users             []UserCharacter `gorm:"foreignKey:CharacterID" json:"users"`
}

type UserCharacter struct {
	gorm.Model
	UserID            uint    `gorm:"index" json:"user_id"`
	CharacterID       uint    `gorm:"index" json:"character_id"`
	AttackPower       int     `gorm:"column:attack_power" json:"attack_power"`
	Defense           int     `gorm:"column:defense" json:"defense"`
	Speed             int     `json:"speed"`
	HitPoints         int     `gorm:"column:hit_points" json:"hit_points"`
	CriticalHitRate   float64 `gorm:"column:critical_hit_rate" json:"critical_hit_rate"`
	ElementalAffinity string  `json:"elemental_affinity"`
	Rarity            string  `json:"rarity"`
	Synergy           bool    `json:"synergy"`
	Evolution         bool    `json:"evolution"`
}

type GachaDrawRequest struct {
	NumTrials int `json:"num_trials"`
}

type GachaDrawResponse struct {
	Results []CharacterResponse `json:"results"`
}
type CharacterResponse struct {
	CharacterID string `json:"characterID"`
	Name        string `json:"name"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&User{})
	db.Model(&User{}).ModifyColumn("name", "varchar(30)")
	db.Model(&User{}).ModifyColumn("token", "char(30)")
	db.AutoMigrate(&Character{})
	db.Model(&Character{}).ModifyColumn("name", "varchar(30)")
	db.AutoMigrate(&UserCharacter{})
}

func (u *User) CreateUser() *User {
	rand.Seed(time.Now().UnixNano())

	u.Token = generateRandomString(30)
	db.NewRecord(u)
	db.Create(&u)
	return u
}

// Function defined to generate random string for token

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GetAllUsers() []User {
	var Users []User
	db.Find(&Users)
	return Users
}

func GetUserById(Id int64) (*User, *gorm.DB) {
	var getUser User
	db := db.Where("ID=?", Id).Find(&getUser)
	return &getUser, db
}

func GetAllCharacters() []Character {
	var Characters []Character
	db.Find(&Characters)
	return Characters
}
func (uc *UserCharacter) CreateUserCharacter() *UserCharacter {
	db.Create(&uc)
	return uc
}

func DrawCharacter(characters []Character) Character {
	rand.Seed(time.Now().UnixNano())

	var rarityPool []Character

	// Create a pool of characters based on rarity and generate the probabilities
	for _, character := range characters {
		rarity := character.Rarity

		// Assign the probability based on rarity
		var probability float64
		switch rarity {
		case "SSR":
			probability = 0.05
		case "SR":
			probability = 0.15
		case "R":
			probability = 0.8
		}

		// Add the character to the pool multiple times based on its probability
		for i := 0; i < int(probability*100); i++ {
			rarityPool = append(rarityPool, character)
			// fmt.Println(rarityPool)
		}
	}

	// Select a random character from the rarity pool
	index := rand.Intn(len(rarityPool))
	print()
	return rarityPool[index]
}

/*
The below function is used to find the maximum probable characters
*/
// func DrawCharacter(characters []Character) Character {
// 	rand.Seed(time.Now().UnixNano())

// 	var maxProbability float64
// 	var maxProbabilityCharacters []Character

// 	// Create a pool of characters based on rarity and generate the probabilities
// 	for _, character := range characters {
// 		rarity := character.Rarity

// 		// Assign the probability based on rarity
// 		var probability float64
// 		switch rarity {
// 		case "SSR":
// 			probability = 0.05
// 		case "SR":
// 			probability = 0.15
// 		case "R":
// 			probability = 0.8
// 		}

// 		// Update the maximum probability and reset the character pool if a higher probability is found
// 		if probability > maxProbability {
// 			maxProbability = probability
// 			maxProbabilityCharacters = []Character{character}
// 		} else if probability == maxProbability {
// 			// Add the character to the pool only if it is not already present
// 			isDuplicate := false
// 			for _, c := range maxProbabilityCharacters {
// 				if c.ID == character.ID {
// 					isDuplicate = true
// 					break
// 				}
// 			}
// 			if !isDuplicate {
// 				maxProbabilityCharacters = append(maxProbabilityCharacters, character)
// 			}
// 		}
// 	}

// 	// Select a random character from the pool of characters with the maximum probability
// 	index := rand.Intn(len(maxProbabilityCharacters))
// 	return maxProbabilityCharacters[index]
// }
