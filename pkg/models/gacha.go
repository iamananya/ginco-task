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
	Name  string `gorm:"type:varchar(30);size:30" json:"name"`
	Token string `gorm:"type:char(30)" json:"token"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&User{})
	db.Model(&User{}).ModifyColumn("name", "varchar(30)")
	db.Model(&User{}).ModifyColumn("token", "char(30)")
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

// func DeleteUser(Id int64) User {
// 	var user User
// 	db.Where("ID=?", Id).Delete(user)
// 	return user
// }
