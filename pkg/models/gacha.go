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
	Rarity      string `json:"rarity"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	if err := db.AutoMigrate(&User{}).Error; err != nil {
		panic("Failed to migrate User model: " + err.Error())
	}
	if err := db.Model(&User{}).ModifyColumn("name", "varchar(30)").Error; err != nil {
		panic("Failed to modify column 'name' in User model: " + err.Error())
	}
	if err := db.Model(&User{}).ModifyColumn("token", "char(30)").Error; err != nil {
		panic("Failed to modify column 'token' in User model: " + err.Error())
	}
	if err := db.AutoMigrate(&Character{}).Error; err != nil {
		panic("Failed to migrate Character model: " + err.Error())
	}
	if err := db.Model(&Character{}).ModifyColumn("name", "varchar(30)").Error; err != nil {
		panic("Failed to modify column 'name' in Character model: " + err.Error())
	}
	if err := db.AutoMigrate(&UserCharacter{}).Error; err != nil {
		panic("Failed to migrate UserCharacter model: " + err.Error())
	}
	if err := db.AutoMigrate(&GachaResult{}).Error; err != nil {
		panic("Failed to migrate GachaResult model: " + err.Error())
	}
}

func (u *User) CreateUser() (*User, error) {
	u.Token = generateRandomString(30)
	db.NewRecord(u)
	if err := db.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
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
func GetUserByToken(token string, user *User) error {
	if err := db.Where("token = ?", token).First(user).Error; err != nil {
		return err
	}
	return nil
}

func UpdateUser(user *User) error {

	if err := db.Save(user).Error; err != nil {
		return err
	}

	return nil
}

func GetAllCharacters() ([]Character, error) {
	var characters []Character
	if err := db.Find(&characters).Error; err != nil {
		return nil, err
	}
	return characters, nil
}
func (uc *UserCharacter) CreateUserCharacter() error {
	db := config.GetDB()
	if err := db.Create(uc).Error; err != nil {
		return err
	}
	return nil
}
func (gr *GachaResult) SaveGachaResult() error {
	db := config.GetDB()
	err := db.Create(gr).Error
	return err
}
