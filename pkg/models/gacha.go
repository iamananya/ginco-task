package models

import (
	"crypto/rand"
	"math/big"

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
type GachaResult struct {
	gorm.Model
	CharacterID   uint   `gorm:"index" json:"character_id"`
	CharacterName string `json:"character_name"`
}

type GachaDrawRequest struct {
	Times int `json:"times"`
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
	db.AutoMigrate(&GachaResult{})
}

func (u *User) CreateUser() *User {
	u.Token = generateRandomString(30)
	db.NewRecord(u)
	db.Create(&u)
	return u
}

// Function defined to generate random string for token---
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLength := big.NewInt(int64(len(charset)))
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, charsetLength)
		b[i] = charset[randomIndex.Int64()]
	}
	return string(b)
}

// Function used to verify user token-----
func GetUserByToken(token string) []User {
	var users []User
	db.Where("token = ?", token).Find(&users)
	return users
}

func GetAllUsers(token string) []User {

	if token != "" {
		return GetUserByToken(token)
	}
	return nil

	// If no token is provided, retrieve all users
	// var users []User
	// db.Find(&users)
	// return users
}

func UpdateUser(user *User) error {

	if err := db.Save(user).Error; err != nil {
		return err
	}

	return nil
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

func DrawCharacter(characters []Character, characterPool []Character) Character {
	randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(characterPool))))

	index := int(randIndex.Int64())
	selectedCharacter := characterPool[index]

	// Save the gacha result in the database
	gachaResult := GachaResult{
		CharacterID:   selectedCharacter.ID,
		CharacterName: selectedCharacter.Name,
	}
	_ = gachaResult.SaveGachaResult()

	return selectedCharacter
}
func (gr *GachaResult) SaveGachaResult() error {
	db := config.GetDB()
	err := db.Create(gr).Error
	return err
}
