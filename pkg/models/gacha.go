package models

import (
	"crypto/rand"
	"errors"
	"fmt"
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
	Name              string           `gorm:"type:varchar(30);size:30" json:"name"`
	AttackPower       int              `gorm:"column:attack_power" json:"attack_power"`
	Defense           int              `gorm:"column:defense" json:"defense"`
	Speed             int              `json:"speed"`
	HitPoints         int              `gorm:"column:hit_points" json:"hit_points"`
	CriticalHitRate   float64          `gorm:"column:critical_hit_rate" json:"critical_hit_rate"`
	ElementalAffinity string           `json:"elemental_affinity"`
	Rarity            string           `json:"rarity"`
	Synergy           bool             `json:"synergy"`
	Evolution         bool             `json:"evolution"`
	Users             []*UserCharacter `gorm:"foreignKey:CharacterID" json:"users"`
	GachaCharacters   []*GachaCharacter
}

type UserCharacter struct {
	gorm.Model
	UserID            uint    `gorm:"index" json:"user_id"`
	CharacterID       uint    `gorm:"index" json:"character_id"`
	GachaID           uint    `gorm:"index" json:"gacha_id"`
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
type Gacha struct {
	gorm.Model
	Name       string  `gorm:"type:varchar(30);size:30" json:"name"`
	RarityR    float64 `gorm:"column:rarity_r" json:"rarity_r"`
	RaritySR   float64 `gorm:"column:rarity_sr" json:"rarity_sr"`
	RaritySSR  float64 `gorm:"column:rarity_ssr" json:"rarity_ssr"`
	Characters []*GachaCharacter
}

type GachaCharacter struct {
	gorm.Model
	GachaID     uint `gorm:"index" json:"gacha_id"`
	CharacterID uint `gorm:"index" json:"character_id"`
	Character   *Character
}
type GachaResult struct {
	gorm.Model
	CharacterID   uint   `gorm:"index" json:"character_id"`
	CharacterName string `json:"character_name"`
}

type GachaDrawRequest struct {
	UserID  int `json:"user_id"`
	Times   int `json:"times"`
	GachaID int `json:"gacha_id"`
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
	// Add necessary columns to existing tables
	db.Exec("ALTER TABLE gacha_results ADD COLUMN IF NOT EXISTS gacha_id INT")
	db.Exec("ALTER TABLE gacha_results ADD COLUMN IF NOT EXISTS character_id INT")

	// Migrate gacha tables
	if err := db.AutoMigrate(&Gacha{}).Error; err != nil {
		panic("Failed to migrate Gacha model: " + err.Error())
	}
	if err := db.AutoMigrate(&GachaCharacter{}).Error; err != nil {
		panic("Failed to migrate GachaCharacter model: " + err.Error())
	}

	// Add necessary columns to existing tables
	db.Exec("ALTER TABLE gacha_characters ADD COLUMN IF NOT EXISTS gacha_id INT")
	db.Exec("ALTER TABLE gacha_characters ADD COLUMN IF NOT EXISTS character_id INT")

	// Add 'gacha_id' column to user_characters table
	db.Exec("ALTER TABLE user_characters ADD COLUMN IF NOT EXISTS gacha_id INT")
	db.Exec("ALTER TABLE user_characters ADD CONSTRAINT fk_gacha_id FOREIGN KEY (gacha_id) REFERENCES gachas(id) ON DELETE RESTRICT ON UPDATE RESTRICT")
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
func CreateUserCharacterBatch(userCharacters []*UserCharacter) error {
	tx := db.Begin()
	batchSize := 1000 // Adjust the batch size as per your requirements

	for i := 0; i < len(userCharacters); i += batchSize {
		end := i + batchSize
		if end > len(userCharacters) {
			end = len(userCharacters)
		}

		batch := userCharacters[i:end]
		if err := batchCreation(tx, batch); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func batchCreation(tx *gorm.DB, batch []*UserCharacter) error {
	for _, uc := range batch {
		if err := tx.Create(&uc).Error; err != nil {
			return err
		}
	}
	return nil
}
func (gr *GachaResult) SaveGachaResult() error {
	db := config.GetDB()
	err := db.Create(gr).Error
	return err
}

func GetGachaByID(gachaID int) (*Gacha, error) {
	// Create a new Gacha object to store the retrieved data
	gacha := &Gacha{}

	// Retrieve the gacha by ID from the database
	if err := db.First(gacha, gachaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Handle case when no gacha with the specified ID is found
			return nil, fmt.Errorf("gacha with ID %d not found", gachaID)
		}
		// Handle other errors that may occur during the database query
		return nil, err
	}

	return gacha, nil
}
