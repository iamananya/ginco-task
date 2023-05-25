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

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&User{})
	db.Model(&User{}).ModifyColumn("name", "varchar(30)")
	db.Model(&User{}).ModifyColumn("token", "char(30)")
	db.AutoMigrate(&Character{})
	db.Model(&Character{}).ModifyColumn("name", "varchar(30)")

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

func (c *Character) CreateCharacter() *Character {
	db.Create(&c)
	return c
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

// func DeleteUser(Id int64) User {
// 	var user User
// 	db.Where("ID=?", Id).Delete(user)
// 	return user
// }
